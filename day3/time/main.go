package main

import(
	"time"
	"fmt"
)

func main() {
	start := time.Now().UnixNano()
	/*
	业务代码
	*/
	time.Sleep(10*time.Millisecond)
	end := time.Now().UnixNano()
	cost := (end - start)/1000
	fmt.Printf("cost:%dus\n", cost)
}