package main

import (
	"strconv"
)
var a string

func test(x int, y string) string {
	var res string
	//strconv.Itoa 就是将 int 类型 转成 stirng
	res = strconv.Itoa(x) + y
	return res
}

func add(x int, y string) string {
	//strconv.Atoi 就是将 string 类型 转成 int
	i, err := strconv.Atoi(y)
	if err != nil {
		panic(err)
	}
	i = i + x
	var res string
	res=strconv.Itoa(i)
	return res
}

