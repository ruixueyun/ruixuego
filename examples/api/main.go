// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package main

import (
	"fmt"

	"git.jiaxianghudong.com/ruixuesdk/ruixuego"
)

const (
	testAppID  = "wltestapp"
	testAppKey = "a463deade4b15d5ac5398f97cdaeab65"
)

func main() {
	err := ruixuego.Init(&ruixuego.Config{
		APIDomain: "http://ruixue.weiletest.com",
		AppKeys:   map[string]string{testAppID: testAppKey},
		CPKey:     "0984cde09ebe42fd167510c727f57f71",
		CPID:      1000049,
	})
	if err != nil {
		panic(err)
	}
	err = ruixuego.GetDefaultClient().SetCustom(testAppID, "rxuR4bwM27Y1JQwtAQn6H39y_9VrkEgR", "123")
	if err != nil {
		panic(err)
	}

	err = ruixuego.GetDefaultClient().AddFriend("rxuR4bwM27Y1JQwtAQn6H39y_9VrkEgR", "rxufPJeZoyrX3eWuNxLMSNK5N6x04jLl", "aaaa1", "bbb13")
	if err != nil {
		panic(err)
	}

	// err = ruixuego.GetDefaultClient().DelFriend("rxuR4bwM27Y1JQwtAQn6H39y_9VrkEgR", "rxufPJeZoyrX3eWuNxLMSNK5N6x04jLl")
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println("done")
}
