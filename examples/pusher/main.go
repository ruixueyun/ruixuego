// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package main

import (
	"fmt"

	"github.com/ruixueyun/ruixuego"
)

func main() {
	err := ruixuego.Init(&ruixuego.Config{
		APIDomain: "http://domain.com",
		CPKey:     "00000000000000000000000",
		CPID:      1000049,
	})

	err = ruixuego.GetDefaultClient().PusherPush(&ruixuego.PusherPushReq{
		PushInfo: ruixuego.PusherPushReqInfo{
			Title:          "鱼啊鱼啊鱼", // 推送标题
			Content:        "鱼鱼快动呀", // 推送内容
			Classification: 0,       // 0 营销广告类消息  1 系统类消息
			Action:         0,       // 0 打开app 1 打开url 2 打开自定义页面
			Payload:        "",      // action=1 时 传url链接 action=2 传自定义页面IntentURI
		},
		Target: ruixuego.PusherPushReqTarget{
			Type:   1, // 推送目标类型 0 全部设备 1 openid 2 openid列表
			Openid: []string{"openid"},
		},
	}, "264", "204") // "264", "204")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("done")
}
