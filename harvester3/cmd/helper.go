package helper

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetZapLogger() *zap.Logger {
	zapCfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		// NOTE: set this false to enable stack trace
		DisableStacktrace: true,
	}

	logger, err := zapCfg.Build()
	if err != nil {
		panic(err)
	}

	return logger
}
