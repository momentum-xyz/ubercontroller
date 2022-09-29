package config

type Common struct {
	IntrospectURL string `yaml:"introspect_url" envconfig:"BACKEND_INTROSPECT_URL"`
	APIPrefix     string `yaml:"api_prefix" envconfig:"API_PREFIX"`
}

func (x *Common) Init() {
	x.IntrospectURL = "http://backend-service-service.default:4000/api/v3/backend/auth/introspect"
	x.APIPrefix = "/api/v4"
}
