package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM_USERS_SET_TRANSFER_OUT") != "YES" {
		log.Fatal("refuse to set transfer out permission without OKX_CONFIRM_USERS_SET_TRANSFER_OUT=YES")
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	subAcct := os.Getenv("OKX_SUB_ACCT")
	if subAcct == "" {
		log.Fatal("missing env: OKX_SUB_ACCT (supports multiple, comma separated)")
	}

	canTransOutRaw := os.Getenv("OKX_CAN_TRANS_OUT")
	if canTransOutRaw == "" {
		log.Fatal("missing env: OKX_CAN_TRANS_OUT (true/false)")
	}
	canTransOut, err := strconv.ParseBool(canTransOutRaw)
	if err != nil {
		log.Fatalf("invalid env OKX_CAN_TRANS_OUT=%q: %v", canTransOutRaw, err)
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

	items, err := c.NewUsersSubaccountSetTransferOutService().SubAcct(subAcct).CanTransOut(canTransOut).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("users_subaccount_set_transfer_out: subAcct=%s canTransOut=%t items=%d", subAcct, canTransOut, len(items))
	for i := 0; i < len(items) && i < 5; i++ {
		x := items[i]
		log.Printf("item[%d]: subAcct=%s canTransOut=%t", i, x.SubAcct, x.CanTransOut)
	}
}
