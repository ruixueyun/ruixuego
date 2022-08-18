// Copyright (c) 2022. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

func NewProducer(c *Client, conf *BigDataConfig) (*Producer, error) {
	conf.done()

	w := newBatchWriter(c, conf)
	err := w.Init()
	if err != nil {
		return nil, err
	}

	return &Producer{
		writer:     w,
		isShutDown: &Bool{},
	}, nil
}

type logWriter interface {
	Init() error
	Write(*BigDataLog) error
	Flush() error
	Close() error
}

type Producer struct {
	writer     logWriter
	wg         sync.WaitGroup
	isShutDown *Bool
}

// Track 大数据埋点事件上报
// 		devicecode 设备码
// 		distinctID 用户标识, 通常为瑞雪 OpenID
// 		event 事件名
//		properties 自定义事件属性
func (p *Producer) Track(devicecode, distinctID, event string, properties map[string]interface{}) error {
	return p.track(typeTrack, devicecode, distinctID, event, properties)
}

// UserTrack 大数据埋点用户上报
// 		devicecode 设备码
// 		distinctID 用户标识, 通常为瑞雪 OpenID
// 		updateType user_setonce,user_set
//		properties 自定义事件属性
func (p *Producer) UserTrack(devicecode, distinctID, updateType string, properties map[string]interface{}) error {
	return p.track(typeUser, devicecode, distinctID, updateType, properties)
}

func (p *Producer) track(eventType, devicecode, distinctID, event string, properties map[string]interface{}) error {
	if devicecode == "" {
		return ErrInvalidDevicecode
	}
	if event == "" {
		return ErrInvalidEvent
	}
	if p.isShutDown.Load() {
		return errProducerShutdown
	}
	p.wg.Add(1)
	defer p.wg.Done()

	cpID := extractCPID(properties)
	if cpID == 0 {
		return ErrInvalidCPID
	}
	appID, channelID, subChannelID, platformID :=
		extractStringProperty(properties, PresetKeyAppID),
		extractStringProperty(properties, PresetKeyChannelID),
		extractStringProperty(properties, PresetKeySubChannelID),
		extractInt32(properties, PresetKeyPlatformID)

	uuidStr, timeStr, ipStr :=
		extractUUID(properties),
		extractTime(properties),
		extractStringProperty(properties, PresetKeyIP)

	logData := &BigDataLog{
		Type:         eventType,
		Time:         timeStr,
		DistinctID:   distinctID,
		Devicecode:   devicecode,
		Event:        event,
		UUID:         uuidStr,
		IP:           ipStr,
		Properties:   properties,
		PlatformID:   platformID,
		AppID:        appID,
		ChannelID:    channelID,
		SubChannelID: subChannelID,
		CPID:         cpID,
	}

	return p.writer.Write(logData)
}

// Close 服务停止前必须显式调用该方法, 不然可能造成数据丢失
func (p *Producer) Close() error {
	p.isShutDown.Store(true)
	p.wg.Wait()
	return p.writer.Close()
}

func extractStringProperty(properties map[string]interface{}, key string) string {
	if t, ok := properties[key]; ok {
		delete(properties, key)
		if v, ok := t.(string); ok {
			return v
		}
	}
	return ""
}

func extractInt32(properties map[string]interface{}, key string) int32 {
	if t, ok := properties[key]; ok {
		delete(properties, key)
		if v, ok := t.(int32); ok {
			return v
		}
	}
	return 0
}

func extractCPID(properties map[string]interface{}) uint32 {
	if t, ok := properties[PresetKeyCPID]; ok {
		delete(properties, PresetKeyCPID)
		if v, ok := t.(uint32); ok {
			return v
		}
	}
	return config.CPID
}

func extractUUID(properties map[string]interface{}) string {
	if t, ok := properties[PresetKeyUUID]; ok {
		delete(properties, PresetKeyUUID)
		if v, ok := t.(string); ok && v != "" {
			return v
		}
	}
	return uuid.New().String()
}

func extractTime(properties map[string]interface{}) string {
	if t, ok := properties[PresetKeyTime]; ok {
		delete(properties, PresetKeyTime)
		switch v := t.(type) {
		case string:
			return v
		case time.Time:
			return v.Format(dateTimeFormat)
		}
	}
	return time.Now().Format(dateTimeFormat)
}
