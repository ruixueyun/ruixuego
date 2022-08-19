// Copyright (c) 2022. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

const (
	dateTimeFormat = "2006-01-02 15:04:05.000"
)

// 预置 Key 定义
const (
	PresetKeyPlatformID   = "$platformid"
	PresetKeyCPID         = "$cpid"
	PresetKeyAppID        = "$appid"
	PresetKeyChannelID    = "$channelid"
	PresetKeySubChannelID = "$subchannelid"
	PresetKeyUUID         = "$uuid"
	PresetKeyTime         = "$time"
	PresetKeyIP           = "$ip"
)

const (
	typeTrack   = "track"
	typeUser    = "user"
	UserSet     = "user_set"
	UserSetOnce = "user_setonce"
	UserAdd     = "user_add"
	UserMin     = "user_min"
	UserMax     = "user_max"
)

type BigDataLog struct {
	Type         string                 `json:"type"`
	Time         string                 `json:"time"`
	DistinctID   string                 `json:"distinct_id"`
	Devicecode   string                 `json:"devicecode"`
	Event        string                 `json:"event"`
	UUID         string                 `json:"uuid"`
	IP           string                 `json:"ip,omitempty"`
	Properties   map[string]interface{} `json:"properties"`
	AppID        string                 `json:"app_id,omitempty"`
	ChannelID    string                 `json:"channel_id,omitempty"`
	SubChannelID string                 `json:"sub_channel_id,omitempty"`
	CPID         uint32                 `json:"cpid"`
	PlatformID   int32                  `json:"platform_id"`
}
