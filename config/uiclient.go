package config

type UIClient struct {
	FrontendURL                   string `yaml:"frontend_url" json:"-" envconfig:"FRONTEND_URL"`
	AgoraAppID                    string `yaml:"agora_app_id" json:"AGORA_APP_ID" envconfig:"AGORA_APP_ID"`
	UnityClientStreamingAssetsURL string `yaml:"unity_client_streaming_assets_url" json:"UNITY_CLIENT_STREAMING_ASSETS_URL" envconfig:"UNITY_CLIENT_STREAMING_ASSETS_URL"`
	UnityClientCompanyName        string `yaml:"unity_client_company_name" json:"UNITY_CLIENT_COMPANY_NAME" envconfig:"UNITY_CLIENT_COMPANY_NAME"`
	UnityClientProductName        string `yaml:"unity_client_product_name" json:"UNITY_CLIENT_PRODUCT_NAME" envconfig:"UNITY_CLIENT_PRODUCT_NAME"`
	UnityClientProductVersion     string `yaml:"unity_client_product_version" json:"UNITY_CLIENT_PRODUCT_VERSION" envconfig:"UNITY_CLIENT_PRODUCT_VERSION"`
}

func (c *UIClient) Init() {
	c.UnityClientStreamingAssetsURL = "StreamingAssets"
	c.UnityClientCompanyName = "Odyssey"
	c.UnityClientProductName = "Odyssey Momentum"
	c.UnityClientProductVersion = "0.1"
}
