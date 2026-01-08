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

	month := os.Getenv("OKX_MONTH")
	if month == "" {
		log.Fatal("missing env: OKX_MONTH (Jan/Feb/.../Dec)")
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

	stmts, err := c.NewAssetMonthlyStatementService().Month(month).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("monthly_statement: month=%s items=%d", month, len(stmts))
	for i := 0; i < len(stmts) && i < 5; i++ {
		x := stmts[i]
		log.Printf("stmt[%d]: state=%s ts=%d href=%s", i, x.State, x.TS, x.FileHref)
	}
}
