package config

type Arbitrum struct {
	BlockchainID          string `yaml:"arbitrum_chain_id" envconfig:"ARBITRUM_CHAIN_ID"`
	MOMTokenAddress       string `yaml:"arbitrum_mom_token_address" envconfig:"ARBITRUM_MOM_TOKEN_ADDRESS"`
	DADTokenAddress       string `yaml:"arbitrum_dad_token_address" envconfig:"ARBITRUM_DAD_TOKEN_ADDRESS"`
	StakeContractAddress  string `yaml:"arbitrum_stake_token_address" envconfig:"ARBITRUM_STAKE_TOKEN_ADDRESS"`
	NFTContractAddress    string `yaml:"arbitrum_nft_contract_address" envconfig:"ARBITRUM_NFT_CONTRACT_ADDRESS"`
	FaucetContractAddress string `yaml:"arbitrum_faucet_contract_address" envconfig:"ARBITRUM_FAUCET_ADDRESS"`
	RPCURL                string `yaml:"arbitrum_rpc_url" envconfig:"ARBITRUM_RPC_URL"`
	WSURL                 string `yaml:"arbitrum_ws_url" envconfig:"ARBITRUM_WS_URL"`
}

func (a *Arbitrum) Init() {
	// Default values connect to a private testnet.
	a.BlockchainID = "412346"
	a.MOMTokenAddress = "0x567d4e8264dC890571D5392fDB9fbd0e3FCBEe56"
	a.DADTokenAddress = "0x0244BbA6fcB25eFed05955C4A1B86A458986D2e0"
	a.StakeContractAddress = "0x047C0A154271498ee718162b718b3D4F464855e0"
	a.NFTContractAddress = "0x97E0B10D89a494Eb5cfFCc72853FB0750BD64AcD"
	a.RPCURL = "https://bcdev.antst.net:8547"
	a.WSURL = "wss://bcdev.antst.net:8548"
}
