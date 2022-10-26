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
