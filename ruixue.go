// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

// Package ruixuego 瑞雪服务端 SDK
package ruixuego

import "fmt"

const Version = "v0.1.31"

var config *Config

// Init 初始化 SDK
func Init(conf *Config) (err error) {
	conf.done()

	config = conf
	defaultClient, err = NewClient()
	if err != nil {
		return err
	}
	return loadAppKeys(config.AppKeys)
}

func loadAppKeys(m map[string]map[string]string) error {
	for productID, channelKeys := range m {
		for channelID, appKey := range channelKeys {
			err := AddAESKey(productID, channelID, []byte(appKey))
			if err != nil {
				return fmt.Errorf("invalid appkey: %s, productid: %s, channelid: %s, error: %s",
					appKey, productID, channelID, err.Error())
			}
		}

	}
	return nil
}

func GetDefaultClient() *Client {
	return defaultClient
}

func Close() error {
	return defaultClient.Close()
}
