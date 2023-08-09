package config

type Media struct {
	Fontpath  string `yaml:"fontpath" envconfig:"RENDER_FONT_PATH"`
	Imagepath string `yaml:"image_path" envconfig:"RENDER_IMAGE_PATH"`
	Audiopath string `yaml:"audio_path" envconfig:"RENDER_AUDIO_PATH"`
	Videopath string `yaml:"video_path" envconfig:"RENDER_VIDEO_PATH"`
	Assetpath string `yaml:"asset_path" envconfig:"RENDER_ASSET_PATH"`
	LogLevel  int8   `yaml:"loglevel"  envconfig:"RENDER_LOGLEVEL"`
}

func (x *Media) Init() {
	x.Fontpath = "./storage/fonts"
	x.Imagepath = "./storage/images"
	x.Videopath = "./storage/videos"
	x.Audiopath = "./storage/tracks"
	x.Assetpath = "./storage/assets"
	x.LogLevel = 0
}
