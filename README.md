# 瑞雪 Go 服务端 SDK

## 已实现功能
- 通行证 LoginOpenID 加解密（LoginOpenID 是由瑞雪服务端在登录成功后下发的加密登录 OpenID，CP 方可通过此字段进行安全检查）;
- 自定/好友关系操作（添加、删除，修改备注）；
- 给用户设置社交模块的自定义信息；
- LBS 相关接口；
- 大数据埋点（参考 examples/bigdata/main.go）；


## 使用方法

通过 `go mod` 引入 SDK

```shell
go get ruixuego
```

```go
package main

import (
	"fmt"

	"ruixuego"
)

const (
	testAppID     = "wltestapp"
	testChannelID = "wltestchannel"
	testAppKey    = "a463deade4b15d5ac5398f97cdaeab65"
)

func main() {
	// SDK 初始化
	err := ruixuego.Init(&ruixuego.Config{
		APIDomain: "http://ruixue.weiletest.com",       // 瑞雪云接口域名 URL, 以 http[s]:// 开头
		CPID:      1000000,                             // CPID 由瑞雪团队提供
		CPKey:     "f3c7907d161764daf97fdaaea1a72261",  // CPKey 由瑞雪团队提供
		AppKeys:   map[string]map[string]string{        // CP 方应用配置
			testAppID: {
				testChannelID: testAppKey,
			},
		},
		BigData: &ruixuego.BigDataConfig{ // 要使用大数据埋点功能必须配置此参数
			AutoFlush: true,
		},
	})
	if err != nil {
		panic(err)
	}

	defer func() {
		// 使用大数据埋点功能上传数据后, 必须在程序退出前显式调用 ruixuego.Close()
		// 不然可能导致数据丢失
		fmt.Println("close result:", ruixuego.Close())
	}()

	// 初始化成功后通过以下方式进行接口调用
	// ruixuego.GetDefaultClient().SomeAPI()
}
```