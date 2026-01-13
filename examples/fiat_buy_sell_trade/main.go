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
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to place fiat buy-sell trade; set OKX_CONFIRM=YES to continue")
	}

	clOrdId := os.Getenv("OKX_CL_ORD_ID")
	quoteId := os.Getenv("OKX_QUOTE_ID")
	side := os.Getenv("OKX_SIDE") // buy/sell
	fromCcy := os.Getenv("OKX_FROM_CCY")
	toCcy := os.Getenv("OKX_TO_CCY")
	rfqAmt := os.Getenv("OKX_RFQ_AMT")
	rfqCcy := os.Getenv("OKX_RFQ_CCY")
	paymentMethod := os.Getenv("OKX_PAYMENT_METHOD")
	if paymentMethod == "" {
		paymentMethod = "balance"
	}
	if clOrdId == "" || quoteId == "" || side == "" || fromCcy == "" || toCcy == "" || rfqAmt == "" || rfqCcy == "" {
		log.Fatal("missing env: OKX_CL_ORD_ID / OKX_QUOTE_ID / OKX_SIDE / OKX_FROM_CCY / OKX_TO_CCY / OKX_RFQ_AMT / OKX_RFQ_CCY")
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

	order, err := c.NewFiatBuySellTradeService().
		ClOrdId(clOrdId).
		Side(side).
		FromCcy(fromCcy).
		ToCcy(toCcy).
		RfqAmt(rfqAmt).
		RfqCcy(rfqCcy).
		PaymentMethod(paymentMethod).
		QuoteId(quoteId).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("fiat_buy_sell_trade: ordId=%s clOrdId=%s state=%s cTime=%d", order.OrdId, order.ClOrdId, order.State, order.CTime)
}
