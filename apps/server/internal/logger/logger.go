package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(env string) {
	var cfg zap.Config
	if env == "prod" || env == "production" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var err error
	log, err = cfg.Build()
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}

	zap.ReplaceGlobals(log)
}

func L() *zap.Logger {
	if log == nil {
		return zap.NewNop()
	}
	return log
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
