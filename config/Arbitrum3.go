package config

type Arbitrum3 struct {
	RPCURL string `yaml:"arbitrum_rpc_url" envconfig:"ARBITRUM_RPC_URL"`
}

func (a *Arbitrum3) Init() {
	a.RPCURL = "https://bcdev.antst.net:8547"
}
