// Copyright (c) 2021. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"fmt"
	"log"
	"os"
)

// Logger 框架内部使用日志接口定义
type Logger interface {
	// Infof 格式化输出信息级日志
	Infof(tmp string, v ...interface{})

	// Errorf 格式化输出错误日志
	Errorf(tmp string, v ...interface{})

	// Debugf 格式化输出 Debug 日志
	Debugf(tmp string, v ...interface{})
}

// defaultLogger 默认日志接口定义, 使用标准库的 log 包
type defaultLogger struct {
	stdout *log.Logger
}

// Infof 输出信息级日志
func (l *defaultLogger) Infof(tmp string, v ...interface{}) {
	l.stdout.Println("[INFO]", fmt.Sprintf(tmp, v...))
}

// Errorf 输出错误日志
func (l *defaultLogger) Errorf(tmp string, v ...interface{}) {
	l.stdout.Println("[ERROR]", fmt.Sprintf(tmp, v...))
}

// Debugf 输出 Debug 日志
func (l *defaultLogger) Debugf(tmp string, v ...interface{}) {
	l.stdout.Println("[DEBUG]", fmt.Sprintf(tmp, v...))
}

// 框架默认日志接口对象缓存
var logger Logger = &defaultLogger{
	stdout: log.New(os.Stdout, "", log.LstdFlags),
}

// GetLogger 获取框架当前使用的日志处理接口
func GetLogger() Logger {
	return logger
}

// RegisterLogger 注册框架当前使用的日志处理接口
func RegisterLogger(l Logger) {
	logger = l
}
