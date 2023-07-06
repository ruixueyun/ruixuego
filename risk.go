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

type WeiXinMediaCheckReq struct {
	MediaURL  string   `json:"media_url"`  // 要检测的图片或音频的url，支持图片格式包括 jpg , jepg, png, bmp, gif（取首帧），支持的音频格式包括mp3, aac, ac3, wma, flac, vorbis, opus, wav
	MediaType int      `json:"media_type"` // 1:音频;2:图片
	Version   int      `json:"version"`    // 接口版本号，2.0版本为固定值2
	Scene     int      `json:"scene"`      // 场景枚举值（1 资料；2 评论；3 论坛；4 社交日志）
	Openid    string   `json:"openid"`     // 不是瑞雪openid,是微信小程序openid 用户的openid（用户需在近两小时访问过小程序）
	Scenes    []string `json:"-"`          // 阿里的场景，微信验证通过后还会走一遍阿里验证 [ "porn"(鉴黄),"terrorism"(暴恐涉政),"ad"(广告),"live"(不良场景),"qrcode"(二维码),"logo"(标志)]
}

type WeiXinResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Result  struct {
		Suggest string `json:"suggest,omitempty"`
		Label   int    `json:"label,omitempty"`
	} `json:"result,omitempty"`
	Detail []struct {
		Strategy string `json:"strategy,omitempty"`
		ErrCode  int    `json:"errcode"`
		Suggest  string `json:"suggest,omitempty"`
		Label    int    `json:"label"`
		Prob     int    `json:"prob,omitempty"`
		Level    int    `json:"level,omitempty"`
		Keyword  string `json:"keyword,omitempty"`
	} `json:"detail,omitempty"`
	TraceID string `json:"trace_id,omitempty"`
}
