package config

type Auth struct {
	rawData map[string]any `yaml:"-"`
}

func (x *Auth) Init() {
	x.rawData = make(map[string]any)
}
