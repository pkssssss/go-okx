package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	transId := os.Getenv("OKX_TRANS_ID")
	clientId := os.Getenv("OKX_CLIENT_ID")
	if transId == "" && clientId == "" {
		log.Fatal("missing env: OKX_TRANS_ID or OKX_CLIENT_ID")
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

	svc := c.NewAssetTransferStateService()
	if transId != "" {
		svc.TransId(transId)
	} else {
		svc.ClientId(clientId)
	}
	if v := os.Getenv("OKX_TYPE"); v != "" {
		svc.Type(v)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("asset_transfer_state: items=%d", len(items))
	if len(items) > 0 {
		x := items[0]
		log.Printf("first: transId=%s clientId=%s ccy=%s amt=%s state=%s from=%s to=%s", x.TransId, x.ClientId, x.Ccy, x.Amt, x.State, x.From, x.To)
	}
}
