package ruixuego

type ArgsUserInSiyu struct {
	CPUserID string `json:"cp_user_id"` //
	RxOpenID string `json:"rx_open_id"` // rx_open_id
}

type RespUserInSiyu struct {
	CPUserID   string `json:"cp_user_id"`  //
	RxOpenID   string `json:"rx_open_id"`  // rx_open_id
	BindStatus int64  `json:"bind_status"` // 0:不在 1:存在&未建联 2:存在&已建联
}
