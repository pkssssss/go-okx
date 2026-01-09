package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM_USERS_CREATE_APIKEY") != "YES" {
		log.Fatal("refuse to create subaccount apikey without OKX_CONFIRM_USERS_CREATE_APIKEY=YES")
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
	label := os.Getenv("OKX_APIKEY_LABEL")
	if label == "" {
		log.Fatal("missing env: OKX_APIKEY_LABEL")
	}
	apiKeyPassphrase := os.Getenv("OKX_APIKEY_PASSPHRASE")
	if apiKeyPassphrase == "" {
		log.Fatal("missing env: OKX_APIKEY_PASSPHRASE")
	}
	perm := os.Getenv("OKX_APIKEY_PERM") // optional
	ip := os.Getenv("OKX_APIKEY_IP")     // optional

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

	svc := c.NewUsersSubaccountCreateAPIKeyService().
		SubAcct(subAcct).
		Label(label).
		Passphrase(apiKeyPassphrase)
	if perm != "" {
		svc.Perm(perm)
	}
	if ip != "" {
		svc.IP(ip)
	}

	res, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("users_subaccount_create_apikey: subAcct=%s label=%s apiKey=%s perm=%s ip=%s ts=%s", res.SubAcct, res.Label, res.APIKey, res.Perm, res.IP, res.TS)
	log.Printf("users_subaccount_create_apikey: secretKey=%s passphrase=%s", res.SecretKey, res.Passphrase)
}
