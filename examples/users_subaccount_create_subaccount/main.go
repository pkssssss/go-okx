package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM_USERS_CREATE_SUBACCOUNT") != "YES" {
		log.Fatal("refuse to create subaccount without OKX_CONFIRM_USERS_CREATE_SUBACCOUNT=YES")
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	subAcct := os.Getenv("OKX_NEW_SUB_ACCT")
	if subAcct == "" {
		log.Fatal("missing env: OKX_NEW_SUB_ACCT")
	}
	typ := os.Getenv("OKX_SUB_ACCT_TYPE")
	if typ == "" {
		log.Fatal("missing env: OKX_SUB_ACCT_TYPE")
	}
	label := os.Getenv("OKX_LABEL")
	if label == "" {
		log.Fatal("missing env: OKX_LABEL")
	}
	pwd := os.Getenv("OKX_PWD") // optional; KYB accounts required

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

	svc := c.NewUsersSubaccountCreateSubaccountService().
		SubAcct(subAcct).
		Type(typ).
		Label(label)
	if pwd != "" {
		svc.Pwd(pwd)
	}

	res, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("users_subaccount_create_subaccount: subAcct=%s uid=%s ts=%s label=%s", res.SubAcct, res.UID, res.TS, res.Label)
}
