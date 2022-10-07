package types

const (
	ContextLoggerKey = "logger"
)

type ConfigAuthProviderType string

const (
	ConfigAuthUnknownProviderType ConfigAuthProviderType = ""
	ConfigAuthWeb3ProviderType    ConfigAuthProviderType = "web3"
)
