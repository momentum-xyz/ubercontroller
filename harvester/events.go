package harvester

type HarvesterEvent string

const NewBlock HarvesterEvent = "new_block"
const BalanceChange HarvesterEvent = "balance_change"
