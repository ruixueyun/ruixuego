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
	typeTrack = "track"
)

type BigDataLog struct {
	Type         string                 `json:"type"`
	Time         string                 `json:"time"`
	DistinctID   string                 `json:"distinct_id,omitempty"`
	Event        string                 `json:"event,omitempty"`
	UUID         string                 `json:"uuid,omitempty"`
	IP           string                 `json:"ip,omitempty"`
	Properties   map[string]interface{} `json:"properties"`
	AppID        string                 `json:"appid,omitempty"`
	ChannelID    string                 `json:"channelid,omitempty"`
	SubChannelID string                 `json:"subchannelid,omitempty"`
	CPID         uint32                 `json:"cpid"`
	PlatformID   int32                  `json:"platformid"`
	IsLogined    bool                   `json:"islogined"`
}
