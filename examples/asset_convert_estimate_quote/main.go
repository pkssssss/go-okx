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

	baseCcy := os.Getenv("OKX_BASE_CCY")
	quoteCcy := os.Getenv("OKX_QUOTE_CCY")
	side := os.Getenv("OKX_SIDE") // buy/sell
	rfqSz := os.Getenv("OKX_RFQ_SZ")
	rfqSzCcy := os.Getenv("OKX_RFQ_SZ_CCY")
	if baseCcy == "" || quoteCcy == "" || side == "" || rfqSz == "" || rfqSzCcy == "" {
		log.Fatal("missing env: OKX_BASE_CCY / OKX_QUOTE_CCY / OKX_SIDE / OKX_RFQ_SZ / OKX_RFQ_SZ_CCY")
	}

	clQReqId := os.Getenv("OKX_CL_Q_REQ_ID")
	tag := os.Getenv("OKX_TAG")
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

	svc := c.NewAssetConvertEstimateQuoteService().
		BaseCcy(baseCcy).
		QuoteCcy(quoteCcy).
		Side(side).
		RfqSz(rfqSz).
		RfqSzCcy(rfqSzCcy)
	if clQReqId != "" {
		svc.ClQReqId(clQReqId)
	}
	if tag != "" {
		svc.Tag(tag)
	}

	q, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("quote: quoteId=%s ttlMs=%s quoteTime=%s base=%s quote=%s side=%s rfqSz=%s rfqSzCcy=%s cnvtPx=%s baseSz=%s quoteSz=%s", q.QuoteId, q.TtlMs, q.QuoteTime, q.BaseCcy, q.QuoteCcy, q.Side, q.RfqSz, q.RfqSzCcy, q.CnvtPx, q.BaseSz, q.QuoteSz)
}
