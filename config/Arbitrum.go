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
	a.ArbitrumMOMTokenAddress = "0x310c2B16c304109f32BABB5f47cC562813765744"
	a.ArbitrumDADTokenAddress = "0x5B328f060Ac623A8e9EB9C6F5A7947F3Cdd82b37"
	a.ArbitrumStakeContractAddress = "0x3A85e361917180567F6a0fb8c68B2b5065126aCA"
	a.ArbitrumNFTContractAddress = "0xa662897d53ff3cFA9cb44Ed635dB6e152C68C677"
	a.ArbitrumRPCURL = "https://bcdev.antst.net:8547"
	a.ArbitrumWSURL = "wss://bcdev.antst.net:8548"
}
