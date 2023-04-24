package config

type Arbitrum struct {
	ArbitrumMOMTokenAddress      string `yaml:"arbitrum_mom_token_address" envconfig:"ARBITRUM_MOM_TOKEN_ADDRESS"`
	ArbitrumDADTokenAddress      string `yaml:"arbitrum_dad_token_address" envconfig:"ARBITRUM_DAD_TOKEN_ADDRESS"`
	ArbitrumStakeContractAddress string `yaml:"arbitrum_stake_token_address" envconfig:"ARBITRUM_STAKE_TOKEN_ADDRESS"`
	ArbitrumRPCURL               string `yaml:"arbitrum_rpc_url" envconfig:"ARBITRUM_RPC_URL"`
	ArbitrumWSURL                string `yaml:"arbitrum_ws_url" envconfig:"ARBITRUM_WS_URL"`
}

func (a *Arbitrum) Init() {
	a.ArbitrumMOMTokenAddress = "0x310c2B16c304109f32BABB5f47cC562813765744"
	a.ArbitrumDADTokenAddress = "0xB647d3a893E7e0534827B5E795d3BF7cb80FF16f"
	a.ArbitrumStakeContractAddress = "0xC4497d6c0f94dc427cE0B8F825c91F25e2845B91"
	a.ArbitrumRPCURL = "https://bcdev.antst.net:8547"
	a.ArbitrumWSURL = "wss://bcdev.antst.net:8548"
}
