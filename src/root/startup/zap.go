package startup

import (
	"github.com/wgjtyu/logMansion/lib"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"wfRadius/src/config"
)

func NewZap(cfg *config.MConfig) *zap.Logger {
	var zlogger *zap.Logger
	if cfg.Environment == config.EnvirIsProd {
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   os.Args[1] + "/zaplog",
			MaxSize:    2, // megabytes
			MaxBackups: 20,
		})
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			w,
			zap.DebugLevel,
		)
		zlogger = zap.New(core)
		lmCore := lib.NewCore(cfg.LogBackend, cfg.LogProjectID, cfg.LogKey, zapcore.DebugLevel)
		zlogger = zlogger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(c, lmCore)
		}))
	} else {
		zlogger, _ = zap.NewDevelopment()
		/*
			lmCore := lib.NewCore(cfg.Log.LogBackend,
				cfg.Log.LogProjectID,
				cfg.Log.LogKey,
				zapcore.InfoLevel)
			zlogger = zlogger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
				return zapcore.NewTee(c, lmCore)
			}))
		*/
	}
	zap.ReplaceGlobals(zlogger)
	return zlogger
}
