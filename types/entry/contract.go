package entry

type Contract struct {
	ContractID []byte `db:"contract_id"`
	Name       string `db:"name"`
}
