package config

type Common struct {
	AgoraAppCertificate string `yaml:"agora_app_certificate" envconfig:"AGORA_APP_CERTIFICATE"`
	MnemonicPhrase      string `yaml:"mnemonic_phrase" envconfig:"MNEMONIC_PHRASE"`
	AllowCORS           bool   `yaml:"allow_cors" envconfig:"ALLOW_CORS"`

	// Enabled pprof http endpoints. If enabled, point your browser at /debug/pprof
	PProfAPI bool `yaml:"pprof_api" envconfig:"PPROF_API"`
}

func (x *Common) Init() {
}
