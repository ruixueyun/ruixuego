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

type argCustom struct {
	AppID  string `json:"appid"`
	OpenID string `json:"openid"`
	Custom string `json:"custom"`
}

type argRelation struct {
	Types         RelationTypes `json:"types,omitempty"` // map[CP自定类型]是否为双向关系
	Type          string        `json:"type"`
	Target        string        `json:"target"`
	TargetRemarks string        `json:"target_remarks,omitempty"`
	OpenID        string        `json:"openid,omitempty"`
	UserRemarks   string        `json:"user_remarks,omitempty"`
}

type argLocation struct {
	OpenID    string   `json:"openid,omitempty"`
	Type      string   `json:"type,omitempty"`
	Types     []string `json:"types,omitempty"`
	Longitude float64  `json:"lon,omitempty"`
	Latitude  float64  `json:"lat,omitempty"`
	Radius    float64  `json:"radius,omitempty"`
	Count     int      `json:"count,omitempty"`
	Page      int      `json:"page,omitempty"`
	PageSize  int      `json:"page_size,omitempty"`
}

type RelationUser struct {
	OpenID   string  `json:"OpenID,omitempty"`
	NickName string  `json:"NickName,omitempty"`
	Birthday string  `json:"Birthday,omitempty"`
	Remarks  string  `json:"Remarks,omitempty"`
	Avatar   string  `json:"Avatar,omitempty"`
	Custom   string  `json:"Custom,omitempty"`
	Dist     float64 `json:"Dist,omitempty"`
	Time     int64   `json:"Time,omitempty"`
	Score    int64   `json:"Score,omitempty"`
	CPID     uint32  `json:"CPID,omitempty"`
	Sex      int32   `json:"Sex,omitempty"`
}

type UserInfo struct {
	AppID    string `json:"appid,omitempty"`
	OpenID   string `json:"openid,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Region   string `json:"region,omitempty"`   // Format: 220101
	Birthday string `json:"birthday,omitempty"` // Format: 2006-01-02
	Sex      string `json:"sex,omitempty"`      // 0: female, 1: male
}
