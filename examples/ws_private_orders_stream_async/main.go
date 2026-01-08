package main

import (
	"context"
	"log"
	"os"
	"strconv"
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

	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "ANY"
	}
	instId := os.Getenv("OKX_INST_ID")
	instFamily := os.Getenv("OKX_INST_FAMILY")

	timeout := 60 * time.Second
	if v := os.Getenv("OKX_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("invalid OKX_TIMEOUT: %v", err)
		}
		timeout = d
	}

	asyncBuffer := 1024
	if v := os.Getenv("OKX_WS_ASYNC_BUFFER"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			log.Fatalf("invalid OKX_WS_ASYNC_BUFFER: %q", v)
		}
		asyncBuffer = n
	}

	handlerSleep := 500 * time.Millisecond
	if v := os.Getenv("OKX_HANDLER_SLEEP"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("invalid OKX_HANDLER_SLEEP: %v", err)
		}
		handlerSleep = d
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

	orderCh := make(chan okx.TradeOrder, 1)
	errCh := make(chan error, 16)

	ws := c.NewWSPrivate(okx.WithWSTypedHandlerAsync(asyncBuffer))
	ws.OnOrders(func(order okx.TradeOrder) {
		select {
		case orderCh <- order:
		default:
		}

		if handlerSleep > 0 {
			time.Sleep(handlerSleep) // 模拟耗时逻辑：启用 async 后不会阻塞 WS read goroutine。
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := ws.Start(ctx, nil, func(err error) {
		select {
		case errCh <- err:
		default:
		}
	}); err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	subCtx, subCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer subCancel()

	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{
		Channel:    okx.WSChannelOrders,
		InstType:   instType,
		InstId:     instId,
		InstFamily: instFamily,
	}); err != nil {
		log.Fatal(err)
	}
	log.Printf("subscribed: channel=orders instType=%s instId=%s instFamily=%s asyncBuffer=%d handlerSleep=%s", instType, instId, instFamily, asyncBuffer, handlerSleep)

	select {
	case o := <-orderCh:
		log.Printf("order update: instId=%s ordId=%s side=%s ordType=%s px=%s sz=%s state=%s accFillSz=%s avgPx=%s uTime=%d",
			o.InstId, o.OrdId, o.Side, o.OrdType, o.Px, o.Sz, o.State, o.AccFillSz, o.AvgPx, o.UTime,
		)
	case err := <-errCh:
		log.Fatal(err)
	case <-time.After(timeout):
		log.Printf("no order update within %s (place/cancel/amend an order to trigger push)", timeout)
	}

	cancel()
	<-ws.Done()
}
