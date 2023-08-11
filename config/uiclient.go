package config

type UIFeatureFlags struct {
	NewsFeed bool `json:"newsfeed" envconfig:"FEATURE_NEWSFEED"`
	BuyNft   bool `yaml:"buy_nft" json:"buy_nft" envconfig:"FEATURE_BUY_NFT"`
}

type UIClient struct {
	AgoraAppID     string         `yaml:"agora_app_id" json:"AGORA_APP_ID" envconfig:"AGORA_APP_ID"`
	BlockchainID   string         `json:"BLOCKCHAIN_ID" envconfig:"ARBITRUM_CHAIN_ID"`
	ContractDAD    string         `json:"CONTRACT_DAD_ADDRESS" envconfig:"ARBITRUM_DAD_TOKEN_ADDRESS"`
	ContractFaucet string         `json:"CONTRACT_FAUCET_ADDRESS" envconfig:"ARBITRUM_FAUCET_ADDRESS"`
	ContractMOM    string         `json:"CONTRACT_MOM_ADDRESS" envconfig:"ARBITRUM_MOM_TOKEN_ADDRESS"`
	ContractNFT    string         `json:"CONTRACT_NFT_ADDRESS" envconfig:"ARBITRUM_NFT_ADDRESS"`
	ContractStake  string         `json:"CONTRACT_STAKING_ADDRESS" envconfig:"ARBITRUM_STAKE_ADDRESS"`
	MintNFTAmount  string         `json:"MINT_NFT_AMOUNT" envconfig:"ARBITRUM_MINT_NFT_AMOUNT"`
	MintNFTDeposit string         `json:"MINT_NFT_DEPOSIT_ADDRESS" envconfig:"ARBITRUM_MINT_NFT_DEPOSIT_ADDRESS"`
	FrontendURL    string         `yaml:"frontend_url" json:"-" envconfig:"FRONTEND_URL"`
	StreamchatKey  string         `yaml:"streamchat_key" json:"STREAMCHAT_KEY" envconfig:"STREAMCHAT_KEY"`
	FeatureFlags   UIFeatureFlags `yaml:"feature_flags" json:"FEATURE_FLAGS" envconfig:"FEATURE_FLAGS"`
}

func (c *UIClient) Init(arb Arbitrum) {
	c.BlockchainID = arb.BlockchainID
	c.ContractMOM = arb.MOMTokenAddress
	c.ContractDAD = arb.DADTokenAddress
	c.ContractStake = arb.StakeAddress
	c.ContractNFT = arb.NFTAddress
	c.ContractFaucet = arb.FaucetAddress
	c.MintNFTAmount = arb.MintNFTAmount
	c.MintNFTDeposit = arb.MintNFTDeposit
}
