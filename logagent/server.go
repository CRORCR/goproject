package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
)

var waitGroup sync.WaitGroup
//可以收集多个日志,新建一个结构
//一个tail挂了,需要记录偏移量,下次启动再接着收集,每个tail一个偏移量 offset
type TailObj struct {
	tail     *tail.Tail
	secLimit *SecondLimit
	offset   int64
	/*
		filename string
		service string
		sendRate int
	*/
	logConf  LogConfig
	exitChan chan bool
}
//保存所有的tail实例
type TailMgr struct {
	//使用map可以检测重复,如果重复就不管,很好的去重
	tailObjMap map[string]*TailObj
	lock       sync.Mutex
}

var tailMgr *TailMgr

func NewTailMgr() *TailMgr {
	return &TailMgr{
		tailObjMap: make(map[string]*TailObj, 16),
	}
}
//如果配置文件是动态的,随时更改,就需要加锁
func (t *TailMgr) AddLogFile(conf LogConfig) (err error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	//如果存在,就不添加 避免重复收集
	_, ok := t.tailObjMap[conf.LogPath]
	if ok {
		err = fmt.Errorf("duplicate filename:%s", conf.LogPath)
		return
	}
	//如果不存在,就初始化一个tail实例
	//location参数作用:当程序奔溃了,重新tailf读取数据的时候就会用到,否则每次都是从头读取,当程序挂了,会写入
	//kafka两遍,对于日志分析来说,这个日志其实就乱了,而且也浪费资源
	tail, err := tail.TailFile(conf.LogPath, tail.Config{
		ReOpen:    true, //重新打开
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, //从哪个位置读取
		MustExist: false,
		Poll:      true,//轮循查询
	})

	tailObj := &TailObj{
		secLimit: NewSecondLimit(int32(conf.SendRate)),
		logConf:  conf,
		offset:   0,
		tail:     tail,
		exitChan: make(chan bool, 1),
	}

	t.tailObjMap[conf.LogPath] = tailObj
	go tailObj.readLog()
	return
}

func (t *TailMgr) reloadConfig(logConfArr []LogConfig) (err error) {

	for _, conf := range logConfArr {
		tailObj, ok := t.tailObjMap[conf.LogPath]
		if !ok {
			err = t.AddLogFile(conf)
			if err != nil {
				logs.Error("add log file failed, err:%v", err)
				continue
			}
			continue
		}

		tailObj.logConf = conf
		tailObj.secLimit.limit = int32(conf.SendRate)
		t.tailObjMap[conf.LogPath] = tailObj
	}

	//处理删除的日记收集配置
	for key, tailObj := range t.tailObjMap {
		var found = false
		for _, newValue := range logConfArr {
			if key == newValue.LogPath {
				found = true
				break
			}
		}
		if found == false {
			logs.Warn("log path:%s is remove", key)
			tailObj.exitChan <- true
			delete(t.tailObjMap, key)
		}
	}
	return
}

func (t *TailMgr) Process() {
	logChan := GetLogConfChan()
	for conf := range logChan {
		logs.Debug("log conf:%v", conf)
		var logConfArr []LogConfig
		err := json.Unmarshal([]byte(conf), &logConfArr)
		if err != nil {
			logs.Error("unmarshal failed, err:%v conf:%s", err, conf)
			continue
		}

		err = t.reloadConfig(logConfArr)
		if err != nil {
			logs.Error("reload config from etcd failed, err:%v", err)
			continue
		}

		logs.Debug("reload from etcd succ, config:%v", logConfArr)

	}

	/*
		for _, tailObj := range t.tailObjMap {
			waitGroup.Add(1)
			go tailObj.readLog()
		}
	*/
}

func (t *TailObj) readLog() {
	for line := range t.tail.Lines {
		if line.Err != nil {
			logs.Error("read line failed, err:%v", line.Err)
			continue
		}

		str := strings.TrimSpace(line.Text)
		if len(str) == 0 || str[0] == '\n' {
			continue
		}
		kafkaSender.addMessage(line.Text, t.logConf.Topic)
		t.secLimit.Add(1)
		t.secLimit.Wait()

		select {
		case <-t.exitChan:
			logs.Warn("tail obj is exited, config:%v", t.logConf)
			return
		default:
		}
	}
	waitGroup.Done()
}

func RunServer() {
	tailMgr = NewTailMgr()
	/*
		var logfiles []string
		for _, filename := range logfiles {
			err := tailMgr.AddLogFile(filename)
			if err != nil {
				logs.Error("add log file %s failed, err:%v", filename, err)
				continue
			}
			logs.Debug("add log file %s succ", filename)
		}
	*/
	tailMgr.Process()
	waitGroup.Wait()
}
