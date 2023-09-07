// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

// Package ruixuego 瑞雪服务端 SDK
package ruixuego

import (
	"fmt"
	"sync"
	"time"
)

type SyncEventAttrsResp struct {
	EventAttrs map[string][]string `json:"public,omitempty"` // 公共属性
	Refresh    int32               `json:"refresh"`          // 刷新时间 单位 ms
	Version    int64               `json:"version"`          // 版本号
}

type PublicAttrHandler func() interface{}

var dataPublicAttr *EventPublicAttr
var eventPublicAttrData *EventAttrsData

type EventAttrsData struct {
	EventAttrs map[string][]string // 公共属性
	Version    int64               // 版本号
	Mutx       *sync.Mutex         //
}

type EventPublicAttr struct {
	Handlers map[string]PublicAttrHandler // key: 属性名 value: 执行的函数
}

func (e *EventPublicAttr) WithPublicAttrHandler(attr string, handler PublicAttrHandler) *EventPublicAttr {
	e.Handlers[attr] = handler
	return e
}

// InitEventPublicAttr 设置公共属性默认值
func InitEventPublicAttr() *EventPublicAttr {
	dataPublicAttr = new(EventPublicAttr)
	dataPublicAttr.Handlers = make(map[string]PublicAttrHandler)
	return dataPublicAttr
}

// GetEventPublicAttrValue 获取公共属性数值
func GetEventPublicAttrValue(attr string) interface{} {
	if dataPublicAttr == nil || dataPublicAttr.Handlers == nil {
		return ""
	}

	handler, ok := dataPublicAttr.Handlers[attr]
	if !ok {
		return ""
	}
	return handler()
}

func SyncEventPublicAttr(c *Client) error {

	eventPublicAttrData = &EventAttrsData{
		Mutx:    &sync.Mutex{},
		Version: 0,
	}

	go func() {
		for {
			resp, err := c.syncEventPublicAttr(eventPublicAttrData.Version)
			if err != nil {
				fmt.Errorf("SyncEventPublicAttr error. err=%+v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			// 版本号相同 不需要同步
			if resp.Version == eventPublicAttrData.Version {
				time.Sleep(time.Duration(resp.Refresh) * time.Millisecond)
				continue
			}

			// 版本号不同需要同步属性
			eventPublicAttrData.Mutx.Lock()
			eventPublicAttrData.Version = resp.Version
			eventPublicAttrData.EventAttrs = resp.EventAttrs
			eventPublicAttrData.Mutx.Unlock()

			time.Sleep(time.Duration(resp.Refresh) * time.Millisecond)
		}
	}()
	return nil
}

func FetchEventPublicAttrs(event string) []string {

	if eventPublicAttrData == nil {
		return nil
	}

	eventPublicAttrData.Mutx.Lock()
	defer eventPublicAttrData.Mutx.Unlock()

	attrs, ok := eventPublicAttrData.EventAttrs[event]
	if ok {
		return attrs
	}

	return nil
}
