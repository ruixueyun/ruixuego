// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package main

import (
	"fmt"

	"github.com/ruixueyun/ruixuego"
)

const (
	testProductID = "wltestapp"
	testChannelID = "wltestchannel"
	testAppKey    = "00000000000000000000000"
)

func main() {
	// SDK 初始化
	err := ruixuego.Init(&ruixuego.Config{
		APIDomain: "http://domain.com",
		AppKeys:   map[string]map[string]string{testProductID: {testChannelID: testAppKey}},
		CPKey:     "00000000000000000000000",
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

	dc := ruixuego.GetDefaultClient()

	// 事件埋点
	err = dc.Tracks(
		"abcdef",
		"123456",
		ruixuego.SetEvent("login"),
		ruixuego.SetPreset(dc, map[string]interface{}{
			ruixuego.PresetKeyProductID:    "123", // 设置 ProducdID 请用预置 Key
			ruixuego.PresetKeyChannelID:    "456", // 设置渠道 ID 请用预置 Key
			ruixuego.PresetKeySubChannelID: "789", // 设置子渠道 ID 请用预置 Key
		}),
		ruixuego.SetProperties(map[string]interface{}{
			"key1": "val",
		}))
	if err != nil {
		panic(err)
	}
	// 用户属性埋点
	err = dc.Tracks(
		"abcdef",
		"123456",
		ruixuego.SetUserUpdateType("user_setonce"),
		ruixuego.SetPreset(dc, map[string]interface{}{
			ruixuego.PresetKeyProductID:    "123", // 设置 ProducdID 请用预置 Key
			ruixuego.PresetKeyChannelID:    "456", // 设置渠道 ID 请用预置 Key
			ruixuego.PresetKeySubChannelID: "789", // 设置子渠道 ID 请用预置 Key
		}),
		ruixuego.SetProperties(map[string]interface{}{
			"key1": "val",
		}))
	if err != nil {
		panic(err)
	}

	fmt.Println("done")

	select {}

}
