// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

// Config 瑞雪配置
type Config struct {
	APIDomain string                       `yaml:"api_domain" json:"api_domain"` // API 接口域名
	AppKeys   map[string]map[string]string `yaml:"appkeys" json:"appkeys"`       // map[瑞雪AChanap[瑞雪ChannelID]瑞雪App密钥
	CPKey     string                       `yaml:"cpkey" json:"cpkey"`
	CPID      uint32                       `yaml:"cpid" json:"cpid"`
}
