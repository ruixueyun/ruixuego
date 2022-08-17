// Copyright (c) 2022. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"crypto/sha1"
	"encoding/hex"
	"hash"
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
