package config

type Auth struct {
	OIDCProviders             string `yaml:"oidc_providers" envconfig:"OIDC_PROVIDERS"`
	OIDCWeb3ID                string `yaml:"oidc_web3_id" envconfig:"OIDC_WEB3_ID"`
	OIDCURL                   string `yaml:"oidc_url" envconfig:"OIDC_URL"`
	OIDCWeb3IntrospectionURL  string `yaml:"oidc_web3_introspection_url" envconfig:"OIDC_WEB3_INTROSPECTION_URL"`
	OIDCWeb3Secret            string `yaml:"oidc_web3_secret" envconfig:"OIDC_WEB3_SECRET"`
	OIDCWeb3AdditionalParty   string `yaml:"oidc_web3_additional_party" envconfig:"OIDC_WEB3_ADDITIONAL_PARTY"`
	OIDCGuestId               string `yaml:"oidc_guest_id" envconfig:"OIDC_GUEST_ID"`
	OIDCGuestIntrospectionURL string `yaml:"oidc_guest_introspection_url" envconfig:"OIDC_GUEST_INTROSPECTION_URL"`
	OIDCGuestSecret           string `yaml:"oidc_guest_secret" envconfig:"OIDC_GUEST_SECRET"`
	OIDCGuestAdditionalParty  string `yaml:"oidc_guest_additional_party" envconfig:"OIDC_GUEST_ADDITIONAL_PARTY"`
}

func (x *Auth) Init() {
	x.OIDCProviders = "momentum,web3"
	x.OIDCWeb3ID = ""
	x.OIDCURL = ""
	x.OIDCWeb3IntrospectionURL = ""
	x.OIDCWeb3Secret = ""
	x.OIDCWeb3AdditionalParty = ""
	x.OIDCGuestId = ""
	x.OIDCGuestIntrospectionURL = ""
	x.OIDCGuestSecret = ""
	x.OIDCGuestAdditionalParty = ""
}
