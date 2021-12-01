// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package httpclient

import "time"

const (
	defaultRequestTimeout     = 5000 // default timeout. unit: millisecond
	defaultRequestConcurrency = 1024
)

type Option func(opts *options)

func getDefaultOptions() *options {
	return &options{
		Timeout:     defaultRequestTimeout * time.Millisecond,
		Concurrency: defaultRequestConcurrency,
	}
}

type options struct {
	Timeout     time.Duration // 请求超时时间(毫秒)
	Concurrency int           // 并发请求限制
}

// WithConcurrency 设置请求限制
func WithConcurrency(v int) Option {
	return func(o *options) {
		o.Concurrency = v
	}
}

// WithTimeout 设置请求超时时间(毫秒)
func WithTimeout(v int64) Option {
	return func(o *options) {
		o.Timeout = time.Duration(v) * time.Millisecond
	}
}
