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

	cfg, err := c.NewAccountConfigService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_config: uid=%s mainUid=%s acctLv=%s posMode=%s autoLoan=%v feeType=%s greeksType=%s perm=%s settleCcy=%s settleCcyList=%v", cfg.Uid, cfg.MainUid, cfg.AcctLv, cfg.PosMode, cfg.AutoLoan, cfg.FeeType, cfg.GreeksType, cfg.Perm, cfg.SettleCcy, cfg.SettleCcyList)
}
