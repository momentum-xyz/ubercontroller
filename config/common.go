package config

type Common struct {
	RenderDefaultUrl  string `yaml:"render_default_url" envconfig:"RENDER_DEFAULT_URL"`
	RenderInternalUrl string `yaml:"render_internal_url" envconfig:"RENDER_INTERNAL_URL"`
}

func (x *Common) Init() {
}
