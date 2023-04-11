package config

type Common struct {
	AgoraAppCertificate string `yaml:"agora_app_certificate" envconfig:"AGORA_APP_CERTIFICATE"`
	RenderInternalURL   string `yaml:"render_internal_url" envconfig:"RENDER_INTERNAL_URL"`
	MnemonicPhrase      string `yaml:"mnemonic_phrase" envconfig:"MNEMONIC_PHRASE"`
	DisableCORS         bool   `yaml:"disable_cors" envconfig:"DISABLE_CORS"`
}

func (x *Common) Init() {
}
