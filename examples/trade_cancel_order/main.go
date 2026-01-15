package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to cancel order; set OKX_CONFIRM=YES to continue")
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		log.Fatal("missing env: OKX_INST_ID")
	}

	ordId, hasOrdId := os.LookupEnv("OKX_ORD_ID")
	clOrdId, hasClOrdId := os.LookupEnv("OKX_CL_ORD_ID")
	if (hasOrdId && hasClOrdId) || (!hasOrdId && !hasClOrdId) {
		log.Fatal("require exactly one of OKX_ORD_ID or OKX_CL_ORD_ID")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

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

	svc := c.NewCancelOrderService().InstId(instId)
	if hasOrdId {
		svc.OrdId(ordId)
	} else {
		svc.ClOrdId(clOrdId)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("trade_cancel_order: instId=%s clOrdId=%s ordId=%s sCode=%s sMsg=%s", instId, ack.ClOrdId, ack.OrdId, ack.SCode, ack.SMsg)
}
