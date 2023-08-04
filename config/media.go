package config

type Media struct {
	Address   string `yaml:"bind_address" envconfig:"RENDER_BIND_ADDRESS"`
	Port      uint   `yaml:"bind_port" envconfig:"RENDER_BIND_PORT"`
	Fontpath  string `yaml:"fontpath" envconfig:"RENDER_FONT_PATH"`
	Imagepath string `yaml:"image_path" envconfig:"RENDER_IMAGE_PATH"`
	Audiopath string `yaml:"audio_path" envconfig:"RENDER_AUDIO_PATH"`
	Videopath string `yaml:"video_path" envconfig:"RENDER_VIDEO_PATH"`
	Assetpath string `yaml:"asset_path" envconfig:"RENDER_ASSET_PATH"`
	LogLevel  int8   `yaml:"loglevel"  envconfig:"RENDER_LOGLEVEL"`
	// Enabled pprof http endpoints. If enabled, point your browser at /debug/pprof
	PProfAPI bool `yaml:"pprof_api" envconfig:"RENDER_PPROF_API"`
}

func (x *Media) Init() {
	x.Address = "0.0.0.0"
	x.Port = 4000
	x.Fontpath = "./fonts"
	x.Imagepath = "./images"
	x.Videopath = "./videos"
	x.Audiopath = "./images/tracks"
	x.Assetpath = "./images/assets"
	x.LogLevel = 0
}
