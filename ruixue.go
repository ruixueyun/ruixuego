// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

// Package ruixuego 瑞雪服务端 SDK
package ruixuego

import "fmt"

var config *Config

// Init 初始化 SDK
func Init(conf *Config) error {
	config = conf
	defaultClient = NewClient()
	return loadAppKeys(config.AppKeys)
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

func GetDefaultClient() *Client {
	return defaultClient
}
