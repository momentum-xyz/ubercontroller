package config

type Common struct {
	IntrospectURL string `yaml:"introspect_url" envconfig:"BACKEND_INTROSPECT_URL"`
}

func (x *Common) Init() {
	x.IntrospectURL = "http://backend-service-service.default:4000/api/v3/backend/auth/introspect"
}
