package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	demo := os.Getenv("OKX_DEMO") == "1"
	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "ANY"
	}

	timeout := 60 * time.Second
	if v := os.Getenv("OKX_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("invalid OKX_TIMEOUT: %v", err)
		}
		timeout = d
	}

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

	ws := c.NewWSBusinessPrivate()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ws.OnGridOrdersContract(func(order okx.WSGridOrder) {
		fmt.Printf("channel=%s instType=%s instId=%s algoId=%s algoOrdType=%s state=%s gridNum=%s runPx=%s pnlRatio=%s cTime=%d uTime=%d\n",
			okx.WSChannelGridOrdersContract, order.InstType, order.InstId, order.AlgoId, order.AlgoOrdType, order.State, order.GridNum, order.RunPx, order.PnlRatio, order.CTime, order.UTime)
	})

	if err := ws.Start(ctx, nil, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{
		Channel:  okx.WSChannelGridOrdersContract,
		InstType: instType,
	}); err != nil {
		log.Fatal(err)
	}

	log.Printf("subscribed: channel=%s instType=%s", okx.WSChannelGridOrdersContract, instType)

	select {
	case <-ws.Done():
	case <-time.After(timeout):
		log.Printf("timeout after %s", timeout)
		ws.Close()
		<-ws.Done()
	}
}
