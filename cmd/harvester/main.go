package main

import (
	"fmt"
	"time"

	"github.com/momentum-xyz/ubercontroller/harvester"
	"github.com/momentum-xyz/ubercontroller/harvester/ethereum_adapter"
	"github.com/momentum-xyz/ubercontroller/harvester/polkadot_adapter"
)

func main() {
	fmt.Println("Harvester Debugger")

	// ** Harvester
	harv := harvester.NewHarvester()
	var harvForClient harvester.HarvesterAPI
	harvForClient = harv
	var harvForAdapter harvester.BCAdapterAPI
	harvForAdapter = harv

	// ** Ethereum Adapter
	ethereumAdapter := ethereum_adapter.NewEthereumAdapter(harvForAdapter)
	go ethereumAdapter.Run()

	// ** Polkadot Adapter
	polkadotAdapter := polkadot_adapter.NewPolkadotAdapter(harvForAdapter)
	go polkadotAdapter.Run()

	// ** Harvester Clients
	testHandler1 := testHandler1
	ptrTestHandler1 := &testHandler1
	harvForClient.Subscribe(harvester.Ethereum, harvester.NewBlock, ptrTestHandler1)
	harvForClient.Subscribe(harvester.Polkadot, harvester.NewBlock, ptrTestHandler1)

	testHandler2 := testHandler2
	ptrTestHandler2 := &testHandler2
	harvForClient.Subscribe(harvester.Ethereum, harvester.NewBlock, ptrTestHandler2)

	time.Sleep(time.Second * 30)
	harvForClient.Unsubscribe(harvester.Ethereum, harvester.NewBlock, ptrTestHandler2)

	time.Sleep(time.Second * 50)
}

func testHandler1(p any) {
	fmt.Printf("testHandler1: %+v \n", p)
}

func testHandler2(p any) {
	fmt.Printf("testHandler2: %+v \n", p)
}
