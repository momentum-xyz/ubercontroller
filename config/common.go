package config

type Common struct {
	AgoraAppCertificate string `yaml:"agora_app_certificate" envconfig:"AGORA_APP_CERTIFICATE"`
	RenderDefaultURL    string `yaml:"render_default_url" envconfig:"RENDER_DEFAULT_URL"`
	RenderInternalURL   string `yaml:"render_internal_url" envconfig:"RENDER_INTERNAL_URL"`
}

func (x *Common) Init() {
}
