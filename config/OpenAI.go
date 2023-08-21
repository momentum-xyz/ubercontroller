package config

type OpenAI struct {
	URL string `yaml:"open_ai_url" envconfig:"OPEN_AI_URL"`
}

func (a *OpenAI) Init() {
	a.URL = "https://api.openai.com/v1/chat/completions"
}
