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

	svc := c.NewUsersSubaccountListService()

	if s := os.Getenv("OKX_ENABLE"); s != "" {
		b, err := strconv.ParseBool(s)
		if err != nil {
			log.Fatalf("invalid env OKX_ENABLE=%q: %v", s, err)
		}
		svc.Enable(b)
	}
	if s := os.Getenv("OKX_SUB_ACCT"); s != "" {
		svc.SubAcct(s)
	}
	if s := os.Getenv("OKX_AFTER"); s != "" {
		svc.After(s)
	}
	if s := os.Getenv("OKX_BEFORE"); s != "" {
		svc.Before(s)
	}
	if s := os.Getenv("OKX_LIMIT"); s != "" {
		svc.Limit(s)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("users_subaccount_list: items=%d", len(items))
	for i := 0; i < len(items) && i < 5; i++ {
		x := items[i]
		log.Printf("item[%d]: subAcct=%s type=%s enable=%t canTransOut=%t uid=%s ts=%s", i, x.SubAcct, x.Type, x.Enable, x.CanTransOut, x.UID, x.TS)
	}
}
