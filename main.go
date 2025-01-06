package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"

	"icylight/uniswap/config"
	"icylight/uniswap/controllers/transaction"
	"icylight/uniswap/service"
)

const (
	subStr = `{"jsonrpc": "2.0", "method": "eth_subscribe","params": ["alchemy_minedTransactions", {"addresses": [{"to": "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"}]}],"id": 1}`
)

type inputData struct {
	Params struct {
		Transaction struct {
			InputData string `json:"input"`
		} `json:"transaction"`
	} `json:"params"`
}

func main() {
	config.Load("./config.yml")
	service.Init()

	c, _, err := websocket.DefaultDialer.Dial(config.Conf.Eth.WS, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			go transaction.GetPairs(message)
		}
	}()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	err = c.WriteMessage(websocket.TextMessage, []byte(subStr))
	if err != nil {
		log.Fatal("write:", err)
	}
	log.Println("subscribe success")

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
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
