// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import "time"

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
	ReqHeader
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

type CPRoleInfo struct {
	ReqHeader
	RxOpenID  string                 `json:"rx_openid"`           // 瑞雪openid
	RegionTag string                 `json:"region_tag"`          // 区服id
	CPRoleID  string                 `json:"cp_role_id"`          // 角色id
	Extension map[string]interface{} `json:"extension,omitempty"` // 扩展字段:json格式
}

type CPRoleInfoDel struct {
	ReqHeader
	RxOpenID  string `json:"rx_openid"`  // 瑞雪openid
	RegionTag string `json:"region_tag"` // 区服id
	CPRoleID  string `json:"cp_role_id"` // 角色id
}

type CPRoleList struct {
	ReqHeader
	RxOpenID  string `json:"rx_openid"` // 瑞雪openid
	Extension string `json:"extension"` // *:全部、nickname,avatar,sex:指定字段 、空:忽略
}

type CPRoleRes struct {
	RxOpenID   string                 `json:"rx_openid"`           // 瑞雪openid
	RegionTag  string                 `json:"region_tag"`          // 区服id
	CPRoleID   string                 `json:"cp_role_id"`          // 角色id
	ReportTime int64                  `json:"report_time"`         // 上报时间
	Extension  map[string]interface{} `json:"extension,omitempty"` // 扩展字段:json格式
}

type ReqSetCustom struct {
	ReqHeader
	ProductID string
	OpenID    string
	CpUserID  string
	Custom    string
}

type ReqAddRelation struct {
	ReqHeader
	Types        RelationTypes
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
	Remarks      []string
}

type ReqDelRelation struct {
	ReqHeader
	Types        RelationTypes
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
}

type ReqUpdateRelationRemarks struct {
	ReqHeader
	Type         string
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
	Remarks      string
}

type ReqRelationList struct {
	ReqHeader
	Type   string
	OpenID string
	UserID string
}

type ReqHasRelation struct {
	ReqHeader
	Type         string
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
}

type ReqAddFriend struct {
	ReqHeader
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
	Remarks      []string
}

type ReqDelFriend struct {
	ReqHeader
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
}

type ReqUpdateFriendRemarks struct {
	ReqHeader
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
	Remarks      string
}

type ReqFriendList struct {
	ReqHeader
	OpenID string
	UserID string
}

type ReqGetRelationUser struct {
	ReqHeader
	Type         string
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
}

type ReqIsFriend struct {
	ReqHeader
	OpenID       string
	UserID       string
	TargetOpenID string
	TargetUserID string
}

type ReqLBSUpdate struct {
	ReqHeader
	OpenID string
	UserID string
	Types  []string
	Lon    float64
	Lat    float64
}

type ReqLBSDelete struct {
	ReqHeader
	OpenID string
	UserID string
	Types  []string
}

type ReqLBSRadius struct {
	ReqHeader
	OpenID   string
	UserID   string
	Typ      string
	Lon      float64
	Lat      float64
	Radius   float64
	Page     int
	PageSize int
	Count    []int
}

type ReqTrack struct {
	ReqHeader
	Data     []byte
	LogCount int
	Compress bool
}

type ReqSyncTrack struct {
	ReqHeader
	DeviceCode string
	DistinctID string
	Opts       []BigdataOptions
}

type ReqCreateRank struct {
	ReqHeader
	RankID      string
	StartTime   time.Time
	DestroyTime time.Time
}

type ReqCloseRank struct {
	ReqHeader
	RankID string
}

type ReqRankAddScore struct {
	ReqHeader
	RankID   string
	OpenId   string
	CpUserID string
	Score    int64
}

type ReqRankSetScore struct {
	ReqHeader
	RankID   string
	OpenId   string
	CpUserID string
	Score    int64
}

type ReqDeleteRankUser struct {
	ReqHeader
	RankID   string
	OpenId   string
	CpUserID string
}

type ReqQueryUserRank struct {
	ReqHeader
	RankID   string
	OpenId   string
	CpUserID string
}

type ReqGetRankList struct {
	ReqHeader
	RankID string
	Start  int32
	End    int32
}

type ReqGetFriendRankList struct {
	ReqHeader
	RankID   string
	OpenId   string
	CpUserID string
}

type ReqGetRankDetail struct {
	ReqHeader
	RankID string
}

type ReqExtensionGameDisplay struct {
	ReqHeader
	GameID string `json:"game_id"`
}

type ReqTradeOrderStatusByNo struct {
	ReqHeader
	OrderNo string `json:"order_no"`
}

type ReqCheckUserInSiyu struct {
	ReqHeader
	RxOpenID string `json:"rx_open_id"`
	CpUserID string `json:"cp_user_id"`
}
