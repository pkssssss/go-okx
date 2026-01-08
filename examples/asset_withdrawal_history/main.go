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
	wdId := os.Getenv("OKX_WD_ID")
	clientId := os.Getenv("OKX_WITHDRAWAL_CLIENT_ID")
	txId := os.Getenv("OKX_WITHDRAWAL_TX_ID")
	typ := os.Getenv("OKX_WITHDRAWAL_TYPE")
	state := os.Getenv("OKX_WITHDRAWAL_STATE")
	after := os.Getenv("OKX_WITHDRAWAL_AFTER")
	before := os.Getenv("OKX_WITHDRAWAL_BEFORE")

	var limit *int
	if v := os.Getenv("OKX_WITHDRAWAL_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_WITHDRAWAL_LIMIT: %v", err)
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

	svc := c.NewAssetWithdrawalHistoryService()
	if ccy != "" {
		svc.Ccy(ccy)
	}
	if wdId != "" {
		svc.WdId(wdId)
	}
	if clientId != "" {
		svc.ClientId(clientId)
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

	wds, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("withdrawals=%d", len(wds))
	for i := 0; i < len(wds) && i < 5; i++ {
		w := wds[i]
		log.Printf("wd[%d]: wdId=%s ccy=%s chain=%s amt=%s state=%s ts=%d txId=%s fee=%s feeCcy=%s", i, w.WdId, w.Ccy, w.Chain, w.Amt, w.State, w.TS, w.TxId, w.Fee, w.FeeCcy)
	}
}
