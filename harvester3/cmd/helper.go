package helper

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/momentum-xyz/ubercontroller/config"
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

func MustGetConfig() *config.Config {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
