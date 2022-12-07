// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package main

import (
	"fmt"

	"ruixuego"
)

const (
	testAppID     = "test_product"
	testChannelID = "test_channel"
)

//apiDomain = "https://rx-api.weilemks.com"

func main() {
	err := ruixuego.Init(&ruixuego.Config{
		APIDomain: "https://ruixue.weiletest.com",
		CPKey:     "4c6d8d2af29e1fbda9e1fc992df13141",
		CPID:      1000000,
	})
	if err != nil {
		panic(err)
	}

	ack, err := ruixuego.GetDefaultClient().IMSSendMessage(&ruixuego.IMSMessage{
		Option:         3,
		Type:           2,
		Sender:         "testkk2012",
		UUID:           "34567890991",
		ClientType:     256,
		ConversationID: "$4$worldChannel",
		ConvType:       4,
		ProductID:      testAppID,
		ChannelID:      testChannelID,
		Content:        "{\"text\":\"22222\"}",
		Ext:            map[string]string{"userData": "{\"game\" : \"22222\"}"},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", ack)
}
