// Copyright (c) 2022. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type BigdataOptions func(p *BigDataLog) error

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

// SetPreset 预制属性
func SetPreset(preset map[string]interface{}) BigdataOptions {
	return func(logData *BigDataLog) error {
		cpID := extractCPID(preset)
		if cpID == 0 {
			return ErrInvalidCPID
		}
		logData.CPID = cpID
		logData.UUID = extractUUID(preset)
		logData.Time = extractTime(preset)
		if preset != nil {
			logData.PlatformID = extractInt32(preset, PresetKeyPlatformID)
			logData.AppID = extractStringProperty(preset, PresetKeyAppID)
			logData.ChannelID = extractStringProperty(preset, PresetKeyChannelID)
			logData.SubChannelID = extractStringProperty(preset, PresetKeySubChannelID)
			logData.IP = extractStringProperty(preset, PresetKeyIP)
		}
		return nil
	}
}

// SetProperties 自定义属性
func SetProperties(properties map[string]interface{}) BigdataOptions {
	return func(p *BigDataLog) error {
		p.Properties = properties
		return nil
	}
}

// SetEvent 事件名
func SetEvent(event string) BigdataOptions {
	return func(logData *BigDataLog) error {
		if event == "" {
			return ErrInvalidEvent
		}
		logData.Type = typeTrack
		logData.Event = event
		return nil
	}
}

// SetUserUpdateType 用户更新类型：user_setonce,user_set
func SetUserUpdateType(updateType string) BigdataOptions {
	return func(p *BigDataLog) error {
		p.Type = typeUser
		p.Event = updateType
		return nil
	}
}

// Tracks 大数据埋点事件上报
// 		devicecode 设备码
// 		distinctID 用户标识, 通常为瑞雪 OpenID
// 		opts 埋点动态参数设置
func (p *Producer) Tracks(devicecode, distinctID string, opts ...BigdataOptions) error {
	if devicecode == "" && distinctID == "" {
		return ErrInvalidDevicecode
	}
	if p.isShutDown.Load() {
		return errProducerShutdown
	}
	p.wg.Add(1)
	defer p.wg.Done()

	logData := &BigDataLog{
		DistinctID: distinctID,
		Devicecode: devicecode,
	}
	for _, opt := range opts {
		err := opt(logData)
		if err != nil {
			return err
		}
	}
	if logData.Type == "" {
		return ErrInvalidType
	}
	if logData.CPID == 0 {
		if config.CPID == 0 {
			return ErrInvalidCPID
		}
		logData.CPID = config.CPID
	}
	if logData.PlatformID <= 0 {
		logData.PlatformID = 10
	}
	if logData.UUID == "" {
		logData.UUID = uuid.New().String()
	}
	if logData.Time == "" {
		logData.Time = time.Now().Format(dateTimeFormat)
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
		if v, ok := t.(string); ok {
			return v
		}
	}
	return ""
}

func extractInt32(properties map[string]interface{}, key string) int32 {
	if t, ok := properties[key]; ok {
		if v, ok := t.(int32); ok {
			return v
		}
	}
	return 0
}

func extractCPID(properties map[string]interface{}) uint32 {
	if t, ok := properties[PresetKeyCPID]; ok {
		if v, ok := t.(uint32); ok {
			return v
		}
	}
	return config.CPID
}

func extractUUID(properties map[string]interface{}) string {
	if t, ok := properties[PresetKeyUUID]; ok {
		if v, ok := t.(string); ok && v != "" {
			return v
		}
	}
	return uuid.New().String()
}

func extractTime(properties map[string]interface{}) string {
	if t, ok := properties[PresetKeyTime]; ok {
		switch v := t.(type) {
		case string:
			return v
		case time.Time:
			return v.Format(dateTimeFormat)
		}
	}
	return time.Now().Format(dateTimeFormat)
}
