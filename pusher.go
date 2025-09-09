// Copyright (c) 2022. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

// PusherPushReq 推送请求参数
type PusherPushReq struct {
	PushInfo       PusherPushReqInfo   `json:"push_info"`
	Target         PusherPushReqTarget `json:"target"`
	TargetUserType int32               `json:"target_user_type"` // 推送通道 0：正式通道，1：测试通道(vivo，honor、huawei，苹果 厂商支持)
}

type PusherPushRes struct {
	TaskID        uint64                 `json:"task_id"`                   // 任务ID 正式通道异步发送
	DeviceTypeMap map[string]interface{} `json:"device_type_map,omitempty"` // 设备类型映射 测试通道异步发送 直接返回对应设备结果
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
