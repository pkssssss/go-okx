package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SWAP"
	}

	sortType := os.Getenv("OKX_SORT_TYPE")
	if sortType == "" {
		sortType = "overview"
	}

	svc := okx.NewClient().NewCopyTradingPublicLeadTradersService().
		InstType(instType).
		SortType(sortType)

	if v := os.Getenv("OKX_STATE"); v != "" {
		svc.State(v)
	}
	if v := os.Getenv("OKX_DATA_VER"); v != "" {
		svc.DataVer(v)
	}
	if v := os.Getenv("OKX_MIN_LEAD_DAYS"); v != "" {
		svc.MinLeadDays(v)
	}
	if v := os.Getenv("OKX_MIN_ASSETS"); v != "" {
		svc.MinAssets(v)
	}
	if v := os.Getenv("OKX_MAX_ASSETS"); v != "" {
		svc.MaxAssets(v)
	}
	if v := os.Getenv("OKX_MIN_AUM"); v != "" {
		svc.MinAum(v)
	}
	if v := os.Getenv("OKX_MAX_AUM"); v != "" {
		svc.MaxAum(v)
	}

	if v := os.Getenv("OKX_PAGE"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_PAGE: %v", err)
		}
		svc.Page(n)
	}
	if v := os.Getenv("OKX_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_LIMIT: %v", err)
		}
		svc.Limit(n)
	}

	res, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("copytrading_public_lead_traders: instType=%s sortType=%s dataVer=%s totalPage=%s ranks=%d", instType, sortType, res.DataVer, res.TotalPage, len(res.Ranks))
	if len(res.Ranks) == 0 {
		return
	}

	r := res.Ranks[0]
	log.Printf("first: uniqueCode=%s nickName=%s aum=%s pnl=%s pnlRatio=%s winRatio=%s", r.UniqueCode, r.NickName, r.Aum, r.Pnl, r.PnlRatio, r.WinRatio)
}
