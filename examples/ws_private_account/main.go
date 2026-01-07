package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/pkssssss/go-okx"
)

func main() {
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	demo := os.Getenv("OKX_DEMO") == "1"
	ccy := os.Getenv("OKX_CCY")
	extraParams := os.Getenv("OKX_WS_EXTRA_PARAMS")

	c := okx.NewClient(
		okx.WithCredentials(okx.Credentials{
			APIKey:     apiKey,
			SecretKey:  secretKey,
			Passphrase: passphrase,
		}),
		okx.WithDemoTrading(demo),
	)

	if _, err := c.SyncTime(context.Background()); err != nil {
		log.Fatal(err)
	}

	ws := c.NewWSPrivate()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := ws.Start(ctx, func(message []byte) {
		dm, ok, err := okx.WSParseAccount(message)
		if err != nil || !ok {
			return
		}

		totalEq := ""
		availEq := ""
		if len(dm.Data) > 0 {
			totalEq = dm.Data[0].TotalEq
			availEq = dm.Data[0].AvailEq
		}
		log.Printf("account push: eventType=%s curPage=%d lastPage=%v totalEq=%s availEq=%s items=%d", dm.EventType, dm.CurPage, dm.LastPage, totalEq, availEq, len(dm.Data))
		cancel()
	}, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, subCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer subCancel()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{
		Channel:     okx.WSChannelAccount,
		Ccy:         ccy,
		ExtraParams: extraParams,
	}); err != nil {
		log.Fatal(err)
	}

	select {
	case <-ws.Done():
	case <-time.After(10 * time.Second):
		log.Printf("timeout waiting account push")
		ws.Close()
		<-ws.Done()
	}
}
