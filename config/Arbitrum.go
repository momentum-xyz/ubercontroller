package config

type Arbitrum struct {
	ArbitrumMOMTokenAddress      string `yaml:"arbitrum_mom_token_address" envconfig:"ARBITRUM_MOM_TOKEN_ADDRESS"`
	ArbitrumDADTokenAddress      string `yaml:"arbitrum_dad_token_address" envconfig:"ARBITRUM_DAD_TOKEN_ADDRESS"`
	ArbitrumStakeContractAddress string `yaml:"arbitrum_stake_token_address" envconfig:"ARBITRUM_STAKE_TOKEN_ADDRESS"`
	ArbitrumNFTContractAddress   string `yaml:"arbitrum_nft_contract_address" envconfig:"ARBITRUM_NFT_CONTRACT_ADDRESS"`
	ArbitrumRPCURL               string `yaml:"arbitrum_rpc_url" envconfig:"ARBITRUM_RPC_URL"`
	ArbitrumWSURL                string `yaml:"arbitrum_ws_url" envconfig:"ARBITRUM_WS_URL"`
}

func (a *Arbitrum) Init() {
	a.ArbitrumMOMTokenAddress = "0x567d4e8264dC890571D5392fDB9fbd0e3FCBEe56"
	a.ArbitrumDADTokenAddress = "0x0244BbA6fcB25eFed05955C4A1B86A458986D2e0"
	a.ArbitrumStakeContractAddress = "0xb187f16656C30580bB0B0b797DaDB9CFab766156"
	a.ArbitrumNFTContractAddress = "0x97E0B10D89a494Eb5cfFCc72853FB0750BD64AcD"
	a.ArbitrumRPCURL = "https://bcdev.antst.net:8547"
	//a.ArbitrumRPCURL = "https://nova.arbitrum.io/rpc"
	a.ArbitrumWSURL = "wss://bcdev.antst.net:8548"
}
