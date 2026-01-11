package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	c := okx.NewClient()
	ws := c.NewWSPublic()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gotCh := make(chan okx.SystemStatus, 1)
	ws.OnStatus(func(status okx.SystemStatus) {
		select {
		case gotCh <- status:
		default:
		}
	})

	if err := ws.Start(ctx, nil, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelStatus}); err != nil {
		log.Fatal(err)
	}

	select {
	case status := <-gotCh:
		fmt.Printf("state=%s begin=%d end=%d title=%q serviceType=%s env=%s ts=%d\n",
			status.State, status.Begin, status.End, status.Title, status.ServiceType, status.Env, status.TS)
		cancel()
	case <-time.After(15 * time.Second):
		log.Fatal("timeout waiting message")
	}

	<-ws.Done()
}
