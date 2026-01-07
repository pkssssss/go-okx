package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkssssss/go-okx"
)

func main() {
	c := okx.NewClient()

	ws := c.NewWSPublic()
	_ = ws.Subscribe(okx.WSArg{
		Channel: "tickers",
		InstId:  "BTC-USDT",
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	msgCh := make(chan []byte, 1)
	if err := ws.Start(ctx, func(message []byte) {
		select {
		case msgCh <- message:
		default:
		}
	}, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	select {
	case msg := <-msgCh:
		fmt.Println(string(msg))
		cancel()
	case <-time.After(10 * time.Second):
		log.Fatal("timeout waiting message")
	}

	<-ws.Done()
}
