package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/examples/internal/redact"
	"github.com/pkssssss/go-okx/v5"
)

func main() {
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	subAcct := os.Getenv("OKX_SUB_ACCT")
	if subAcct == "" {
		log.Fatal("missing env: OKX_SUB_ACCT")
	}
	apiKeyFilter := os.Getenv("OKX_SUB_APIKEY") // optional

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

	svc := c.NewUsersSubaccountAPIKeysService().SubAcct(subAcct)
	if apiKeyFilter != "" {
		svc.APIKey(apiKeyFilter)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("users_subaccount_apikeys: subAcct=%s apiKeyFilter=%s items=%d", subAcct, redact.Mask(apiKeyFilter), len(items))
	for i := 0; i < len(items) && i < 5; i++ {
		x := items[i]
		log.Printf("item[%d]: label=%s apiKey=%s perm=%s ip=%s ts=%s", i, x.Label, redact.Mask(x.APIKey), x.Perm, x.IP, x.TS)
	}
}
