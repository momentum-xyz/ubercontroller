package config

type Arbitrum struct {
	BlockchainID    string `yaml:"arbitrum_chain_id" envconfig:"ARBITRUM_CHAIN_ID"`
	MOMTokenAddress string `yaml:"arbitrum_mom_token_address" envconfig:"ARBITRUM_MOM_TOKEN_ADDRESS"`
	DADTokenAddress string `yaml:"arbitrum_dad_token_address" envconfig:"ARBITRUM_DAD_TOKEN_ADDRESS"`
	StakeAddress    string `yaml:"arbitrum_stake_token_address" envconfig:"ARBITRUM_STAKE_ADDRESS"`
	NFTAddress      string `yaml:"arbitrum_nft_address" envconfig:"ARBITRUM_NFT_ADDRESS"`
	FaucetAddress   string `yaml:"arbitrum_faucet_address" envconfig:"ARBITRUM_FAUCET_ADDRESS"`
	RPCURL          string `yaml:"arbitrum_rpc_url" envconfig:"ARBITRUM_RPC_URL"`
	WSURL           string `yaml:"arbitrum_ws_url" envconfig:"ARBITRUM_WS_URL"`
}

func (a *Arbitrum) Init() {
	// Default values connect to a private testnet.
	a.BlockchainID = "412346"
	a.MOMTokenAddress = "0x0147a5cB13Ab1D75f09385e53118AD00D4B74778"
	a.DADTokenAddress = "0x37EcF9B8A0fAceF220f2Bd3C7C30Cbf416433564"
	a.StakeAddress = "0xBCb202a9B77d9A44B54e30f7adF30d9eBd4Cd145"
	a.NFTAddress = "0x9f89a52a29A9964DE6B2450b8edA92Ad2d2ba146"
	a.FaucetAddress = "0x9E760F1CddA0694B6156076C60657118CF874289"
	a.RPCURL = "https://bcdev.antst.net:8547"
	a.WSURL = "wss://bcdev.antst.net:8548"
}
