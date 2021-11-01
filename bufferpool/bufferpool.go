// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package bufferpool

import (
	"bytes"
	"sync"
)

// const defaultPoolLenght = 256

var defaultPool = &sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func Get() *bytes.Buffer {
	return defaultPool.Get().(*bytes.Buffer)
}

func Put(b *bytes.Buffer) {
	b.Reset()
	defaultPool.Put(b)
}
