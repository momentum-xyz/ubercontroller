package config

type UIClient struct {
	AgoraAppID             string `yaml:"agora_app_id" json:"AGORA_APP_ID" envconfig:"AGORA_APP_ID"`
	BlockchainWsServer     string `yaml:"blockchain_ws_server" json:"BLOCKCHAIN_WS_SERVER" envconfig:"BLOCKCHAIN_WS_SERVER"`
	BlockchainID           string `yaml:"arbitrum_chain_id" json:"BLOCKCHAIN_ID" envconfig:"ARBITRUM_CHAIN_ID"`
	ContractDAD            string `yaml:"arbitrum_dad_token_address" json:"CONTRACT_DAD_ADDRESS" envconfig:"ARBITRUM_DAD_TOKEN_ADDRESS"`
	ContractFaucet         string `yaml:"arbitrum_faucet_contract_address" json:"CONTRACT_FAUCET_ADDRESS" envconfig:"ARBITRUM_FAUCET_CONTRACT_ADDRESS"`
	ContractMOM            string `yaml:"arbitrum_mom_token_address" json:"CONTRACT_MOM_ADDRESS" envconfig:"ARBITRUM_MOM_TOKEN_ADDRESS"`
	ContractNFT            string `yaml:"arbitrum_nft_contract_address" json:"CONTRACT_NFT_ADDRESS" envconfig:"ARBITRUM_NFT_CONTRACT_ADDRESS"`
	ContractStake          string `yaml:"arbitrum_stake_token_address" json:"CONTRACT_STAKING_ADDRESS" envconfig:"ARBITRUM_STAKE_TOKEN_ADDRESS"`
	FrontendURL            string `yaml:"frontend_url" json:"-" envconfig:"FRONTEND_URL"`
	NFTAdminAddress        string `yaml:"nft_admin_address" json:"NFT_ADMIN_ADDRESS" envconfig:"NFT_ADMIN_ADDRESS"`
	NFTCollectionOdysseyID string `yaml:"nft_collection_odyssey_id" json:"NFT_COLLECTION_ODYSSEY_ID" envconfig:"NFT_COLLECTION_ODYSSEY_ID"`
	StreamchatKey          string `yaml:"streamchat_key" json:"STREAMCHAT_KEY" envconfig:"STREAMCHAT_KEY"`
}

func (c *UIClient) Init() {
}
