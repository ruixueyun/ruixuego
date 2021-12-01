// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

const (
	RelationTypeFriend = "friend"
)

// RelationTypes 自定义用户关系类型
// 		map[自定义类型]是否为双向关系
type RelationTypes map[string]bool

type response struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type argRelation struct {
	Types         RelationTypes `json:"types,omitempty"` // map[CP自定类型]是否为双向关系
	Type          string        `json:"type"`
	Target        string        `json:"target"`
	TargetRemarks string        `json:"target_remarks,omitempty"`
	OpenID        string        `json:"openid,omitempty"`
	UserRemarks   string        `json:"user_remarks,omitempty"`
}

type argCustom struct {
	AppID  string `json:"appid"`
	OpenID string `json:"openid"`
	Custom string `json:"custom"`
}
