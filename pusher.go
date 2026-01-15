// Copyright (c) 2022. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

// PusherPushReq 推送请求参数
type PusherPushReq struct {
	PushInfo       PusherPushReqInfo   `json:"push_info"`
	Target         PusherPushReqTarget `json:"target"`
	TargetUserType int32               `json:"target_user_type"` // 推送通道 0：正式通道，1：测试通道(vivo，honor、huawei，苹果 厂商支持)
}

type PusherPushRes struct {
	DeviceTypeMap map[string]string `json:"device_type_map,omitempty"` // 仅测试通道返回，标记当前厂商状态
}

type PusherPushReqInfo struct {
	Title          string `json:"title"`
	Content        string `json:"content"`
	Payload        string `json:"payload"`
	Action         uint8  `json:"action"`
	Classification uint8  `json:"classification"`
}

type PusherPushReqTarget struct {
	Type     uint8    `json:"type"`
	Openid   []string `json:"openid"`
	CPUserID []string `json:"cp_user_id"`
	Tag      string   `json:"tag"`
	Alias    []string `json:"alias"`
}

type ReqPusher struct {
	ReqHeader
	Req       *PusherPushReq
	ProductID string
	ChannelID string
}
