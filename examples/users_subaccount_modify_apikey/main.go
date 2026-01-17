package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/examples/internal/redact"
	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM_USERS_MODIFY_APIKEY") != "YES" {
		log.Fatal("refuse to modify subaccount apikey without OKX_CONFIRM_USERS_MODIFY_APIKEY=YES")
	}

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
	subAPIKey := os.Getenv("OKX_SUB_APIKEY")
	if subAPIKey == "" {
		log.Fatal("missing env: OKX_SUB_APIKEY")
	}

	label := os.Getenv("OKX_APIKEY_LABEL") // optional
	perm := os.Getenv("OKX_APIKEY_PERM")   // optional
	ip, ipOK := os.LookupEnv("OKX_APIKEY_IP")

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

	svc := c.NewUsersSubaccountModifyAPIKeyService().SubAcct(subAcct).APIKey(subAPIKey)
	if label != "" {
		svc.Label(label)
	}
	if perm != "" {
		svc.Perm(perm)
	}
	if ipOK {
		svc.IP(ip) // allow empty string to unbind IP
	}

	res, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("users_subaccount_modify_apikey: subAcct=%s apiKey=%s label=%s perm=%s ip=%s ts=%s", res.SubAcct, redact.Mask(res.APIKey), res.Label, res.Perm, res.IP, res.TS)
}
