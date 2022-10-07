package config

type OIDCProvider struct {
	Name string `yaml:"name"`
}

type Auth struct {
	OIDCProviders     []string          `yaml:"oidc_providers" envconfig:"OIDC_PROVIDERS"`
	OIDCURLs          map[string]string `yaml:"oidc_urls" envconfig:"OIDC_URLS"`
	OIDCIntospectURLs map[string]string `yaml:"oidc_introspect_urls" envconfig:"OIDC_INTROSPECT_URLS"`
	OIDCClientIDs     map[string]string `yaml:"oidc_client_ids" envconfig:"OIDC_CLIENT_IDS"`
	OIDCSecrets       map[string]string `yaml:"oidc_secrets" envconfig:"OIDC_SECRETS"`
}

func (x *Auth) Init() {
	x.OIDCProviders = []string{"web3", "momentum"}
	x.OIDCURLs = map[string]string{}
	x.OIDCIntospectURLs = map[string]string{}
	x.OIDCClientIDs = map[string]string{}
	x.OIDCSecrets = map[string]string{}
}
