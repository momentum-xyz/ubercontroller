package config

type UIClient struct {
	FrontendURL                   string `yaml:"frontend_url" json:"-" envconfig:"FRONTEND_URL"`
	KeycloakOpenIDConnectURL      string `yaml:"keycloak_open_id_connect_url" json:"KEYCLOAK_OPENID_CONNECT_URL" envconfig:"KEYCLOAK_OPENID_CONNECT_URL"`
	KeycloakOpenIDClientID        string `yaml:"keycloak_open_id_client_id" json:"KEYCLOAK_OPENID_CLIENT_ID" envconfig:"KEYCLOAK_OPENID_CLIENT_ID"`
	KeycloakOpenIDScope           string `yaml:"keycloak_open_id_scope" json:"KEYCLOAK_OPENID_SCOPE" envconfig:"KEYCLOAK_OPENID_SCOPE"`
	HydraOpenIDConnectURL         string `yaml:"hydra_open_id_connect_url" json:"HYDRA_OPENID_CONNECT_URL" envconfig:"HYDRA_OPENID_CONNECT_URL"`
	HydraOpenIDClientID           string `yaml:"hydra_open_id_client_id" json:"HYDRA_OPENID_CLIENT_ID" envconfig:"HYDRA_OPENID_CLIENT_ID"`
	HydraOpenIDGuestClientID      string `yaml:"hydra_open_id_guest_client_id" json:"HYDRA_OPENID_GUEST_CLIENT_ID" envconfig:"HYDRA_OPENID_GUEST_CLIENT_ID"`
	HydraOpenIDScope              string `yaml:"hydra_open_id_scope" json:"HYDRA_OPENID_SCOPE" envconfig:"HYDRA_OPENID_SCOPE"`
	Web3IdentityProviderURL       string `yaml:"web_3_identity_provider_url" json:"WEB3_IDENTITY_PROVIDER_URL" envconfig:"WEB3_IDENTITY_PROVIDER_URL"`
	GuestIdentityProviderURL      string `yaml:"guest_identity_provider_url" json:"GUEST_IDENTITY_PROVIDER_URL" envconfig:"GUEST_IDENTITY_PROVIDER_URL"`
	SentryDSN                     string `yaml:"sentry_dsn" json:"SENTRY_DSN" envconfig:"SENTRY_DSN"`
	AgoraAppID                    string `yaml:"agora_app_id" json:"AGORA_APP_ID" envconfig:"AGORA_APP_ID"`
	AuthServiceURL                string `yaml:"auth_service_url" json:"AUTH_SERVICE_URL" envconfig:"AUTH_SERVICE_URL"`
	GoogleAPIClientID             string `yaml:"google_api_client_id" json:"GOOGLE_API_CLIENT_ID" envconfig:"GOOGLE_API_CLIENT_ID"`
	GoogleAPIDeveloperKey         string `yaml:"google_api_developer_key" json:"GOOGLE_API_DEVELOPER_KEY" envconfig:"GOOGLE_API_DEVELOPER_KEY"`
	MiroAppID                     string `yaml:"miro_app_id" json:"MIRO_APP_ID" envconfig:"MIRO_APP_ID"`
	ReactAppYoutubeKey            string `yaml:"react_app_youtube_key" json:"YOUTUBE_KEY" envconfig:"REACT_APP_YOUTUBE_KEY"`
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
