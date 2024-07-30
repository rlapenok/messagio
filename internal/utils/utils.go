package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {
	var err error
	encoder := zapcore.EncoderConfig{

		TimeKey:      "time",
		LevelKey:     "level",
		NameKey:      "logger",
		CallerKey:    "caller",
		FunctionKey:  zapcore.OmitKey,
		MessageKey:   "msg",
		LineEnding:   zapcore.DefaultLineEnding,
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = encoder
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.Development = true
	cfg.Encoding = "json"
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stdout"}
	Logger, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}
