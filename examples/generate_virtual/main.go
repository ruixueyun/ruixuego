// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package main

import (
	"fmt"

	"git.jiaxianghudong.com/ruixuesdk/ruixuego"
)

func main() {
	err := ruixuego.Init(&ruixuego.Config{})
	if err != nil {
		panic(err)
	}

	err = ruixuego.AddAESKey(
		"wltestapp", "wltestapp", []byte("a463deade4b15d5ac5398f97cdaeab65"))
	if err != nil {
		panic(err)
	}

	loginData, err := ruixuego.GenerateVirtualLoginData(
		"11111", "wltestapp", "wltestapp", "1234567")
	if err != nil {
		panic(err)
	}

	fmt.Println(loginData)
}
