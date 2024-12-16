package logger

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

const (
	ContextKey = "context"
)

type Entry struct {
	logger *Logger
	Data   map[string]interface{}
}

func (e *Entry) Info(msg string) {
	e.logger.SugaredLogger.Logw(zapcore.InfoLevel, msg, ContextKey, e.Data)
}

func (e *Entry) Debug(msg string) {
	e.logger.SugaredLogger.Logw(zapcore.InfoLevel, msg, ContextKey, e.Data)
}

func (e *Entry) Warn(msg string) {
	e.logger.SugaredLogger.Logw(zapcore.WarnLevel, msg, ContextKey, e.Data)
}

func (e *Entry) Error(msg string) {
	e.logger.SugaredLogger.Logw(zapcore.ErrorLevel, msg, ContextKey, e.Data)
}

func (e *Entry) Panic(msg string) {
	e.logger.SugaredLogger.Logw(zapcore.PanicLevel, ContextKey, msg)
}

func (e *Entry) WithError(err error) *Entry {
	return e.WithFields(map[string]interface{}{
		"error": err.Error(),
	})
}

func (e *Entry) WithField(field string, value interface{}) *Entry {
	return e.WithFields(map[string]interface{}{
		field: value,
	})
}

func (e *Entry) WithFields(fields map[string]interface{}) *Entry {
	data := make(map[string]interface{}, len(e.Data)+len(fields))

	for k, v := range e.Data {
		data[k] = v
	}
	for k, v := range fields {
		v = fmt.Sprintf("%v", v)
		data[k] = v
	}

	return &Entry{logger: e.logger, Data: data}
}
