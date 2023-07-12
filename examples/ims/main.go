// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package main

import (
	"fmt"

	"github.com/ruixueyun/ruixuego"
)

const (
	testAppID     = "test_product"
	testChannelID = "test_channel"
)

// apiDomain = "https://rx-api.weilemks.com"

func main() {
	err := ruixuego.Init(&ruixuego.Config{
		APIDomain: "https://domain.com",
		CPKey:     "00000000000000000000000",
		CPID:      1000005,
		ProductID: "425",
	})
	if err != nil {
		panic(err)
	}

	// ack, err := ruixuego.GetDefaultClient().IMSSendMessage(&ruixuego.IMSMessage{
	//	Option:         3,
	//	Type:           2,
	//	Sender:         "testkk2012",
	//	UUID:           "34567890991",
	//	ClientType:     256,
	//	ConversationID: "$4$worldChannel",
	//	ConvType:       4,
	//	ProductID:      testAppID,
	//	ChannelID:      testChannelID,
	//	Content:        "{\"text\":\"22222\"}",
	//	Ext:            map[string]string{"userData": "{\"game\" : \"22222\"}"},
	// })
	// if err != nil {
	//	panic(err)
	// }
	//
	// fmt.Printf("%+v\n", ack)

	// tuser, err := ruixuego.GetDefaultClient().GetRelationUser("pp", "rxuSl4QZoNk0G1HY2-Za6GlO7wO-p_ej", "rxufb2GKDBm5n4u7gYkcOLbOv7J4RrVP")
	// if err != nil {
	//	panic(err)
	// }
	//
	// fmt.Printf("%+v\n", tuser)

	tuser, err := ruixuego.GetDefaultClient().RiskGreenAsyncScan([]string{"porn"}, []*ruixuego.GreenRequestTask{{Tag: 1, URL: "https://oss.ruixueyun.com/image/4444_09ad12c35307b04e894290bf1a75d060.jpeg"}}, "1234", "")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", tuser)
}
