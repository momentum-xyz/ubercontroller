package main

import (
	"fmt"
	"github.com/momentum-xyz/ubercontroller/harvester2"
	"time"
)

func main() {
	c := harvester2.NewCallbacks()
	fmt.Println(c)

	vListener := listener
	pListener := &vListener
	c.Add("arbitrum_nova", "event", pListener)
	c.Trigger("arbitrum_nova", "event", "TTT")
	time.Sleep(time.Second)
}

func listener(p any) {
	fmt.Println("listener")
	fmt.Println(p)
}
