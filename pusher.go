// Copyright (c) 2022. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

// PusherPushReq 推送请求参数
type PusherPushReq struct {
	PushInfo       PusherPushReqInfo   `json:"push_info"`
	Target         PusherPushReqTarget `json:"target"`
	TargetUserType int32               `json:"target_user_type"` // 0：普通消息（默认值），1：测试消息(vivo，honor、huawei 厂商支持)
}

type PusherPushReqInfo struct {
	Title          string `json:"title"`
	Content        string `json:"content"`
	Payload        string `json:"payload"`
	Action         uint8  `json:"action"`
	Classification uint8  `json:"classification"`
}

type PusherPushReqTarget struct {
	Type   uint8    `json:"type"`
	Openid []string `json:"openid"`
	Tag    string   `json:"tag"`
	Alias  []string `json:"alias"`
}
