// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package main

import (
	"fmt"

	"git.jiaxianghudong.com/ruixuesdk/ruixuego"
)

const (
	testAppID     = "test_product"
	testChannelID = "test_channel"
)

func main() {
	err := ruixuego.Init(&ruixuego.Config{
		APIDomain: "https://ruixue.weiletest.com",
		CPKey:     "f3c7907d161764daf97fdaaea1a72261",
		CPID:      1000000,
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
