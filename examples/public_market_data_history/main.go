package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	module := os.Getenv("OKX_MODULE")
	if module == "" {
		module = "1"
	}

	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SWAP"
	}

	dateAggrType := os.Getenv("OKX_DATE_AGGR_TYPE")
	if dateAggrType == "" {
		dateAggrType = "daily"
	}

	begin := os.Getenv("OKX_BEGIN")
	end := os.Getenv("OKX_END")
	if begin == "" || end == "" {
		now := time.Now().UTC()
		// 默认查询近 2 天范围（按 OKX 规则仅使用日期部分）。
		endMs := now.Add(-48 * time.Hour).UnixMilli()
		beginMs := now.Add(-72 * time.Hour).UnixMilli()
		begin = strconv.FormatInt(beginMs, 10)
		end = strconv.FormatInt(endMs, 10)
	}

	svc := okx.NewClient().
		NewPublicMarketDataHistoryService().
		Module(module).
		InstType(instType).
		DateAggrType(dateAggrType).
		Begin(begin).
		End(end)

	if instType == "SPOT" {
		instIdList := os.Getenv("OKX_INST_ID")
		if instIdList == "" {
			instIdList = "BTC-USDT"
		}
		svc.InstIdList(instIdList)
	} else {
		instFamilyList := os.Getenv("OKX_INST_FAMILY")
		if instFamilyList == "" {
			instFamilyList = "BTC-USDT"
		}
		svc.InstFamilyList(instFamilyList)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("module=%s instType=%s dateAggrType=%s count=%d\n", module, instType, dateAggrType, len(items))
	if len(items) == 0 {
		return
	}

	it := items[0]
	fmt.Printf("first ts=%d totalSizeMB=%s details=%d\n", it.TS, it.TotalSizeMB, len(it.Details))
	if len(it.Details) == 0 || len(it.Details[0].GroupDetails) == 0 {
		return
	}
	f := it.Details[0].GroupDetails[0]
	fmt.Printf("first file=%s sizeMB=%s url=%s\n", f.Filename, f.SizeMB, f.URL)
}
