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
		"abcdef",
		"123456",
		"game",
		map[string]interface{}{
			ruixuego.PresetKeyAppID:        "123", // 设置 AppID 请用预置 Key
			ruixuego.PresetKeyChannelID:    "456", // 设置渠道 ID 请用预置 Key
			ruixuego.PresetKeySubChannelID: "789", // 设置子渠道 ID 请用预置 Key
		}, map[string]interface{}{
			"key1": "val",
		})
	if err != nil {
		panic(err)
	}

	err = ruixuego.GetDefaultClient().TrackType(
		"abcdef",
		"123456",
		"user_setonce",
		map[string]interface{}{ // 预制属性
			ruixuego.PresetKeyAppID:        "123", // 设置 AppID 请用预置 Key
			ruixuego.PresetKeyChannelID:    "456", // 设置渠道 ID 请用预置 Key
			ruixuego.PresetKeySubChannelID: "789", // 设置子渠道 ID 请用预置 Key
		}, map[string]interface{}{ // 自定义属性
			"key2": 888,
		})
	if err != nil {
		panic(err)
	}

	fmt.Println("done")
}
