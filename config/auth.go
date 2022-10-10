package config

import (
	"fmt"
	"github.com/momentum-xyz/ubercontroller/utils"
)

type Auth struct {
	OIDCProviders []string       `yaml:"oidc_providers" envconfig:"OIDC_PROVIDERS"`
	OIDCURL       string         `yaml:"oidc_url" envconfig:"OIDC_URL"`
	rawData       map[string]any `yaml:"-"`
}

func (x *Auth) Init() {
	x.OIDCProviders = []string{"web3,guest"}
	x.rawData = make(map[string]any)
}

func (x *Auth) GetIDByProvider(provider string) string {
	return x.GetByKey(fmt.Sprintf("oidc_%s_id", provider))
}

func (x *Auth) GetSecretByProvider(provider string) string {
	return x.GetByKey(fmt.Sprintf("oidc_%s_secret", provider))
}

func (x *Auth) GetIntrospectURLByProvider(provider string) string {
	return x.GetByKey(fmt.Sprintf("oidc_%s_introspection_url", provider))
}

func (x *Auth) GetAdditionalPartyByProvider(provider string) string {
	return x.GetByKey(fmt.Sprintf("oidc_%s_additional_party", provider))
}

func (x *Auth) GetByKey(key string) string {
	if env, ok := getEnv(key); ok {
		return env
	}
	return utils.GetFromAnyMap(x.rawData, key, "")
}
