package utils

import (
	"github.com/sid-sun/ioctl-api/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

var logLevels map[string]zapcore.Level = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"warn":  zap.WarnLevel,
	"error": zap.ErrorLevel,
}

func InitLogger(cfg *config.Config) {
	lgrCfg := zap.NewProductionConfig()
	if cfg.App.GetEnv() == config.Environments["dev"] {
		lgrCfg.Development = true
	}
	lgrCfg.Level = zap.NewAtomicLevelAt(logLevels[cfg.App.GetLogLevel()])

	lgr, err := lgrCfg.Build()
	if err != nil {
		panic(err)
	}
	Logger = lgr
}
