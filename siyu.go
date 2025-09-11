package ruixuego

type ArgsUserInSiyu struct {
	CPUserID string `json:"cp_user_id"` // 游戏 ID
	RxOpenID string `json:"rx_open_id"` // 瑞雪 openid
}

type RespUserInSiyu struct {
	CPUserID     string `json:"cp_user_id"`    // 游戏 ID
	RxOpenID     string `json:"rx_open_id"`    // 瑞雪 openid
	BindStatus   int64  `json:"bind_status"`   // 0:不在 1:存在&未建联 2:存在&已建联
	MinigamePath string `json:"minigame_path"` // 私域小游戏路径
}
