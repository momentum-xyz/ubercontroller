package config

import (
	"fmt"
	"io"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/pborman/getopt/v2"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"

	"github.com/momentum-xyz/ubercontroller/logger"
)

// Config : structure to hold configuration
type Config struct {
	Common     Common     `yaml:"common"`
	Settings   Local      `yaml:"settings"`
	Postgres   Postgres   `yaml:"postgres"`
	Influx     Influx     `yaml:"influx"`
	UIClient   UIClient   `yaml:"ui_client"`
	Streamchat Streamchat `yaml:"streamchat"`
	Arbitrum   Arbitrum   `yaml:"arbitrum"`
}

const configFileName = "config.yaml"

func (x *Config) Init() {
	x.Common.Init()
	x.Postgres.Init()
	x.Settings.Init()
	x.Influx.Init()
	x.Streamchat.Init()
	x.Arbitrum.Init()
	x.UIClient.Init(x.Arbitrum)
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
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		if err != io.EOF {
			return fmt.Errorf("failed to read config file: %w", err)
		}
		return nil
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to unmarhal config: %w", err)
	}

	return nil
}

func readEnv(cfg *Config) error {
	return envconfig.Process("", cfg)
}

func prettyPrint(cfg *Config) {
	d, _ := yaml.Marshal(cfg)
	fmt.Printf("--- Config ---\n%s\n\n", string(d))
}

// GetConfig : get config file
func GetConfig() (*Config, error) {
	cfg := defConfig()

	if err := readFile(cfg); err != nil {
		return nil, fmt.Errorf("GetConfig: %w", err)
	}
	if err := readEnv(cfg); err != nil {
		return nil, fmt.Errorf("GetConfig: %w", err)
	}
	readOpts(cfg)

	logger.SetLevel(zapcore.Level(cfg.Settings.LogLevel))
	prettyPrint(cfg)

	return cfg, nil
}
