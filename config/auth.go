package config

import "github.com/momentum-xyz/ubercontroller/types"

type Auth struct {
	OIDCProviders         map[string]types.ConfigAuthProviderType `yaml:"oidc_providers" envconfig:"OIDC_PROVIDERS"`
	OIDCURLs              map[string]string                       `yaml:"oidc_urls" envconfig:"OIDC_URLS"`
	OIDCIntospectURLs     map[string]string                       `yaml:"oidc_introspect_urls" envconfig:"OIDC_INTROSPECT_URLS"`
	OIDCClientIDs         map[string]string                       `yaml:"oidc_client_ids" envconfig:"OIDC_CLIENT_IDS"`
	OIDCSecrets           map[string]string                       `yaml:"oidc_secrets" envconfig:"OIDC_SECRETS"`
	OIDCAdditionalParties map[string]string                       `yaml:"oidc_additional_parties" envconfig:"OIDC_ADDITIONAL_PARTIES"`
}

func (x *Auth) Init() {
	x.OIDCProviders = map[string]types.ConfigAuthProviderType{}
	x.OIDCURLs = map[string]string{}
	x.OIDCIntospectURLs = map[string]string{}
	x.OIDCClientIDs = map[string]string{}
	x.OIDCSecrets = map[string]string{}
	x.OIDCAdditionalParties = map[string]string{}
}
