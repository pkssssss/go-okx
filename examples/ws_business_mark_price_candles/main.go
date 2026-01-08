package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		instId = "BTC-USDT-SWAP"
	}

	bar := os.Getenv("OKX_CANDLE_BAR")
	if bar == "" {
		bar = "1m"
	}

	c := okx.NewClient()
	ws := c.NewWSBusiness()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	candleCh := make(chan okx.PriceCandle, 1)
	channelCh := make(chan string, 1)

	if err := ws.Start(ctx, func(message []byte) {
		dm, ok, err := okx.WSParseMarkPriceCandles(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		select {
		case channelCh <- dm.Arg.Channel:
		default:
		}
		select {
		case candleCh <- dm.Data[0]:
		default:
		}
	}, func(err error) {
		log.Printf("ws error: %v", err)
	}); err != nil {
		log.Fatal(err)
	}

	subCtx, cancelSub := context.WithTimeout(ctx, 10*time.Second)
	defer cancelSub()

	channel := okx.WSMarkPriceCandleChannel(bar)
	if err := ws.SubscribeAndWait(subCtx, okx.WSArg{Channel: channel, InstId: instId}); err != nil {
		log.Fatal(err)
	}

	select {
	case candle := <-candleCh:
		gotChannel := channel
		select {
		case gotChannel = <-channelCh:
		default:
		}
		fmt.Printf("channel=%s instId=%s ts=%d o=%s h=%s l=%s c=%s confirm=%s\n",
			gotChannel, instId, candle.TS, candle.Open, candle.High, candle.Low, candle.Close, candle.Confirm)
		cancel()
	case <-time.After(15 * time.Second):
		log.Fatal("timeout waiting message")
	}

	<-ws.Done()
}
