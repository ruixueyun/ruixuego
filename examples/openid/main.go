// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package main

import (
	"fmt"

	"git.jiaxianghudong.com/ruixue/sdk/ruixuego"
)

const ciphertext = "x3kcmn9xUXwj3hMzz4O7FFTvlpQTB2rsPXBREal4PpFXvjVU+sRKBUYfPbg3lcIUgidIIZWBYuMmdkwJBtk0BfzLo2aFTaER8e9Bc0bjpIM="

func main() {
	err := ruixuego.Init(&ruixuego.Config{})
	if err != nil {
		panic(err)
	}

	err = ruixuego.AddAESKey("test", []byte("098f6bcd4621d373cade4e832627b4f6"))
	if err != nil {
		panic(err)
	}

	openIDData, err := ruixuego.DecryptOpenIDData("test", ciphertext)
	if err != nil {
		panic(err)
	}

	fmt.Printf("appid: %s, token: %s, openid: %s, ts: %d\n", openIDData.AppID, openIDData.Token, openIDData.OpenID, openIDData.Timestamp)
}
