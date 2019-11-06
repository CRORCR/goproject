package main

import(
	"fmt"
)


func main() {
	weibo_test()
}

const (
	HongMing = 1 << 0
	DaRen = 1 << 1
	Vip = 1 << 2
)

type User struct {
	name string
	flag uint8
}

func set_flag(user User, isSet bool, flag uint8) User {
	if isSet == true {
		user.flag = user.flag | flag
	} else {
		user.flag = user.flag ^ flag
	}
	return user
}


func is_set_flag(user User, flag uint8) bool {
	result := user.flag & flag
	return result == 1
}

func weibo_test() {
	var user User
	user.name = "test01"
	user.flag = 0

	result := is_set_flag(user, HongMing)
	fmt.Printf("user is hongming:%t\n", result)

	user = set_flag(user, true, HongMing)
	result = is_set_flag(user, HongMing)
	fmt.Printf("user is hongming:%t\n", result)

	user = set_flag(user, false,HongMing )
	result = is_set_flag(user, HongMing)
	fmt.Printf("user is hongming:%t\n", result)
}