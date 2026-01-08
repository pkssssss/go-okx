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

	demo := os.Getenv("OKX_DEMO") == "1"

	wdId := os.Getenv("OKX_WD_ID")
	txId := os.Getenv("OKX_TX_ID")

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

	svc := c.NewAssetDepositWithdrawStatusService()
	if wdId != "" {
		svc.WdId(wdId)
	} else {
		ccy := os.Getenv("OKX_CCY")
		chain := os.Getenv("OKX_CHAIN")
		to := os.Getenv("OKX_DEPOSIT_TO")
		if txId == "" || ccy == "" || chain == "" || to == "" {
			log.Fatal("missing env for deposit query: OKX_TX_ID / OKX_CCY / OKX_CHAIN / OKX_DEPOSIT_TO (or set OKX_WD_ID for withdrawal query)")
		}
		svc.TxId(txId).Ccy(ccy).Chain(chain).To(to)
	}

	st, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("status=%d", len(st))
	for i := 0; i < len(st) && i < 5; i++ {
		x := st[i]
		log.Printf("item[%d]: wdId=%s txId=%s estCompleteTime=%s state=%s", i, x.WdId, x.TxId, x.EstCompleteTime, x.State)
	}
}
