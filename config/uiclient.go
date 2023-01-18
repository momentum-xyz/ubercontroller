package config

type UIClient struct {
	AgoraAppID                    string `yaml:"agora_app_id" json:"AGORA_APP_ID" envconfig:"AGORA_APP_ID"`
	BlockchainWsServer            string `yaml:"blockchain_ws_server" json:"BLOCKCHAIN_WS_SERVER" envconfig:"BLOCKCHAIN_WS_SERVER"`
	FrontendURL                   string `yaml:"frontend_url" json:"-" envconfig:"FRONTEND_URL"`
	NFTAdminAddress               string `yaml:"nft_admin_address" json:"NFT_ADMIN_ADDRESS" envconfig:"NFT_ADMIN_ADDRESS"`
	NFTCollectionOdysseyID        string `yaml:"nft_collection_odyssey_id" json:"NFT_COLLECTION_ODYSSEY_ID" envconfig:"NFT_COLLECTION_ODYSSEY_ID"`
	StreamchatKey                 string `yaml:"streamchat_key" json:"STREAMCHAT_KEY" envconfig:"STREAMCHAT_KEY"`
	UnityClientCompanyName        string `yaml:"unity_client_company_name" json:"UNITY_CLIENT_COMPANY_NAME" envconfig:"UNITY_CLIENT_COMPANY_NAME"`
	UnityClientProductName        string `yaml:"unity_client_product_name" json:"UNITY_CLIENT_PRODUCT_NAME" envconfig:"UNITY_CLIENT_PRODUCT_NAME"`
	UnityClientProductVersion     string `yaml:"unity_client_product_version" json:"UNITY_CLIENT_PRODUCT_VERSION" envconfig:"UNITY_CLIENT_PRODUCT_VERSION"`
	UnityClientStreamingAssetsURL string `yaml:"unity_client_streaming_assets_url" json:"UNITY_CLIENT_STREAMING_ASSETS_URL" envconfig:"UNITY_CLIENT_STREAMING_ASSETS_URL"`
	UnityClientURL                string `yaml:"unity_client_url" json:"UNITY_CLIENT_URL" envconfig:"UNITY_CLIENT_URL"`
	UnityLoaderFileName           string `yaml:"unity_loader_file_name" json:"UNITY_LOADER_FILE_NAME" envconfig:"UNITY_LOADER_FILE_NAME"`
	UnityFrameworkFileName        string `yaml:"unity_framework_file_name" json:"UNITY_FRAMEWORK_FILE_NAME" envconfig:"UNITY_FRAMEWORK_FILE_NAME"`
	UnityDataFileName             string `yaml:"unity_data_file_name" json:"UNITY_DATA_FILE_NAME" envconfig:"UNITY_DATA_FILE_NAME"`
	UnityCodeFileName             string `yaml:"unity_code_file_name" json:"UNITY_CODE_FILE_NAME" envconfig:"UNITY_CODE_FILE_NAME"`
}

func (c *UIClient) Init() {
	c.UnityClientStreamingAssetsURL = "StreamingAssets"
	c.UnityClientCompanyName = "Odyssey"
	c.UnityClientProductName = "Odyssey Momentum"
	c.UnityClientProductVersion = "0.1"
	c.UnityLoaderFileName = "WebGL.loader.js"
	c.UnityFrameworkFileName = "WebGL.framework.js.gz"
	c.UnityCodeFileName = "WebGL.wasm.gz"
	c.UnityDataFileName = "WebGL.data.gz"
}
