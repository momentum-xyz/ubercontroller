package config

type Streamchat struct {
	APIKey    string `yaml:"key" envconfig:"STREAMCHAT_KEY"`
	APISecret string `yaml:"secret" envconfig:"STREAMCHAT_SECRET"`
}

func (s *Streamchat) Init() {
}
