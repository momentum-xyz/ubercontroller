package config

type UIClient struct {
	AgoraAppID                    string `yaml:"agora_app_id" json:"AGORA_APP_ID" envconfig:"AGORA_APP_ID"`
	BlockchainWsServer            string `yaml:"blockchain_ws_server" json:"BLOCKCHAIN_WS_SERVER" envconfig:"BLOCKCHAIN_WS_SERVER"`
	FrontendURL                   string `yaml:"frontend_url" json:"-" envconfig:"FRONTEND_URL"`
	NFTAdminAddress               string `yaml:"nft_admin_address" json:"NFT_ADMIN_ADDRESS" envconfig:"NFT_ADMIN_ADDRESS"`
	NFTCollectionOdysseyID        string `yaml:"nft_collection_odyssey_id" json:"NFT_COLLECTION_ODYSSEY_ID" envconfig:"NFT_COLLECTION_ODYSSEY_ID"`
	StreamchatKey                 string `yaml:"streamchat_key" json:"STREAMCHAT_KEY" envconfig:"STREAMCHAT_KEY"`
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
