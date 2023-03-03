package main

import (
	"github.com/momentum-xyz/ubercontroller/harvester"
	"github.com/momentum-xyz/ubercontroller/harvester/polkadot_adapter"
)

func main() {
	harv := harvester.NewHarvester()
	a := polkadot_adapter.NewPolkadotAdapter(harv)
	a.Run()
}
