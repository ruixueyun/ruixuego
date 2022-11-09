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
		APIDomain: "https://rx-api.weilemks.com",
		CPKey:     "236ad548c691522990bafb4990291a53",
		CPID:      1000005,
	})
	if err != nil {
		panic(err)
	}

	ack, err := ruixuego.GetDefaultClient().IMSSendMessage(&ruixuego.IMSMessage{
		Option:         3,
		Type:           2,
		ConversationID: "$4$worldChannel",
		ProductID:      testAppID,
		ChannelID:      testChannelID,
		Content:        "{\"text\":\"message content\"}",
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", ack)
}
