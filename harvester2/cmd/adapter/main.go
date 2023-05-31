package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/harvester2"
	"github.com/momentum-xyz/ubercontroller/harvester2/arbitrum_nova_adapter2"
	"log"
	"time"
)

func main() {
	cfg := config.GetConfig()
	a := arbitrum_nova_adapter2.NewArbitrumNovaAdapter(cfg)

	a.Run()

	n, err := a.GetLastBlockNumber()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Last Block: %+v \n", n)

	var l harvester2.AdapterListener
	l = func(blockNumber uint64) {

	}

	a.RegisterNewBlockListener(l)

	contracts := []common.Address{
		//common.HexToAddress("0x7F85fB7f42A0c0D40431cc0f7DFDf88be6495e67"),
		//common.HexToAddress("0x567d4e8264dC890571D5392fDB9fbd0e3FCBEe56"), //mom
		common.HexToAddress("0xDA10009cBd5D07dd0CeCc66161FC93D7c9000da1"), //mom
	}
	_ = contracts
	//diffs, err := a.GetLogs(0, 12, contracts)
	c := arbitrum_nova_adapter2.NewContracts(&cfg.Arbitrum)

	for k, v := range contracts {
		fmt.Println(k, v)
	}
	for k, v := range c.AllAddresses {
		fmt.Println(k, v)
	}
	logs, err := a.GetLogs(0, 8817433, contracts)

	for _, log := range logs {
		switch log.(type) {
		case *harvester2.TransferERC20Log:
			l := log.(*harvester2.TransferERC20Log)
			fmt.Printf("%s %s %s \n", l.From, l.To, l.Value)
			//fmt.Println(log.(*harvester.TransferERC20Log).Value)
		}
	}

	for _, log := range logs {
		switch log.(type) {
		case *harvester2.StakeLog:
			l := log.(*harvester2.StakeLog)
			fmt.Printf("  stake: %s %s %s %s %s \n", l.TxHash, l.UserWallet, l.OdysseyID, l.AmountStaked, l.TotalStaked)
			//fmt.Println(log.(*harvester.TransferERC20Log).Value)
		case *harvester2.UnstakeLog:
			l := log.(*harvester2.UnstakeLog)
			fmt.Printf("unstake: %s %s %s %s \n", l.UserWallet, l.OdysseyID, l.AmountUnstaked, l.TotalStaked)
			//fmt.Println(log.(*harvester.TransferERC20Log).Value)
		}
	}

	mom := common.HexToAddress("0x567d4e8264dC890571D5392fDB9fbd0e3FCBEe56")
	nft := common.HexToAddress("0x97E0B10D89a494Eb5cfFCc72853FB0750BD64AcD")
	stake := common.HexToAddress("0xe9C6d7Cd04614Dde6Ca68B62E6fbf23AC2ECe2F8")
	_ = mom
	_ = nft
	_ = stake

	w04 := common.HexToAddress("0xA058Aa2fCf33993e17D074E6843202E7C94bf267")
	w78 := common.HexToAddress("0x78B00B17E7e5619113A4e922BC3c8cb290355043")
	w68 := common.HexToAddress("0x683642c22feDE752415D4793832Ab75EFdF6223c")
	_ = w04
	_ = w78
	_ = w68

	ids, err := a.GetNFTBalance(1000, &w78, &nft)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ids)

	b, err := a.GetBalance(&w78, &mom, n)
	fmt.Println(b.String())

	s, err := a.GetStakeBalance(1000, &w68, &stake)
	if err != nil {
		log.Fatal(err)
	}

	for odyssey, stakes := range s {
		fmt.Println(odyssey, stakes[0].String(), stakes[1].String(), stakes[2].String())
	}

	time.Sleep(time.Second * 300)
}