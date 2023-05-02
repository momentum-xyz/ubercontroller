package config

type UIClient struct {
	AgoraAppID             string `yaml:"agora_app_id" json:"AGORA_APP_ID" envconfig:"AGORA_APP_ID"`
	BlockchainWsServer     string `yaml:"blockchain_ws_server" json:"BLOCKCHAIN_WS_SERVER" envconfig:"BLOCKCHAIN_WS_SERVER"`
	FrontendURL            string `yaml:"frontend_url" json:"-" envconfig:"FRONTEND_URL"`
	NFTAdminAddress        string `yaml:"nft_admin_address" json:"NFT_ADMIN_ADDRESS" envconfig:"NFT_ADMIN_ADDRESS"`
	NFTCollectionOdysseyID string `yaml:"nft_collection_odyssey_id" json:"NFT_COLLECTION_ODYSSEY_ID" envconfig:"NFT_COLLECTION_ODYSSEY_ID"`
	StreamchatKey          string `yaml:"streamchat_key" json:"STREAMCHAT_KEY" envconfig:"STREAMCHAT_KEY"`
}

func (c *UIClient) Init() {
}
