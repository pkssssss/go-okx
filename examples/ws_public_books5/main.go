package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pkssssss/go-okx"
)

type bookUpdate struct {
	Action string
	Book   okx.WSOrderBook
}

func main() {
	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		instId = "BTC-USDT"
	}

	c := okx.NewClient()
	ws := c.NewWSPublic()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bookCh := make(chan bookUpdate, 1)
	if err := ws.Start(ctx, func(message []byte) {
		dm, ok, err := okx.WSParseOrderBook(message)
		if err != nil || !ok || dm.Arg.Channel != okx.WSChannelBooks5 || len(dm.Data) == 0 {
			return
		}
		select {
		case bookCh <- bookUpdate{Action: dm.Action, Book: dm.Data[0]}:
		default:
		}
	}, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: okx.WSChannelBooks5, InstId: instId}); err != nil {
		log.Fatal(err)
	}

	select {
	case u := <-bookCh:
		book := u.Book
		id := book.InstId
		if id == "" {
			id = instId
		}
		var bestAsk, bestBid okx.OrderBookLevel
		if len(book.Asks) > 0 {
			bestAsk = book.Asks[0]
		}
		if len(book.Bids) > 0 {
			bestBid = book.Bids[0]
		}
		fmt.Printf("channel=%s instId=%s action=%s bid=%s/%s ask=%s/%s ts=%d checksum=%d prevSeqId=%d seqId=%d\n",
			okx.WSChannelBooks5, id, u.Action, bestBid.Px, bestBid.Sz, bestAsk.Px, bestAsk.Sz, book.TS, book.Checksum, book.PrevSeqId, book.SeqId)
		cancel()
	case <-time.After(10 * time.Second):
		log.Fatal("timeout waiting message")
	}

	<-ws.Done()
}
