package config

type Local struct {
	Address          string `yaml:"bind_address" envconfig:"UBERCONTROLLER_BIND_ADDRESS"`
	Port             uint   `yaml:"bind_port" envconfig:"UBERCONTROLLER_BIND_PORT"`
	LogLevel         int    `yaml:"loglevel"  envconfig:"UBERCONTROLLER_LOGLEVEL"`
	ExtensionStorage string `yaml:"storage"  envconfig:"UBERCONTROLLER_STORAGE"`
	SeedDataFiles    string `yaml:"seed_data_files" envconfig:"UBERCONTROLLER_SEED_DATA_FILES"`
	FrontendServeDir string `yaml:"frontend_serve_dir" envconfig:"FRONTEND_SERVE_DIR"`
	// TODO: rename FrontendURL to avoid confusing this with 'the frontend'
	FrontendURL string `yaml:"frontend_url" json:"-" envconfig:"FRONTEND_URL"` // URL where this instance is reachable (e.g. when behind a proxy)
}

func (x *Local) Init() {
	x.LogLevel = 1
	x.Address = "0.0.0.0"
	x.Port = 4000
	x.ExtensionStorage = "/opt/ubercontroller"
	x.SeedDataFiles = "./seed/data"
	x.FrontendURL = "http://localhost:4000"
}
