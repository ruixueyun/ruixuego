package ruixuego

type ExtensionExchangeReq struct {
	ReqHeader
	CdKey    string `json:"cdkey"`
	OpenID   string `json:"open_id"`
	CpUserID string `json:"cp_user_id"`
}

type ExtensionProp struct {
	PropID   string `json:"prop_id"`   // 道具ID
	PropName string `json:"prop_name"` // 道具名称
	PropNum  string `json:"prop_num"`  // 道具数量
}

type GameDisplayWelfareCodeInfo struct {
	PromoterID      string `json:"promoter_id"`       // 主播id
	PromoID         string `json:"promo_id"`          // 福利码ID
	RefreshPeriod   int    `json:"refresh_period"`    // 刷新周期 单位: 分钟
	PromoName       string `json:"promo_name"`        // 福利码名称
	GiftName        string `json:"gift_name"`         // 礼包名称
	PromoValidStart string `json:"promo_valid_start"` // 福利码有效开始时间
	PromoValidEnd   string `json:"promo_valid_end"`   // 福利码有效结束时间
	PromoCode       string `json:"promo_code"`
}

// GameDisplayWelfareCodeInfoExp 刷新过期时间
type GameDisplayWelfareCodeInfoExp struct {
	*GameDisplayWelfareCodeInfo
	RefreshPeriodExp int64 `json:"refresh_period_exp"` // 刷新周期过期时间
	Polling          int   `json:"polling"`            // 客户端轮询时间 秒
}
