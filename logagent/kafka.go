package main

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

type Message struct {
	line string
	topic string
}

type KafkaSender struct {
	client sarama.SyncProducer
	lineChan chan *Message
}

var kafkaSender *KafkaSender

func NewKafkaSender (kafkaAddr string) (kafka *KafkaSender, err error) {
	kafka = &KafkaSender{
		lineChan: make (chan *Message, 100000),
	}
	//新建一个配置
	config := sarama.NewConfig()
	//是否等待回应  wait是等待/noresponse是不等待
	config.Producer.RequiredAcks = sarama.NoResponse
	//分区  可以把一个topic存在一个分区,kafka是分布式的集群,如果存在一个分区,只能存一台机器
	//随机选择一个分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//需要返回
	config.Producer.Return.Successes = true
	/*
	msg := &sarama.ProducerMessage{}
	msg.Topic = "nginx_log"
	msg.Value = sarama.StringEncoder("this is a good test, my message is good")
	*/
	//生成一个同步的生产者客户端 需要kafka地址可配置信息
	client, err := sarama.NewSyncProducer([]string{kafkaAddr}, config)
	if err != nil {
		logs.Error("init kafka client failed, err:%v", err)
		return
	}

	kafka.client = client
	for i := 0; i < appConfig.KafkaThreadNum; i++ {
		go kafka.sendToKafka()
	}

	return
}

func initKafka() (err error) {
	kafkaSender, err = NewKafkaSender(appConfig.kafkaAddr)
	return
}
//发送消息到kafka
func (k *KafkaSender) sendToKafka() {
	for v := range k.lineChan {
		msg := &sarama.ProducerMessage{}
		msg.Topic = v.topic
		msg.Value = sarama.StringEncoder(v.line)
		//发送消息  用之前生成的同步的客户端发送消息
		//第一个参数:消息存在那个分区
		//第二个参数:偏移量
		_, _, err := k.client.SendMessage(msg)
		if err != nil {
			logs.Error("send message to kafka failed, err:%v", err)
		}
	}
}

func (k *KafkaSender) addMessage(line string, topic string) (err error) {
	k.lineChan <- &Message{line:line, topic:topic}
	return
}
