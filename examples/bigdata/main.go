// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package main

import (
	"fmt"

	"git.jiaxianghudong.com/ruixuesdk/ruixuego"
)

const (
	testAppID     = "wltestapp"
	testChannelID = "wltestchannel"
	testAppKey    = "a463deade4b15d5ac5398f97cdaeab65"
)

func main() {
	err := ruixuego.Init(&ruixuego.Config{
		APIDomain: "http://ruixue.weiletest.com",
		AppKeys:   map[string]map[string]string{testAppID: {testChannelID: testAppKey}},
		CPKey:     "0984cde09ebe42fd167510c727f57f71",
		CPID:      1000049,
		BigData:   &ruixuego.BigDataConfig{},
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		fmt.Println("close result:", ruixuego.Close())
	}()
	err = ruixuego.GetDefaultClient().Track(
		"123456",
		"game",
		map[string]interface{}{
			"key1": "val",
		},
		true)
	if err != nil {
		panic(err)
	}

	err = ruixuego.GetDefaultClient().Track(
		"123456",
		"game",
		map[string]interface{}{
			"key2": 888,
		},
		true)
	if err != nil {
		panic(err)
	}

	fmt.Println("done")
}
