package config

import (
	"io"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/pborman/getopt/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"

	"github.com/momentum-xyz/ubercontroller/logger"
	"github.com/momentum-xyz/ubercontroller/utils"
)

// Config : structure to hold configuration
type Config struct {
	Common     Common     `yaml:"common"`
	Settings   Local      `yaml:"settings"`
	Auth       Auth       `yaml:"auth"`
	Postgres   Postgres   `yaml:"postgres"`
	Influx     Influx     `yaml:"influx"`
	UIClient   UIClient   `yaml:"ui_client"`
	Streamchat Streamchat `yaml:"streamchat"`
}

const configFileName = "config.yaml"

var log = logger.L()

func (x *Config) Init() {
	x.Common.Init()
	x.Auth.Init()
	x.Postgres.Init()
	x.Settings.Init()
	x.UIClient.Init()
	x.Influx.Init()
	x.Streamchat.Init()

}

func defConfig() *Config {
	var cfg Config
	cfg.Init()
	return &cfg
}

func readOpts(cfg *Config) {
	helpFlag := false
	getopt.Flag(&helpFlag, 'h', "display help")
	getopt.Flag(&cfg.Settings.LogLevel, 'l', "be verbose")

	getopt.Parse()
	if helpFlag {
		getopt.Usage()
		os.Exit(0)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func readFile(cfg *Config) error {
	if !fileExists(configFileName) {
		return nil
	}

	f, err := os.Open(configFileName)
	if err != nil {
		return errors.WithMessage(err, "failed to open file")
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		if err != io.EOF {
			return errors.WithMessage(err, "failed to read file")
		}
		return nil
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return errors.WithMessage(err, "failed to unmarshal data")
	}

	// filling config additional raw data
	var rawData map[string]any
	if err := yaml.Unmarshal(data, &rawData); err != nil {
		return errors.WithMessage(err, "failed to unmarshal raw data")
	}

	cfg.Auth.rawData = utils.GetFromAnyMap(rawData, "auth", map[string]any{})

	return nil
}

func readEnv(cfg *Config) error {
	return envconfig.Process("", cfg)
}

func prettyPrint(cfg *Config) {
	d, _ := yaml.Marshal(cfg)
	log.Infof("--- Config ---\n%s\n\n", string(d))
}

// GetConfig : get config file
func GetConfig() *Config {
	cfg := defConfig()

	if err := readFile(cfg); err != nil {
		log.Fatalf("GetConfig: failed to read file: %s", err)
	}
	if err := readEnv(cfg); err != nil {
		log.Fatalf("GetConfig: failed to read env: %s", err)
	}
	readOpts(cfg)

	logger.SetLevel(zapcore.Level(cfg.Settings.LogLevel))
	prettyPrint(cfg)

	return cfg
}

func getEnv(key string) (string, bool) {
	key = strings.ToUpper(key)

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if pair[0] == key {
			return pair[1], true
		}
	}

	return "", false
}
