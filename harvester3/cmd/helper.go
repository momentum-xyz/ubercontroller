package helper

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
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

func GetOdysseyContracts(environment string) (mom, dad, nft, stake common.Address) {
	if environment == "main_net" {
		mom = common.HexToAddress("0x457fd0Ee3Ce35113ee414994f37eE38518d6E7Ee")
		dad = common.HexToAddress("0x457fd0Ee3Ce35113ee414994f37eE38518d6E7Ee")
		nft = common.HexToAddress("0x97E0B10D89a494Eb5cfFCc72853FB0750BD64AcD")
		stake = common.HexToAddress("0xe9C6d7Cd04614Dde6Ca68B62E6fbf23AC2ECe2F8")
	}

	if environment == "main_net" {
		mom = common.HexToAddress("0x457fd0Ee3Ce35113ee414994f37eE38518d6E7Ee")
		dad = common.HexToAddress("0x457fd0Ee3Ce35113ee414994f37eE38518d6E7Ee")
		nft = common.HexToAddress("0x97E0B10D89a494Eb5cfFCc72853FB0750BD64AcD")
		stake = common.HexToAddress("0xe9C6d7Cd04614Dde6Ca68B62E6fbf23AC2ECe2F8")
	}

	return
}
