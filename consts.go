// Copyright (c) 2022. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

const (
	// IMSConvTypeSingle 单聊
	IMSConvTypeSingle int32 = iota + 1

	// IMSConvTypeGroup 群聊
	IMSConvTypeGroup

	// IMSConvTypeCustomSingle 自定义单聊
	IMSConvTypeCustomSingle

	// IMSConvTypeChannel 频道聊天
	// 该类型主要针对游戏业务，不支持离线消息、已读、撤回等功能
	IMSConvTypeChannel
)

var (
	// imsSupportedConversationTypes 支持的 ConversationType 列表
	imsSupportedConversationTypes = map[int32]struct{}{
		IMSConvTypeSingle:       {},
		IMSConvTypeGroup:        {},
		IMSConvTypeCustomSingle: {},
		IMSConvTypeChannel:      {},
	}
)
