package ruixuego

type UpdateCPUserIDRequest struct {
	OpenID   string `json:"open_id"`    // 瑞雪openid
	CPUserID string `json:"cp_user_id"` // cp侧user_id
}
