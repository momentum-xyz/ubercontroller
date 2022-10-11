package config

type Local struct {
	Address          string `yaml:"bind_address" envconfig:"UBERCONTROLLER_BIND_ADDRESS"`
	Port             uint   `yaml:"bind_port" envconfig:"UBERCONTROLLER_BIND_PORT"`
	LogLevel         int    `yaml:"loglevel"  envconfig:"UBERCONTROLLER_LOGLEVEL"`
	ExtensionStorage string `yaml:"storage"  envconfig:"UBERCONTROLLER_STORAGE"`
}

func (x *Local) Init() {
	x.LogLevel = 1
	x.Address = "0.0.0.0"
	x.Port = 4000
	x.ExtensionStorage = "/opt/ubercontroller"
}
