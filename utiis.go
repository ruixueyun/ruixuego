// Copyright (c) 2022. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"regexp"
	"strconv"
	"strings"
)

// Atoi32 string to int32
func Atoi32(s string, d ...int32) int32 {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}

	return int32(i)
}

// Atoi64 string to int64
func Atoi64(s string, d ...int64) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		if len(d) > 0 {
			return d[0]
		} else {
			return 0
		}
	}

	return i
}

// I64toa int64 转字符串
func I64toa(i int64) string {
	return strconv.FormatInt(i, 10)
}

// I32toa int32 转字符串
func I32toa(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

const (
	imsConvTypeJoiner = "$"
	imsMsgIDJoiner    = "."
	imsUserIDJoiner   = ":"
)

var (
	// imsRegexpConvData ConvData 验证器
	imsRegexpConvData = regexp.MustCompile(`^[a-zA-Z\d_-]{1,32}$`)
)

// IMSVerifyUserID 验证 UserID 是否合规
func IMSVerifyUserID(userID string) bool {
	if userID == "" {
		return false
	}
	return imsRegexpConvData.MatchString(userID)
}

// IMSParseMsgID 解析 MsgID，convType 为空表示 msgID 不合规
func IMSParseMsgID(msgID string) (msgSeqID int64, convType int32, conversationID string, ok bool) {
	if msgID == "" {
		return 0, 0, "", false
	}

	var conversationSeqIDStr string
	conversationSeqIDStr, conversationID, ok = strings.Cut(msgID, imsMsgIDJoiner)
	if !ok {
		return 0, 0, "", false
	}

	msgSeqID = Atoi64(conversationSeqIDStr)
	if msgSeqID == 0 {
		return 0, 0, "", false
	}

	convType, _, ok = IMSParseConversationID(conversationID)
	if !ok {
		return 0, 0, "", false
	}
	return
}

// imsGetConversationID 生成 ConversationID
func imsGetConversationID(typ int32, data string) string {
	return imsConvTypeJoiner + I32toa(typ) + imsConvTypeJoiner + data
}

// IMSParseConversationID 解析 ConversationID
func IMSParseConversationID(conversationID string) (convType int32, ConversationIDData string, ok bool) {
	if conversationID == "" || conversationID[0] != '$' {
		return
	}

	var convTypeStr string
	convTypeStr, ConversationIDData, ok = strings.Cut(conversationID[1:], imsConvTypeJoiner)
	if !ok {
		return 0, "", false
	}

	convType = Atoi32(convTypeStr)
	if _, ok := imsSupportedConversationTypes[convType]; !ok {
		return 0, "", false
	}

	if convType == IMSConvTypeSingle { // 单聊
		_, _, ok = IMSParseSingleConversationIDData(ConversationIDData)
		if !ok {
			return 0, "", false
		}
	} else if !imsRegexpConvData.MatchString(ConversationIDData) {
		return 0, "", false
	}

	return convType, ConversationIDData, ok
}

// IMSGetSingleConversationID 生成单聊会话 ID
func IMSGetSingleConversationID(sender, receiver string) string {
	if !IMSVerifyUserID(sender) || !IMSVerifyUserID(receiver) {
		return ""
	}
	if sender < receiver {
		return imsGetConversationID(IMSConvTypeSingle, sender+imsUserIDJoiner+receiver)
	}
	return imsGetConversationID(IMSConvTypeSingle, receiver+imsUserIDJoiner+sender)
}

// IMSGetGroupConversationID 生成群聊会话 ID
func IMSGetGroupConversationID(groupID string) string {
	if !IMSVerifyUserID(groupID) {
		return ""
	}
	return imsGetConversationID(IMSConvTypeGroup, groupID)
}

// IMSGetCustomSingleConversationID 生成自定义单聊会话 ID
func IMSGetCustomSingleConversationID(customID string) string {
	if !IMSVerifyUserID(customID) {
		return ""
	}
	return imsGetConversationID(IMSConvTypeCustomSingle, customID)
}

// IMSParseSingleConversationIDData 解析单聊 ConversationIDData
// 小的 UserID 在前
func IMSParseSingleConversationIDData(data string) (minUserID, maxUserID string, ok bool) {
	minUserID, maxUserID, ok = strings.Cut(data, imsUserIDJoiner)
	if !ok {
		return "", "", false
	}

	if !IMSVerifyUserID(minUserID) || !IMSVerifyUserID(maxUserID) {
		return "", "", false
	}

	// 返回结果，小的 UserID 在前
	if minUserID < maxUserID {
		return minUserID, maxUserID, ok
	}
	return maxUserID, minUserID, ok
}
