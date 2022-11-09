// Copyright (c) 2022. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"net/url"
	"strconv"
	"sync"
)

var (
	sha1Pool = &sync.Pool{
		New: func() interface{} {
			return sha1.New() // nolint:gosec
		},
	}
)

// GetSign 获取请求签名，该签名通过 sha1(TraceID+Timestamp+CPKey) 得来
func GetSign(traceID, ts string) string {
	h := sha1Pool.Get().(hash.Hash)
	_, _ = h.Write([]byte(traceID + ts + config.CPKey))
	ret := hex.EncodeToString(h.Sum(nil))
	h.Reset()
	sha1Pool.Put(h)
	return ret
}

func GetLoginResultSign(result *LoginResult, signFields []string) string {
	params := url.Values{}
	for _, field := range signFields {
		switch field {
		case SignFieldNickname:
			params.Set(field, result.Nickname)
		case SignFieldAvatar:
			params.Set(field, result.Avatar)
		case SignFieldOpenID:
			params.Set(field, result.OpenID)
		case SignFieldRegion:
			params.Set(field, result.Region)
		case SignFieldSex:
			params.Set(field, strconv.FormatInt(int64(result.Sex), 10))
		case SignFieldAge:
			params.Set(field, strconv.FormatInt(int64(result.Age), 10))
		}
	}

	signSource := params.Encode() + config.CPKey

	h := sha1Pool.Get().(hash.Hash)
	h.Write([]byte(signSource))

	sign := hex.EncodeToString(h.Sum(nil))

	h.Reset()
	sha1Pool.Put(h)

	return sign
}

type LoginResult struct {
	Nickname string `json:"nickname"` // 昵称
	Avatar   string `json:"avatar"`   // 头像地址
	OpenID   string `json:"openid"`   // 加密后的瑞雪openid
	Region   string `json:"region"`   // 地区码
	Sex      int8   `json:"sex"`      // 性别
	Age      uint8  `json:"age"`      // 年龄
}

const (
	SignFieldNickname = "nickname"
	SignFieldAvatar   = "avatar"
	SignFieldOpenID   = "openid"
	SignFieldRegion   = "region"
	SignFieldSex      = "sex"
	SignFieldAge      = "age"
)
