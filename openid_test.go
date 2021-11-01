// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import "testing"

func TestOpenID(t *testing.T) {
	err := Init(&Config{})
	if err != nil {
		panic(err)
	}

	err = AddAESKey("test", []byte("098f6bcd4621d373cade4e832627b4f6"))
	if err != nil {
		panic(err)
	}

	ret, err := EncryptOpenIDData("test", "123abc", "aaabbbccc123456")
	if err != nil {
		panic(err)
	}

	t.Logf("ciphertext => %s", ret)

	data, err := DecryptOpenIDData("test", ret)
	if err != nil {
		panic(err)
	}

	t.Logf("appid: %s, token: %s, openid: %s, ts: %d",
		data.AppID, data.Token, data.OpenID, data.Timestamp)
}
