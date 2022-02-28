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
	// SDK 初始化
	err := ruixuego.Init(&ruixuego.Config{
		APIDomain: "http://ruixue.weiletest.com",
		AppKeys:   map[string]map[string]string{testAppID: {testChannelID: testAppKey}},
		CPKey:     "0984cde09ebe42fd167510c727f57f71",
		CPID:      1000049,
		BigData: &ruixuego.BigDataConfig{ // 要使用大数据埋点功能必须配置此参数
			AutoFlush: true,
		},
	})
	if err != nil {
		panic(err)
	}

	defer func() {
		// 使用大数据埋点功能上传数据后, 必须在程序退出前显式调用 ruixuego.Close()
		// 不然可能导致数据丢失
		fmt.Println("close result:", ruixuego.Close())
	}()

	err = ruixuego.GetDefaultClient().Track(
		"123456",
		"game",
		map[string]interface{}{
			ruixuego.PresetKeyAppID:        "123", // 设置 AppID 请用预置 Key
			ruixuego.PresetKeyChannelID:    "456", // 设置渠道 ID 请用预置 Key
			ruixuego.PresetKeySubChannelID: "789", // 设置子渠道 ID 请用预置 Key
			"key1":                         "val",
		},
		true)
	if err != nil {
		panic(err)
	}

	err = ruixuego.GetDefaultClient().Track(
		"123456",
		"game",
		map[string]interface{}{
			ruixuego.PresetKeyAppID:        "123", // 设置 AppID 请用预置 Key
			ruixuego.PresetKeyChannelID:    "456", // 设置渠道 ID 请用预置 Key
			ruixuego.PresetKeySubChannelID: "789", // 设置子渠道 ID 请用预置 Key
			"key2":                         888,
		},
		true)
	if err != nil {
		panic(err)
	}

	fmt.Println("done")
}
