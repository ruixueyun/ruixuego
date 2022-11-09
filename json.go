// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"ruixuego/bufferpool"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// MarshalJSONEscapeHTML 编码 JSON, 不编码 HTML 标记字符
func MarshalJSONEscapeHTML(v interface{}) ([]byte, error) {
	buf := bufferpool.Get()
	defer bufferpool.Put(buf)
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	b := make([]byte, buf.Len())
	copy(b, buf.Bytes())
	return b, nil
}

// MarshalJSON 编码 JSON
func MarshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// UnmarshalJSON 解码 JSON
func UnmarshalJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
