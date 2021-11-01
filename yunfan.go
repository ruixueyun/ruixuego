// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"git.jiaxianghudong.com/ruixuesdk/ruixuego/bytepool"
)

var (
	BytePools, _ = bytepool.NewMultiRatedBytePool(4, 10, 1024)
)
