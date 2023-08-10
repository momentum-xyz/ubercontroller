package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/harvester3/arbitrum_nova_adapter3"
	helper "github.com/momentum-xyz/ubercontroller/harvester3/cmd"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	cfg.Arbitrum3.RPCURL = "https://nova.arbitrum.io/rpc"

	logger := helper.GetZapLogger()
	sugaredLogger := logger.Sugar()

	a := arbitrum_nova_adapter3.NewArbitrumNovaAdapter(&cfg.Arbitrum3, sugaredLogger)

	a.Run()

	n, err := a.GetLastBlockNumber()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Last Block: %+v \n", n)

	mom := common.HexToAddress("0x457fd0Ee3Ce35113ee414994f37eE38518d6E7Ee")
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

	//mom = common.HexToAddress("0x0C270A47D5B00bb8db42ed39fa7D6152496944ca")
	//dad := common.HexToAddress("0x11817050402d2bb1418753ca398fdB3A3bc7CfEA")
	//_ = mom
	//_ = dad

	wAdd := common.HexToAddress("0xAdd2e75c298F34E4d66fBbD4e056DA31502Da5B0")
	_ = wAdd

	//wAnton := common.HexToAddress("0x83FfD8c86e7cC10544403220d857c66bF6CdF8B8")

	b, _, err := a.GetTokenBalance(&mom, &w68, n)
	fmt.Println(n)
	fmt.Println(b.String())

	time.Sleep(time.Second * 300)
}
