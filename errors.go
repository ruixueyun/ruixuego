// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import "errors"

var (
	ErrInvalidCPID     = errors.New("invalid cpid")
	ErrInvalidOpenID   = errors.New("invalid openid")
	ErrInvalidAppID    = errors.New("invalid appid")
	ErrAppKeyNotExistx = errors.New("appkey not exists")
	ErrInvalidType     = errors.New("invalid type")

	errProducerShutdown = errors.New("producer already shut down")
)
