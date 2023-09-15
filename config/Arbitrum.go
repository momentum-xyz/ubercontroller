package config

type Arbitrum struct {
	BlockchainID    string `yaml:"arbitrum_chain_id" envconfig:"ARBITRUM_CHAIN_ID"`
	MOMTokenAddress string `yaml:"arbitrum_mom_token_address" envconfig:"ARBITRUM_MOM_TOKEN_ADDRESS"`
	DADTokenAddress string `yaml:"arbitrum_dad_token_address" envconfig:"ARBITRUM_DAD_TOKEN_ADDRESS"`
	StakeAddress    string `yaml:"arbitrum_stake_token_address" envconfig:"ARBITRUM_STAKE_ADDRESS"`
	NFTAddress      string `yaml:"arbitrum_nft_address" envconfig:"ARBITRUM_NFT_ADDRESS"`
	NodeAddress     string `yaml:"arbitrum_node_address" envconfig:"ARBITRUM_NODE_ADDRESS"`
	MintNFTAmount   string `json:"arbitrum_mint_nft_amount" envconfig:"ARBITRUM_MINT_NFT_AMOUNT"`
	MintNFTDeposit  string `json:"arbitrum_mint_nft_deposit_address" envconfig:"ARBITRUM_MINT_NFT_DEPOSIT_ADDRESS"`
	FaucetAddress   string `yaml:"arbitrum_faucet_address" envconfig:"ARBITRUM_FAUCET_ADDRESS"`
	RPCURL          string `yaml:"arbitrum_rpc_url" envconfig:"ARBITRUM_RPC_URL"`
	WSURL           string `yaml:"arbitrum_ws_url" envconfig:"ARBITRUM_WS_URL"`
}

func (a *Arbitrum) Init() {
	// Default values connect to a private testnet.
	a.BlockchainID = "412346"
	a.MOMTokenAddress = "0x457fd0Ee3Ce35113ee414994f37eE38518d6E7Ee"
	a.DADTokenAddress = "0xfCa1B6bD67AeF9a9E7047bf7D3949f40E8dde18d"
	a.StakeAddress = "0x18f3FEE919DBc22b5a68401298B01dcd46ab8665"
	a.NFTAddress = "0xbc48cb82903f537614E0309CaF6fe8cEeBa3d174"
	a.NodeAddress = "0x19DcA2dd179A260e9Fe88ea6c821d237EA10Bfd8"
	a.FaucetAddress = "0x9E760F1CddA0694B6156076C60657118CF874289"
	a.MintNFTAmount = "4.20"
	a.MintNFTDeposit = "0x683642c22feDE752415D4793832Ab75EFdF6223c"
	a.RPCURL = "https://bcdev.antst.net:8547"
	a.WSURL = "wss://bcdev.antst.net:8548"
}
