// Copyright (c) 2022. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

// IMSLoginReq 登录请求参数
type IMSLoginReq struct {
	UserID     string `json:"user_id"`     // 用户在某 CP 下的唯一标识符
	ProductID  string `json:"product_id"`  // 产品 ID
	ChannelID  string `json:"channel_id"`  // 渠道 ID
	DeviceCode string `json:"device_code"` // 设备码
	CPID       uint32 `json:"cpid"`        // CPID
	ClientType int32  `json:"client_type"` // 客户端类型
}

// IMSLoginResp 登录响应参数
type IMSLoginResp struct {
	AccessToken        string `json:"access_token,omitempty"`         // Entry 连接令牌
	AccessTokenExpire  int64  `json:"access_token_expire,omitempty"`  // 连接另外过期时间，毫秒
	RefreshToken       string `json:"refresh_token,omitempty"`        // Entry 刷新令牌
	RefreshTokenExpire int64  `json:"refresh_token_expire,omitempty"` // 刷新令牌过期时间，毫秒
	AESKey             string `json:"aeskey,omitempty"`               // AES 密钥，返回格式为 16 进制 64 位字符串，需解析成 32 Byte 密钥使用
}

// IMSMessage 聊天信息定义
type IMSMessage struct {
	MilliTS        int64             `json:"milli_ts,omitempty"`
	Attr           int64             `json:"attr,omitempty"`
	Option         int64             `json:"option,omitempty"`
	Status         int32             `json:"status,omitempty"`
	ClientType     int32             `json:"client_type,omitempty"`
	Type           int32             `json:"type,omitempty"`
	SubType        int32             `json:"sub_type,omitempty"`
	InboxID        int64             `json:"inbox_id,omitempty"`
	UUID           string            `json:"uuid,omitempty"`
	MsgID          string            `json:"msg_id,omitempty"`
	Sender         string            `json:"sender,omitempty"`
	Receivers      []string          `json:"receivers,omitempty"`
	UnreadCount    int32             `json:"unread_count,omitempty"`
	ReceiverNum    int32             `json:"receiver_num,omitempty"`
	ConversationID string            `json:"conversation_id,omitempty"`
	ConvType       int32             `json:"conv_type,omitempty"`
	CPID           uint32            `json:"cpid,omitempty"`
	ProductID      string            `json:"product_id,omitempty"`
	ChannelID      string            `json:"channel_id,omitempty"`
	Content        string            `json:"content,omitempty"`
	Ext            map[string]string `json:"ext,omitempty"`
	IMSExt         map[string]string `json:"ims_ext,omitempty"`
}

// IMSMessageAck 聊天消息确认响应
type IMSMessageAck struct {
	MsgID   string `json:"msg_id,omitempty"`
	UUID    string `json:"uuid,omitempty"`
	InboxID int64  `json:"inbox_id,omitempty"`
	MilliTS int64  `json:"milli_ts,omitempty"`
}

// IMSHistoryReq 历史聊天记录请求参数
type IMSHistoryReq struct {
	ConversationID string `json:"conversation_id,omitempty"`
	StartMsgID     string `json:"start_msg_id,omitempty"`
	FetchCount     int32  `json:"fetch_count,omitempty"`
}

// IMSHistoryResp 历史聊天记录响应参数
type IMSHistoryResp struct {
	Messages []*IMSMessage `json:"chat_message,omitempty"`
	Count    int32         `json:"count,omitempty"`
	Done     bool          `json:"done,omitempty"`
}

// IMSCreateConvReq 创建会话请求参数
type IMSCreateConvReq struct {
	ConversationID string            `json:"conversation_id"`   // 会话 ID
	Option         int64             `json:"option,omitempty"`  // 会话选项
	ConvType       int32             `json:"conv_type"`         // 会话类型
	Creator        string            `json:"creator"`           // 会话创建者，即单聊的发起者或群的创建者
	Members        []*MemberInfo     `json:"members,omitempty"` // 除创建者外该会话的其他参与者
	Ext            map[string]string `json:"ext,omitempty"`     // 该会话 CP 自定义扩展信息
	IMSExt         map[string]string `json:"ims_ext,omitempty"` // 该会话 IMS 服务的扩展信息，不可更改
}

