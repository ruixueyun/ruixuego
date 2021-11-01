// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

// Package ruixuego 瑞雪服务端 SDK
package ruixuego

import "fmt"

// Init 初始化 SDK
func Init(conf *Config) error {
	return loadAppKeys(conf.AppKeys)
}

func loadAppKeys(m map[string]string) error {
	for appID, appKey := range m {
		err := AddAESKey(appID, []byte(appKey))
		if err != nil {
			return fmt.Errorf("invalid appkey: %s, appid: %s, error: %s",
				appKey, appID, err.Error())
		}
	}
	return nil
}
