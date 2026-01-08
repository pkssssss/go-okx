package main

import (
	"context"
	"log"
	"os"
	"strconv"

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

	ccy := os.Getenv("OKX_CCY")
	depId := os.Getenv("OKX_DEPOSIT_DEP_ID")
	fromWdId := os.Getenv("OKX_DEPOSIT_FROM_WD_ID")
	txId := os.Getenv("OKX_DEPOSIT_TX_ID")
	typ := os.Getenv("OKX_DEPOSIT_TYPE")
	state := os.Getenv("OKX_DEPOSIT_STATE")
	after := os.Getenv("OKX_DEPOSIT_AFTER")
	before := os.Getenv("OKX_DEPOSIT_BEFORE")

	var limit *int
	if v := os.Getenv("OKX_DEPOSIT_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_DEPOSIT_LIMIT: %v", err)
		}
		limit = &n
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

	svc := c.NewAssetDepositHistoryService()
	if ccy != "" {
		svc.Ccy(ccy)
	}
	if depId != "" {
		svc.DepId(depId)
	}
	if fromWdId != "" {
		svc.FromWdId(fromWdId)
	}
	if txId != "" {
		svc.TxId(txId)
	}
	if typ != "" {
		svc.Type(typ)
	}
	if state != "" {
		svc.State(state)
	}
	if after != "" {
		svc.After(after)
	}
	if before != "" {
		svc.Before(before)
	}
	if limit != nil {
		svc.Limit(*limit)
	}

	deps, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("deposits=%d", len(deps))
	for i := 0; i < len(deps) && i < 5; i++ {
		d := deps[i]
		log.Printf("dep[%d]: depId=%s ccy=%s chain=%s amt=%s state=%s ts=%d txId=%s", i, d.DepId, d.Ccy, d.Chain, d.Amt, d.State, d.TS, d.TxId)
	}
}