// MemberInfo 会话成员属性
type MemberInfo struct {
	UserID string            `json:"user_id"`
	Option int64             `json:"option,omitempty"`
	Ext    map[string]string `json:"ext,omitempty"`
	IMSExt map[string]string `json:"ims_ext,omitempty"`
}

// IMSUpdateConvReq 更新会话信息请求参数
type IMSUpdateConvReq struct {
	ConversationID string            `json:"conversation_id"`   // 会话 ID
	Option         int64             `json:"option,omitempty"`  // 会话选项
	Creator        string            `json:"creator,omitempty"` // 会话创建者，即单聊的发起者或群的创建者
	Ext            map[string]string `json:"ext,omitempty"`     // 该会话 CP 自定义扩展信息
	IMSExt         map[string]string `json:"ims_ext,omitempty"` // 该会话 IMS 服务的扩展信息，不可更改
}

// IMSConvDeleteReq 删除会话请求参数
type IMSConvDeleteReq struct {
	ConversationID string `json:"conversation_id"` // 会话 ID
	UserID         string `json:"user_id"`         // 操作人用户 ID
}

// IMSGetConversationReq 获取会话信息请求参数
type IMSGetConversationReq struct {
	ConversationID string `json:"conversation_id"` // 会话 ID
	UserID         string `json:"user_id"`         // 操作人用户 ID
}

// IMSConversation 会话信息
type IMSConversation struct {
	ConversationID string            `json:"conversation_id"`
	Attr           int64             `json:"attr,omitempty"`
	Option         int64             `json:"option,omitempty"`
	UserAttr       int64             `json:"user_attr,omitempty"`
	UserOption     int64             `json:"user_option,omitempty"`
	CreateMilliTS  int64             `json:"create_milli_ts"`
	UpdateMilliTS  int64             `json:"update_milli_ts,omitempty"`
	JoinMilliTS    int64             `json:"join_milli_ts"`
	Creator        string            `json:"creator"`
	UserExt        map[string]string `json:"user_ext,omitempty"` // 用户扩展信息
	Ext            map[string]string `json:"ext,omitempty"`      // 会话扩展信息
	Members        []string          `json:"members,omitempty"`  // []UserID
	Type           int8              `json:"type,omitempty"`     // 会话类型
	Status         int8              `json:"status,omitempty"`
}

// IMSJoinConversationReq 加入会话请求参数
type IMSJoinConversationReq struct {
	ConversationID string            `json:"conversation_id,omitempty"` // 会话 ID
	UserID         string            `json:"user_id"`                   // 用户 ID
	Option         int64             `json:"option,omitempty"`          // 用户在会话中的选项参数
	Ext            map[string]string `json:"ext,omitempty"`             // 用户扩展信息
	ClientType     int32             `json:"client_type,omitempty"`     // 用户客户端类型
}

// IMSLeaveConversationReq 离开会话请求参数
type IMSLeaveConversationReq struct {
	ConversationIDs string `json:"conversation_ids,omitempty"`
	UserID          string `json:"user_id"`
	ClientType      int32  `json:"client_type,omitempty"` // 用户客户端类型
}

// IMSUpdateConvUserDataReq 更新用户在会话中的属性请求参数
type IMSUpdateConvUserDataReq struct {
	ConversationID string            `json:"conversation_id,omitempty"` // 会话 ID
	UserID         string            `json:"user_id"`                   // 用户 ID
	Option         int64             `json:"option,omitempty"`          // 用户在会话中的选项参数
	Ext            map[string]string `json:"ext,omitempty"`             // 用户扩展信息
	ClientType     int32             `json:"client_type,omitempty"`     // 用户客户端类型
}

// IMSConversationUserListReq 获取用户会话列表请求参数
type IMSConversationUserListReq struct {
	UserID string `json:"user_id"` // 用户 ID
}
