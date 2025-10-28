package ruixuego

type (
	GreenRequestTask struct {
		Tag int32  `json:"tag,omitempty"` // 图片类型
		URL string `json:"url"`           // 检测图片的地址
	}

	GreenRequest struct {
		Seed         string              `json:"-"` // 用于 Callback 加密
		Interval     int32               `json:"-"` // 截帧频率，GIF图、长图检测专用
		MaxFrames    int32               `json:"-"` // 最大截帧数量，GIF图、长图检测专用，默认值为1
		CPID         uint32              `json:"-"`
		ProductID    string              `json:"-"`
		ChannelID    string              `json:"-"`
		BizType      string              `json:"-"`                  // 标识业务场景
		Scenes       []string            `json:"scenes,omitempty"`   // 指定检测场景 [ "porn"(鉴黄),"terrorism"(暴恐涉政),"ad"(广告),"live"(不良场景),"qrcode"(二维码),"logo"(标志)]
		Tasks        []*GreenRequestTask `json:"tasks,omitempty"`    // 检测任务
		TaskID       []string            `json:"taskids,omitempty"`  // 查询检测结果 taskID 最多 100 个
		CPCallback   string              `json:"callback,omitempty"` // 客户回调地址
		Extend       string              `json:"extend,omitempty"`   // 透传数据
		RiskCallback string              `json:"-"`                  // 本服务的回调地址
	}

	GreenFeedbackRequest struct {
		CPID      uint32            `json:"-"`
		ProductID string            `json:"-"`
		ChannelID string            `json:"-"`
		Results   map[string]string `json:"results"` // 指定检测场景 ["porn","terrorism","ad","live","qrcode","logo"]
		URL       string            `json:"url"`
		TaskID    string            `json:"taskid"`
	}
)

type (

	// GreenUsercaseResultScene 不同 scene 返回结果
	GreenUsercaseResultScene struct {
		Scene      string `json:"scene"`
		Suggestion string `json:"suggestion"`
	}

	// GreenCallbackResultTask 阿里回调的检测结果
	GreenCallbackResultTask struct {
		Code         int                         `json:"code"`
		Msg          string                      `json:"msg,omitempty"`
		TaskID       string                      `json:"taskid"`
		URL          string                      `json:"url,omitempty"`
		SceneResults []*GreenUsercaseResultScene `json:"results,omitempty"`
	}

	// GreenUsercaseResultTask 每个 URL 检测结果
	GreenUsercaseResultTask struct {
		Code         int               `json:"code"`
		Msg          string            `json:"msg,omitempty"`
		TaskID       string            `json:"taskid"`
		URL          string            `json:"url,omitempty"`
		Tag          int32             `json:"tag,omitempty"`
		Result       string            `json:"result,omitempty"`
		SceneResults map[string]string `json:"scene_result,omitempty"`
	}

	// GreenUsercaseResult 内容安全检测同步返回结果
	GreenUsercaseResult struct {
		Results []*GreenUsercaseResultTask `json:"taskresult"`
		Extend  string                     `json:"extend,omitempty"`
	}
)

// SensitiveReq 敏感词检测请求
type SensitiveReq struct {
	Content string `json:"check_words"` // 需要检测的语句
}

// SensitiveResponse 敏感词检测返回
type SensitiveResponse struct {
	Content        string   `json:"content"`         // 返回去敏感词的语句
	SensitiveWords []string `json:"sensitive_words"` // 敏感词
}

type MediaCheckReq struct {
	URLs   []string `json:"urls"` // 图片地址
	Scenes []string `json:"-"`    // 阿里的场景，微信验证通过后还会走一遍阿里验证 [ "porn"(鉴黄),"terrorism"(暴恐涉政),"ad"(广告),"live"(不良场景),"qrcode"(二维码),"logo"(标志)]
}

type MediaResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		URL     string `json:"url"`
		TraceID int    `json:"trace_id"`
	} `json:"data"`
}

type RiskContentImageScanReq struct {
	ReqHeader
	URL string `json:"url"`
}

type RiskContentImageScanResp struct {
	URL        string    `json:"url"`        // 图片地址
	Suggestion string    `json:"suggestion"` // 审核结果 pass: 通过；risky: 风险
	Results    []Results `json:"results"`    // 检测结果详情
}

type Results struct {
	Label      string  `json:"label"`      // 标签
	Confidence float64 `json:"confidence"` // 置信分值，0到100分; 部分标签无置信分
}

type RiskContentTextScanReq struct {
	ReqHeader
	OpenID  string `json:"open_id"`                    // 微信小程序openid
	Scene   string `json:"scene" binding:"required"`   // 场景 nick_name: 昵称；private_chat: 私聊；public_chat: 公聊评论
	Content string `json:"content" binding:"required"` // 文本内容
}

type RiskContentTextScanResp struct {
	Suggestion string   `json:"suggestion"` // 审核结果 pass: 通过；risky: 风险
	Words      []string `json:"words"`      // 敏感词
	Content    string   `json:"content"`    // 替换后的文本内容
}

// RealAuthReq 实名请求
type RealAuthReq struct {
	ReqHeader
	CPID      uint32 `json:"cpid"`       // cpid
	ProductID string `json:"product_id"` // 产品id
	IDCard    string `json:"id_card"`    // 身份证
	RealName  string `json:"real_name"`  // 姓名
}

// RealAuthResponse 实名结果
type RealAuthResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Verified int32 `json:"verified,omitempty"`
	}
}
