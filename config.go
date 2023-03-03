// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"time"
)

const (
	defaultRequestTimeout      = 5000  // default timeout. unit: millisecond
	defaultTrackRequestTimeout = 10000 // default timeout. unit: millisecond
	defaultRequestConcurrency  = 1024

	bigDataDefaultCacheCapacity     = 2000             // 默认缓存容量
	bigDataDefaultBatchSize         = 20               // 默认批量发送条数
	bigDataDefaultAutoFlushInterval = 30 * time.Second // 默认自动上传间隔 30 秒
)

// Config 瑞雪配置
type Config struct {
	APIDomain    string                       `yaml:"api_domain" json:"api_domain"`       // API 接口域名
	AppKeys      map[string]map[string]string `yaml:"appkeys" json:"appkeys"`             // map[瑞雪AChanap[瑞雪ChannelID]瑞雪App密钥
	Timeout      time.Duration                `yaml:"timeout" json:"timeout"`             // 请求超时时间(毫秒)
	TrackTimeout time.Duration                `yaml:"track_timeout" json:"track_timeout"` // 请求超时时间(毫秒)
	Concurrency  int                          `yaml:"concurrency" json:"concurrency"`     // 并发请求限制
	CPID         uint32                       `yaml:"cpid" json:"cpid"`
	CPKey        string                       `yaml:"cpkey" json:"cpkey"`
	ProductID    string                       `yaml:"product_id" json:"product_id"`
	BigData      *BigDataConfig               `yaml:"bigdata" json:"bigdata"`
	_done        bool
}

func (conf *Config) done() {
	if conf._done {
		return
	}

	if conf.Timeout == 0 {
		conf.Timeout = defaultRequestTimeout
	}
	conf.Timeout *= time.Millisecond

	if conf.TrackTimeout == 0 {
		conf.TrackTimeout = defaultRequestTimeout
	}
	conf.TrackTimeout *= time.Millisecond

	if conf.Concurrency == 0 {
		conf.Concurrency = defaultRequestConcurrency
	}

	if conf.BigData != nil {
		conf.BigData.done()
	}

	conf._done = true
}

type BigDataConfig struct {
	CacheCapacity     int           `json:"cache_capacity"`      // 缓存容量
	BatchSize         int           `json:"batch_size"`          // 大数据埋点批量发送每批条数
	AutoFlushInterval time.Duration `json:"auto_flush_interval"` // 自动上传间隔, 单位秒
	AutoFlush         bool          `json:"auto_flush"`          // 是否启动自动上传
	DisableCompress   bool          `json:"disable_compress"`    // 是否禁用 GZip 压缩
	_done             bool
}

func (conf *BigDataConfig) done() {
	if conf == nil {
		return
	}
	if conf._done {
		return
	}
	if conf.CacheCapacity == 0 {
		conf.CacheCapacity = bigDataDefaultCacheCapacity
	}
	if conf.BatchSize == 0 {
		conf.BatchSize = bigDataDefaultBatchSize
	}
	if conf.AutoFlushInterval == 0 {
		conf.AutoFlushInterval = bigDataDefaultAutoFlushInterval
	}
	conf._done = true
}
