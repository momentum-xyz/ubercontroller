package config

type UIClient struct {
	AgoraAppID             string `yaml:"agora_app_id" json:"AGORA_APP_ID" envconfig:"AGORA_APP_ID"`
	BlockchainWsServer     string `yaml:"blockchain_ws_server" json:"BLOCKCHAIN_WS_SERVER" envconfig:"BLOCKCHAIN_WS_SERVER"`
	BlockchainID           string `json:"BLOCKCHAIN_ID"`
	ContractDAD            string `json:"CONTRACT_DAD_ADDRESS"`
	ContractFaucet         string `json:"CONTRACT_FAUCET_ADDRESS"`
	ContractMOM            string `json:"CONTRACT_MOM_ADDRESS"`
	ContractNFT            string `json:"CONTRACT_NFT_ADDRESS"`
	ContractStake          string `json:"CONTRACT_STAKING_ADDRESS"`
	FrontendURL            string `yaml:"frontend_url" json:"-" envconfig:"FRONTEND_URL"`
	NFTAdminAddress        string `yaml:"nft_admin_address" json:"NFT_ADMIN_ADDRESS" envconfig:"NFT_ADMIN_ADDRESS"`
	NFTCollectionOdysseyID string `yaml:"nft_collection_odyssey_id" json:"NFT_COLLECTION_ODYSSEY_ID" envconfig:"NFT_COLLECTION_ODYSSEY_ID"`
	StreamchatKey          string `yaml:"streamchat_key" json:"STREAMCHAT_KEY" envconfig:"STREAMCHAT_KEY"`
}

func (c *UIClient) Init(arb Arbitrum) {
	c.BlockchainID = arb.BlockchainID
	c.ContractMOM = arb.MOMTokenAddress
	c.ContractDAD = arb.DADTokenAddress
	c.ContractStake = arb.StakeAddress
	c.ContractNFT = arb.NFTAddress
	c.ContractFaucet = arb.FaucetAddress
}
