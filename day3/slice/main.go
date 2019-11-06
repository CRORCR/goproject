package main


import(
	"fmt"
)


func main(){
	//test1()
	//test3()

	//testSliceCap()
	//testStrReverseUtf8()
	//arr()
}

func arr(){
	var a [10]int
	a[0] = 1
	a[1] = 2
	a[2] = 3

	//以下切片,会修改原值
	result := Sum(a[:])
	//数组传递,不会修改原值
	//result := SumArray(a)
	fmt.Printf("sum=%d\n", result)
	fmt.Printf("a:%#v\n", a)
}

func test1(){
	//这要注意,后续对a的修改,会改变b的值
	var a [5]int 
	b := a[1:3]
	a[0] = 100
	a[1] = 200
	fmt.Printf("b:%#v\n", b)
}

func test3(){
	var a [5]int 
	b := a[1:3]
	//越界访问会panic
	b[100] = 100
	
	fmt.Printf("b:%#v\n", a)
}

func Sum(b []int) int {
	var sum int
	for i := 0; i < len(b); i++ {
		sum = sum + b[i]
	}

	b[0] = 100
	return sum
}

func SumArray(b [10]int) int {
	var sum int
	for i := 0; i < len(b); i++ {
		sum = sum + b[i]
	}

	b[0] = 100
	return sum
}

func testSliceCap() {
	a := make([]int, 5, 10)
	a[4] = 100
	b := a[2:3] //b的长度是1,但是容量是8
	//b[9] = 100 //会角标越界

	fmt.Printf("a=%#v, len(a) = %d, cap(a)=%d\n", a, len(a), cap(a))
	fmt.Printf("b=%#v, len(b) = %d, cap(b)=%d\n", b, len(b), cap(b))
}

func testStrReverseUtf8() string{
	str := "hello world我们爱中国"
	b := []rune(str)

	for i := 0; i < len(b)/2;i++ {
		b[i], b[len(b)-i-1] = b[len(b)-i-1], b[i]
	}

	str1 := string(b)
	fmt.Println(str1)

	fmt.Printf("len(str)=%d, len(rune)=%d\n", len(str), len(b))
	return str1
}
