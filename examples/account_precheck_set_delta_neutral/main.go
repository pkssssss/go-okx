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

	stgyType := os.Getenv("OKX_STGY_TYPE") // 0/1
	if stgyType == "" {
		log.Fatal("missing env: OKX_STGY_TYPE (0 / 1)")
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

	res, err := c.NewAccountPrecheckSetDeltaNeutralService().StgyType(stgyType).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_precheck_set_delta_neutral: stgyType=%s unmatched=%d", stgyType, len(res.UnmatchedInfoCheck))
	for _, it := range res.UnmatchedInfoCheck {
		log.Printf("type=%s deltaLever=%s pos=%d ord=%d", it.Type, it.DeltaLever, len(it.PosList), len(it.OrdList))
	}
}
