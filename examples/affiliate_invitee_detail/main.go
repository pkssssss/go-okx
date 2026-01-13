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

	uid := os.Getenv("OKX_UID")
	if uid == "" {
		log.Fatal("missing env: OKX_UID")
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

	res, err := c.NewAffiliateInviteeDetailService().UID(uid).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("affiliate_invitee_detail: uid=%s level=%s rebateRate=%s totalCommission=%s joinTime=%d", uid, res.InviteeLevel, res.InviteeRebateRate, res.TotalCommission, res.JoinTime)
}
