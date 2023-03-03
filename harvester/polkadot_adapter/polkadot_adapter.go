package polkadot_adapter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"github.com/momentum-xyz/ubercontroller/harvester"
)

type PolkadotAdapter struct {
	harv   harvester.BCAdapterAPI
	blocks []harvester.BCBlock
}

type JSONRPCResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Subscription string `json:"subscription"`
		Result       struct {
			ParentHash     string `json:"parentHash"`
			Number         string `json:"number"`
			StateRoot      string `json:"stateRoot"`
			ExtrinsicsRoot string `json:"extrinsicsRoot"`
			Digest         struct {
				Logs []string `json:"logs"`
			} `json:"digest"`
		} `json:"result"`
	} `json:"params"`
}

func NewPolkadotAdapter(harv harvester.BCAdapterAPI) *PolkadotAdapter {
	return &PolkadotAdapter{
		harv:   harv,
		blocks: make([]harvester.BCBlock, 0),
	}
}

func (pa *PolkadotAdapter) Run() {
	//Create Message Out
	messageOut := make(chan string)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "wss", Host: "drive.antst.net:19947", Path: "/"}
	log.Printf("connecting to Polkadot Block Chain: %s", u.String())
	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("dial:", err)
	}

	//When the program closes close the connection
	defer conn.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				//log.Println("read:", err)
				return
			}
			//log.Printf("recv: %s", message)

			//if string(message) == "Connected" {
			//	log.Printf("Send Sub Details: %s", message)
			//	payload := `{"id":12,"jsonrpc":"2.0","method":"chain_subscribeNewHead","params":[]}`
			//	messageOut <- payload
			//}

			response := make(map[string]any)
			if err := json.Unmarshal(message, &response); err != nil {
				log.Println("Error:", err)
			}

			var id uint64 = 0
			idFloat, ok := response["id"].(float64)
			if ok == true {
				id = uint64(idFloat)
			}

			method, ok := response["method"]
			_ = method
			_ = ok

			if ok == true && method == "chain_finalizedHead" {
				params := response["params"].(map[string]any)
				result := params["result"].(map[string]any)
				numberString := result["number"].(string)
				numberString = strings.Replace(numberString, "0x", "", 1)
				number, err := strconv.ParseUint(numberString, 16, 64)
				if err != nil {
					fmt.Println("Error:", err)
				}
				s := strconv.FormatUint(number, 10)
				payload := `{"id":` + s + `,"jsonrpc":"2.0","method":"chain_getBlockHash","params":[` + s + `]}`
				messageOut <- payload
			}

			hashString, ok := response["result"].(string)
			if ok == true {
				if id > 12 {
					pa.harv.OnNewBlock(harvester.Polkadot, &harvester.BCBlock{
						Hash:   hashString,
						Number: id,
					})
				}
			}
		}

	}()

	time.AfterFunc(time.Second*3, func() {
		//payload := `{"id":12,"jsonrpc":"2.0","method":"chain_subscribeNewHead","params":[]}`
		//payload := `{"id":12,"jsonrpc":"2.0","method":"chain_getBlockHash","params":[117048]}`
		payload := `{"id":12,"jsonrpc":"2.0","method":"chain_subscribeFinalizedHeads","params":[]}`
		messageOut <- payload
	})

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case m := <-messageOut:
			//log.Printf("Send Message %s", m)
			err := conn.WriteMessage(websocket.TextMessage, []byte(m))
			if err != nil {
				//log.Println("write:", err)
				return
			}
		//case t := <-ticker.C:
		//	err := conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
		//	if err != nil {
		//		log.Println("write:", err)
		//		return
		//	}
		case <-interrupt:
			log.Println("interrupt")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
