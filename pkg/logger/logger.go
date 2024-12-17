package logger

import (
	"context"
	"os"
	"runtime/debug"
	"strconv"
	"sync"

	"github.com/danushk97/image-analyzer/pkg/contextkey"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey struct{}

var (
	once   sync.Once
	logger *Logger
)

type Logger struct {
	*zap.SugaredLogger
}

func (l *Logger) newEntry() *Entry {
	return &Entry{
		logger: l,
		Data:   map[string]interface{}{},
	}
}

func NewLogger() *Logger {
	once.Do(func() {
		stdout := zapcore.AddSync(os.Stdout)

		logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
		if err != nil {
			logLevel = int(zap.InfoLevel)
		}

		level := zap.NewAtomicLevelAt(zapcore.Level(logLevel))

		productionCfg := zap.NewProductionEncoderConfig()
		productionCfg.TimeKey = "timestamp"
		productionCfg.TimeKey = "message"
		productionCfg.CallerKey = "caller"
		productionCfg.LevelKey = "level"
		productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		developmentCfg := zap.NewDevelopmentEncoderConfig()
		developmentCfg.TimeKey = "timestamp"
		developmentCfg.MessageKey = "message"
		developmentCfg.CallerKey = "caller"
		developmentCfg.LevelKey = "level"
		jsonEncoder := zapcore.NewJSONEncoder(developmentCfg)

		core := zapcore.NewCore(jsonEncoder, stdout, level)
		if os.Getenv("APP_ENV") != "dev" {
			core = zapcore.NewTee(
				zapcore.NewCore(jsonEncoder, stdout, level),
			)
		}
		var gitRevision string

		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, v := range buildInfo.Settings {
				if v.Key == "vcs.revision" {
					gitRevision = v.Value
					break
				}
			}
		}

		logger = &Logger{
			zap.New(core).
				With(
					zap.String(
						"git_revision",
						gitRevision,
					),
					zap.String("go_version", buildInfo.GoVersion),
				).Sugar(),
		}
	})

	return logger
}

// WithContext returns a copy of ctx with the Logger attached.
func Ctx(ctx context.Context) *Entry {
	if l, ok := ctx.Value(ctxKey{}).(*Entry); ok {
		return l
	}

	logger := NewLogger()
	entry := logger.newEntry()
	ctx = context.WithValue(ctx, ctxKey{}, entry)

	fields := map[string]interface{}{}

	if ctx != nil {
		if val, ok := ctx.Value(contextkey.RequestID).(string); ok {
			fields[contextkey.RequestID.String()] = val
		}

		if val, ok := ctx.Value(contextkey.RequestPath).(string); ok {
			fields[contextkey.RequestPath.String()] = val
		}
	}

	if len(fields) > 0 {
		entry = entry.WithFields(fields)
	}

	return entry
}
