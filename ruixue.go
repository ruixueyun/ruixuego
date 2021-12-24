// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

// Package ruixuego 瑞雪服务端 SDK
package ruixuego

import "fmt"

const Version = "2.0.0"

var config *Config

// Init 初始化 SDK
func Init(conf *Config) error {
	config = conf
	defaultClient = NewClient()
	return loadAppKeys(config.AppKeys)
}

func loadAppKeys(m map[string]map[string]string) error {
	for appID, channelKeys := range m {
		for channelID, appKey := range channelKeys {
			err := AddAESKey(appID, channelID, []byte(appKey))
			if err != nil {
				return fmt.Errorf("invalid appkey: %s, appid: %s, channelid: %s, error: %s",
					appKey, appID, channelID, err.Error())
			}
		}

	}
	return nil
}

func GetDefaultClient() *Client {
	return defaultClient
}
