// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
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
	AppID     string `json:"appid"`
	Token     string `json:"token"`
	OpenID    string `json:"openid"`
	Timestamp int64  `json:"ts"`
}

func (data *OpenIDData) Release() {
	*data = OpenIDData{}
	poolOpenIDData.Put(data)
}

// AddAESKey 添加预置密钥
func AddAESKey(appID string, key []byte) error {
	k, err := NewAESData(key)
	if err != nil {
		return err
	}
	appKeys.Store(appID, k)
	return nil
}

// DelAESKey 删除预置密钥
func DelAESKey(appID string) {
	appKeys.Delete(appID)
}

// EncryptOpenIDData 加密 OpenID 数据, 获取密文
func EncryptOpenIDData(appID, token, openID string) (string, error) {
	k, ok := appKeys.Load(appID)
	if !ok {
		return "", ErrAppKeyNotExistx
	}
	return EncryptOpenIDDataWithKey(k.(*AESData), appID, token, openID)
}

// EncryptOpenIDDataWithKey 加密 OpenID 数据, 获取密文
func EncryptOpenIDDataWithKey(
	aesData *AESData,
	appID, token, openID string) (ret string, err error) {

	data := poolOpenIDData.Get().(*OpenIDData)
	defer data.Release()
	data.AppID = appID
	data.Token = token
	data.OpenID = openID
	data.Timestamp = time.Now().UnixMilli()
	b, err := MarshalJSON(data)
	if err != nil {
		return ret, err
	}
	return AesEncryptBase64String(BytesToString(b), aesData), nil
}

// DecryptOpenIDData 解密 OpenIDData 密文字符串
func DecryptOpenIDData(appID, openIDCipherText string) (*OpenIDData, error) {
	k, ok := appKeys.Load(appID)
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
