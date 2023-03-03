package main

import (
	"fmt"
	"time"

	"github.com/momentum-xyz/ubercontroller/harvester"
	"github.com/momentum-xyz/ubercontroller/harvester/etherium_adapter"
)

func main() {
	fmt.Println("Harvester Debugger")

	// ** Harvester
	harv := harvester.NewHarvester()
	var harvForClient harvester.HarvesterAPI
	harvForClient = harv
	var harvForAdapter harvester.BCAdapterAPI
	harvForAdapter = harv

	// ** Adapter
	etheriumAdapter := etherium_adapter.NewEtheriumAdapter(harvForAdapter)
	go etheriumAdapter.Run()

	// ** Harvester Clients
	testHandler1 := testHandler1
	ptrTestHandler1 := &testHandler1
	harvForClient.Subscribe(harvester.Etherium, harvester.NewBlock, ptrTestHandler1)

	testHandler2 := testHandler2
	ptrTestHandler2 := &testHandler2
	harvForClient.Subscribe(harvester.Etherium, harvester.NewBlock, ptrTestHandler2)

	time.Sleep(time.Second * 20)
	harvForClient.Unsubscribe(harvester.Etherium, harvester.NewBlock, ptrTestHandler2)

	time.Sleep(time.Second * 20)
}

func testHandler1(p any) {
	fmt.Printf("testHandler1: %+v \n", p)
}

func testHandler2(p any) {
	fmt.Printf("testHandler2: %+v \n", p)
}
