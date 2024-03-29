// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	jsonen "encoding/json"
	"sync"
	"time"
)

var (
	poolOpenIDData = &sync.Pool{
		New: func() interface{} {
			return &OpenIDData{}
		},
	}

	appKeys = &sync.Map{}
)

// OpenIDData OpenID 加密 OpenID 数据定义
type OpenIDData struct {
	TraceID       string         `json:"traceid,omitempty"`
	ProducdID     string         `json:"productid,omitempty"`
	ChannelID     string         `json:"channelid,omitempty"`
	Method        string         `json:"method,omitempty"`
	OpenID        string         `json:"openid,omitempty"`
	Ext           string         `json:"ext,omitempty"`
	TokenID       string         `json:"tokenid,omitempty"`
	Timestamp     int64          `json:"ts,omitempty"`
	AccountID     int64          `json:"accountid,omitempty"`
	AccountDetail *AccountDetail `json:"accountdetail,omitempty"`
}

type AccountDetail struct {
	CredentialType   string               `json:"credential_type"`
	CredentialParams jsonen.RawMessage    `json:"credential_params"`
	Identifier       *RXAccountIdentifier `json:"identifier"`
	Info             *RXAccountInfo       `json:"info"`
}

type RXAccountIdentifier struct {
	ID        int64  `json:"id"`
	OpenID    string `json:"openid"`
	RXVersion int    `json:"rxversion"`
}

type RXAccountInfo struct {
	// 密码编码类型
	PasswordType int32 `json:"password_type"`

	Gender int8 `json:"gender"`

	RegPlatformType int32 `json:"reg_platform_type"`

	// 账号注册时间戳，秒
	RegisteredAtTime int64 `json:"registered_at_time"`

	UserID int64 `json:"userid"`

	Password string `json:"password"`
	NickName string `json:"nick_name"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`

	RealName string `json:"real_name"`
	IDCard   string `json:"idcard"`
}

func (data *OpenIDData) Release() {
	*data = OpenIDData{}
	poolOpenIDData.Put(data)
}

func getKey(productID, channelID string) string {
	return productID + "_" + channelID
}

// AddAESKey 添加预置密钥
func AddAESKey(productID, channelID string, key []byte) error {
	k, err := NewAESData(key)
	if err != nil {
		return err
	}
	appKeys.Store(getKey(productID, channelID), k)
	return nil
}

// DelAESKey 删除预置密钥
func DelAESKey(productID, channelID string) {
	appKeys.Delete(getKey(productID, channelID))
}

// EncryptOpenIDData 加密 OpenID 数据, 获取密文
func EncryptOpenIDData(
	traceID, productID, channelID, method, openID, ext string) (string, error) {

	k, ok := appKeys.Load(getKey(productID, channelID))
	if !ok {
		return "", ErrAppKeyNotExistx
	}
	return EncryptOpenIDDataWithKey(
		k.(*AESData), traceID, productID, channelID, method, openID, ext)
}

// EncryptOpenIDDataWithKey 加密 OpenID 数据, 获取密文
func EncryptOpenIDDataWithKey(
	aesData *AESData,
	traceID, productID, channelID, method, openID, ext string) (ret string, err error) {

	data := poolOpenIDData.Get().(*OpenIDData)
	defer data.Release()
	data.TraceID = traceID
	data.ProducdID = productID
	data.ChannelID = channelID
	data.Method = method
	data.OpenID = openID
	data.Ext = ext
	data.Timestamp = time.Now().UnixMilli()
	b, err := MarshalJSON(data)
	if err != nil {
		return ret, err
	}
	return AesEncryptBase64String(BytesToString(b), aesData), nil
}

// DecryptOpenIDData 解密 OpenIDData 密文字符串
func DecryptOpenIDData(productID, channelID, openIDCipherText string) (*OpenIDData, error) {
	k, ok := appKeys.Load(getKey(productID, channelID))
	if !ok {
		return nil, ErrAppKeyNotExistx
	}
	return DecryptOpenIDDataWithKey(k.(*AESData), openIDCipherText)
}

// DecryptOpenIDDataWithKey 解密 OpenIDData 密文字符串
func DecryptOpenIDDataWithKey(
	aesData *AESData, openIDCipherText string) (*OpenIDData, error) {

	s := AesDecryptBase64String(openIDCipherText, aesData)
	if s == "" {
		return nil, ErrInvalidOpenID
	}
	data := poolOpenIDData.Get().(*OpenIDData)
	err := UnmarshalJSON(StringToBytes(s), data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GenerateVirtualLoginData 生成用于虚拟登录瑞雪的登录凭证
// 用于未接入通行证的业务又想实用其他接口的情况
// userID 参数为 CP 自己用户的唯一标识符, 一般情况下推荐使用用户 ID
func GenerateVirtualLoginData(
	traceID, productID, channelID, userID string) (ret string, err error) {
	return EncryptOpenIDData(traceID, productID, channelID, "virtual", "", userID)
}

// GenerateVirtualLoginDataWithKey 生成用于虚拟登录瑞雪的登录凭证
func GenerateVirtualLoginDataWithKey(
	aesData *AESData,
	traceID, productID, channelID, userID string) (ret string, err error) {
	return EncryptOpenIDDataWithKey(aesData, traceID, productID, channelID, "virtual", "", userID)
}
