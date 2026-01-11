package main

import (
	"context"
	"encoding/json"
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
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to set rfq maker instrument settings; set OKX_CONFIRM=YES to continue")
	}

	settingsJSON := os.Getenv("OKX_RFQ_MAKER_SETTINGS") // JSON array
	if settingsJSON == "" {
		log.Fatal("missing env: OKX_RFQ_MAKER_SETTINGS (JSON array)")
	}

	var settings []okx.RFQMakerInstrumentSetting
	if err := json.Unmarshal([]byte(settingsJSON), &settings); err != nil {
		log.Fatalf("invalid OKX_RFQ_MAKER_SETTINGS: %v", err)
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

	ack, err := c.NewRFQSetMakerInstrumentSettingsService().Settings(settings).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("set rfq maker instrument settings: result=%v", ack.Result)
}
