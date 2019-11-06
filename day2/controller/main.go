package main


import(
	"fmt"
	"strings"
)

func StrOperator (){
	str3 := "the,character,represented,by,the corresponding Unicode code point"
	result := strings.Split(str3, ",")
	fmt.Printf("result:%v\n", result)

	str5 := strings.Join(result, ",")
	fmt.Printf("str5:%s\n", str5)

	str4 := "baidu.com"
	index := strings.Index(str4, "du")

	if ret := strings.HasPrefix(str4, "http://"); ret == false {
		str4 = "http://" +str4
	}
}

func TestScanf() {
	var number int
	var str string 
	fmt.Scanf("%d %s", &number, &str)
	fmt.Println(number, str)
}
