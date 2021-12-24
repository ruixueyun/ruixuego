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
	TraceID   string `json:"traceid"`
	AppID     string `json:"appid"`
	ChannelID string `json:"channelid"`
	Method    string `json:"method"`
	OpenID    string `json:"openid"`
	Ext       string `json:"ext,omitempty"`
	Timestamp int64  `json:"ts"`
}

func (data *OpenIDData) Release() {
	*data = OpenIDData{}
	poolOpenIDData.Put(data)
}

func getKey(appID, channelID string) string {
	return appID + "_" + channelID
}

// AddAESKey 添加预置密钥
func AddAESKey(appID, channelID string, key []byte) error {
	k, err := NewAESData(key)
	if err != nil {
		return err
	}
	appKeys.Store(getKey(appID, channelID), k)
	return nil
}

// DelAESKey 删除预置密钥
func DelAESKey(appID, channelID string) {
	appKeys.Delete(getKey(appID, channelID))
}

// EncryptOpenIDData 加密 OpenID 数据, 获取密文
func EncryptOpenIDData(
	traceID, appID, channelID, method, openID, ext string) (string, error) {

	k, ok := appKeys.Load(getKey(appID, channelID))
	if !ok {
		return "", ErrAppKeyNotExistx
	}
	return EncryptOpenIDDataWithKey(
		k.(*AESData), traceID, appID, channelID, method, openID, ext)
}

// EncryptOpenIDDataWithKey 加密 OpenID 数据, 获取密文
func EncryptOpenIDDataWithKey(
	aesData *AESData,
	traceID, appID, channelID, method, openID, ext string) (ret string, err error) {

	data := poolOpenIDData.Get().(*OpenIDData)
	defer data.Release()
	data.TraceID = traceID
	data.AppID = appID
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
func DecryptOpenIDData(appID, channelID, openIDCipherText string) (*OpenIDData, error) {
	k, ok := appKeys.Load(getKey(appID, channelID))
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
