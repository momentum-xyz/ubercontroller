package harvester

type BlockChain struct {
	name                     BCType
	lastProcessedBlockNumber uint64
	rpcURL                   string
}
