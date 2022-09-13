// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import "errors"

var (
	ErrInvalidDevicecode        = errors.New("invalid devicecode or distinctid")
	ErrInvalidEvent             = errors.New("invalid event")
	ErrInvalidCPID              = errors.New("invalid cpid")
	ErrInvalidOpenID            = errors.New("invalid openid")
	ErrInvalidAppID             = errors.New("invalid appid")
	ErrAppKeyNotExistx          = errors.New("appkey not exists")
	ErrInvalidType              = errors.New("invalid type")
	ErrInvalidIMSConversationID = errors.New("invalid ims conversation id")

	errProducerShutdown = errors.New("producer already shut down")
)

type Error struct {
	traceID string
	err     string
}

func (e *Error) Error() string {
	return e.err
}

func errWithTraceID(err error, traceID string) error {
	if err == nil {
		return nil
	}
	return &Error{
		traceID: traceID,
		err:     err.Error(),
	}
}

// ErrTraceID 返回发生错误的 TraceID
func ErrTraceID(err error) string {
	if err == nil {
		return ""
	}

	if e, ok := err.(*Error); ok {
		return e.traceID
	}

	if e := errors.Unwrap(err); e != nil {
		return ErrTraceID(e)
	}

	return ""
}
