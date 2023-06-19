// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package main

import (
	"fmt"

	"github.com/ruixueyun/ruixuego"
)

const ciphertext = "x3kcmn9xUXwj3hMzz4O7FFTvlpQTB2rsPXBREal4PpFXvjVU+sRKBUYfPbg3lcIUgidIIZWBYuMmdkwJBtk0BfzLo2aFTaER8e9Bc0bjpIM="

func main() {
	err := ruixuego.Init(&ruixuego.Config{})
	if err != nil {
		panic(err)
	}

	err = ruixuego.AddAESKey("test", "test", []byte("098f6bcd4621d373cade4e832627b4f6"))
	if err != nil {
		panic(err)
	}

	openIDData, err := ruixuego.DecryptOpenIDData("test", "test", ciphertext)
	if err != nil {
		panic(err)
	}

	fmt.Printf(fmt.Sprintf("openid: %+v, ts: %+v\n", openIDData, openIDData.Timestamp))
}
