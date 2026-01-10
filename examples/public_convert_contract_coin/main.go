package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		instId = "BTC-USD-SWAP"
	}

	sz := os.Getenv("OKX_SZ")
	if sz == "" {
		sz = "0.888"
	}

	svc := okx.NewClient().NewPublicConvertContractCoinService().InstId(instId).Sz(sz)

	if v := os.Getenv("OKX_CONVERT_TYPE"); v != "" {
		svc.Type(v)
	}
	if v := os.Getenv("OKX_PX"); v != "" {
		svc.Px(v)
	}
	if v := os.Getenv("OKX_UNIT"); v != "" {
		svc.Unit(v)
	}
	if v := os.Getenv("OKX_OP_TYPE"); v != "" {
		svc.OpType(v)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instId=%s count=%d\n", instId, len(items))
	if len(items) == 0 {
		return
	}

	it := items[0]
	fmt.Printf("first type=%s instId=%s px=%s sz=%s unit=%s\n", it.Type, it.InstId, it.Px, it.Sz, it.Unit)
}
