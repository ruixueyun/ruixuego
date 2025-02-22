// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

const (
	RelationTypeFriend = "friend"
)

// RelationTypes 自定义用户关系类型
//
//	map[自定义类型]是否为双向关系
type RelationTypes map[string]bool

type response struct {
	Msg  string      `json:"msg"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type argCustom struct {
	ProductID string `json:"productid"`
	OpenID    string `json:"openid"`
	CPUserID  string `json:"cp_user_id"`
	Custom    string `json:"custom"`
}

type argRelation struct {
	Types          RelationTypes `json:"types,omitempty"` // map[CP自定类型]是否为双向关系
	Type           string        `json:"type"`
	Target         string        `json:"target"`
	TargetCPUserID string        `json:"target_cp_user_id"`
	TargetRemarks  string        `json:"target_remarks,omitempty"`
	OpenID         string        `json:"openid,omitempty"`
	CPUserID       string        `json:"cp_user_id"`
	UserRemarks    string        `json:"user_remarks,omitempty"`
}

type argLocation struct {
	OpenID    string   `json:"openid,omitempty"`
	CPUserID  string   `json:"cp_user_id"`
	Type      string   `json:"type,omitempty"`
	Types     []string `json:"types,omitempty"`
	Longitude float64  `json:"lon,omitempty"`
	Latitude  float64  `json:"lat,omitempty"`
	Radius    float64  `json:"radius,omitempty"`
	Count     int      `json:"count,omitempty"`
	Page      int      `json:"page,omitempty"`
	PageSize  int      `json:"page_size,omitempty"`
}

//type RelationUser struct {
//	OpenID   string  `json:"OpenID,omitempty"`
//	NickName string  `json:"NickName,omitempty"`
//	Birthday string  `json:"Birthday,omitempty"`
//	Remarks  string  `json:"Remarks,omitempty"`
//	Avatar   string  `json:"Avatar,omitempty"`
//	Custom   string  `json:"Custom,omitempty"`
//	Dist     float64 `json:"Dist,omitempty"`
//	Time     int64   `json:"Time,omitempty"`
//	Score    int64   `json:"Score,omitempty"`
//	CPID     uint32  `json:"CPID,omitempty"`
//	Sex      int32   `json:"Sex,omitempty"`
//}

type rankAPIArg struct {
	RankID      string `json:"rank_id"`
	Score       int64  `json:"score,omitempty"`
	OpenID      string `json:"open_id,omitempty"`
	CPUserID    string `json:"cp_user_id"`
	StartTime   string `json:"start,omitempty"`
	DestroyTime string `json:"destroy,omitempty"`
	StartRank   int32  `json:"start_rank,omitempty"`
	EndRank     int32  `json:"end_rank,omitempty"`
}

type RelationUser struct {
	OpenID   string  `json:"openid,omitempty"`
	CPUserID string  `json:"cp_user_id"`
	NickName string  `json:"nickname,omitempty"`
	Birthday string  `json:"birthday,omitempty"`
	Remarks  string  `json:"remarks,omitempty"`
	Avatar   string  `json:"avatar,omitempty"`
	Custom   string  `json:"custom,omitempty"`
	Dist     float64 `json:"dist,omitempty"`
	Time     int64   `json:"time,omitempty"`
	Score    int64   `json:"score,omitempty"`
	CPID     uint32  `json:"cpid,omitempty"`
	Sex      int32   `json:"sex,omitempty"`
}

// RankMember 玩家排行对象
type RankMember struct {
	UserName string        `json:"-"`
	Score    int64         `json:"score"`
	Rank     int64         `json:"rank"` // 玩家排名
	UserInfo *RelationUser `json:"user,omitempty"`
}

// ReportCustomAction 投放归因上报自定义action
type ReportCustomAction struct {
	OpenID string `json:"open_id"`
	Action string `json:"action"` // 上报行为
}

type RespRankDetail struct {
	RankID       string `json:"rank_id"`
	Capacity     int32  `json:"capacity"`      // 排行榜容量
	ContinueTime int32  `json:"continue_hour"` // 排行榜持续时间 单位（小时）
	Flag         int64  `json:"flag"`          // CP 自定义 Flag
	StartTime    string `json:"start_time"`    // 排行榜开启时间
	DestroyTime  string `json:"destroy_time"`  // 排行榜销毁时间
}

type RespAllRankID struct {
	RankIDList []string `json:"list"`
}
