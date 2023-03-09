package harvester

type BCAdapter interface {
	GetLastBlockNumber() (uint64, error)
}
