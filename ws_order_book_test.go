package okx

import (
	"errors"
	"testing"
)

func TestWSOrderBookChecksum(t *testing.T) {
	bids := []OrderBookLevel{
		{Px: "3366.1", Sz: "7"},
		{Px: "3366", Sz: "6"},
	}
	asks := []OrderBookLevel{
		{Px: "3366.8", Sz: "9"},
		{Px: "3368", Sz: "8"},
	}

	got := wsOrderBookChecksum(bids, asks)
	// 校验字符串："3366.1:7:3366.8:9:3366:6:3368:8"
	const want int64 = -1881014294
	if got != want {
		t.Fatalf("checksum = %d, want %d", got, want)
	}
}

func TestWSOrderBookStore_Books_SnapshotAndUpdate(t *testing.T) {
	store := NewWSOrderBookStore(WSChannelBooks, "BTC-USDT")

	snapshot := &WSData[WSOrderBook]{
		Arg:    WSArg{Channel: WSChannelBooks, InstId: "BTC-USDT"},
		Action: "snapshot",
		Data: []WSOrderBook{{
			Asks: []OrderBookLevel{
				{Px: "101", Sz: "1", LiqOrd: "0", NumOrders: "1"},
				{Px: "101.5", Sz: "2", LiqOrd: "0", NumOrders: "1"},
				{Px: "102", Sz: "3", LiqOrd: "0", NumOrders: "1"},
			},
			Bids: []OrderBookLevel{
				{Px: "100", Sz: "1", LiqOrd: "0", NumOrders: "1"},
				{Px: "99.5", Sz: "2", LiqOrd: "0", NumOrders: "1"},
				{Px: "99", Sz: "3", LiqOrd: "0", NumOrders: "1"},
			},
			InstId:    "BTC-USDT",
			TS:        1,
			Checksum:  97597700,
			PrevSeqId: -1,
			SeqId:     10,
		}},
	}

	if err := store.Apply(snapshot); err != nil {
		t.Fatalf("Apply(snapshot) error = %v", err)
	}
	if !store.Ready() {
		t.Fatalf("expected ready")
	}
	got := store.Snapshot()
	if got.SeqId != 10 || got.TS != 1 || got.Checksum != 97597700 {
		t.Fatalf("snapshot meta = %#v, want seqId=10 ts=1 checksum=97597700", got)
	}
	if len(got.Bids) != 3 || got.Bids[0].Px != "100" || got.Bids[1].Px != "99.5" || got.Bids[2].Px != "99" {
		t.Fatalf("bids = %#v", got.Bids)
	}
	if len(got.Asks) != 3 || got.Asks[0].Px != "101" || got.Asks[1].Px != "101.5" || got.Asks[2].Px != "102" {
		t.Fatalf("asks = %#v", got.Asks)
	}

	update := &WSData[WSOrderBook]{
		Arg:    WSArg{Channel: WSChannelBooks, InstId: "BTC-USDT"},
		Action: "update",
		Data: []WSOrderBook{{
			Asks: []OrderBookLevel{
				{Px: "101.5", Sz: "7", LiqOrd: "0", NumOrders: "1"},
				{Px: "102", Sz: "0", LiqOrd: "0", NumOrders: "0"},
			},
			Bids: []OrderBookLevel{
				{Px: "99.8", Sz: "5", LiqOrd: "0", NumOrders: "1"},
				{Px: "99.5", Sz: "0", LiqOrd: "0", NumOrders: "0"},
			},
			InstId:    "BTC-USDT",
			TS:        2,
			Checksum:  -665260576,
			PrevSeqId: 10,
			SeqId:     15,
		}},
	}

	if err := store.Apply(update); err != nil {
		t.Fatalf("Apply(update) error = %v", err)
	}

	got = store.Snapshot()
	if got.SeqId != 15 || got.TS != 2 || got.Checksum != -665260576 {
		t.Fatalf("snapshot meta = %#v, want seqId=15 ts=2 checksum=-665260576", got)
	}
	if len(got.Bids) != 3 || got.Bids[0].Px != "100" || got.Bids[1].Px != "99.8" || got.Bids[2].Px != "99" {
		t.Fatalf("bids = %#v, want 100,99.8,99", got.Bids)
	}
	if len(got.Asks) != 2 || got.Asks[0].Px != "101" || got.Asks[1].Px != "101.5" || got.Asks[1].Sz != "7" {
		t.Fatalf("asks = %#v, want 101,101.5(sz=7)", got.Asks)
	}
}

func TestWSOrderBookStore_SequenceMismatch_Resets(t *testing.T) {
	store := NewWSOrderBookStore(WSChannelBooks, "BTC-USDT")

	snapshot := &WSData[WSOrderBook]{
		Arg:    WSArg{Channel: WSChannelBooks, InstId: "BTC-USDT"},
		Action: "snapshot",
		Data: []WSOrderBook{{
			Asks: []OrderBookLevel{{Px: "101", Sz: "1", LiqOrd: "0", NumOrders: "1"}},
			Bids: []OrderBookLevel{{Px: "100", Sz: "1", LiqOrd: "0", NumOrders: "1"}},
			TS:   1,
			// 校验字符串 "100:1:101:1"
			Checksum:  1189976625,
			PrevSeqId: -1,
			SeqId:     10,
		}},
	}
	if err := store.Apply(snapshot); err != nil {
		t.Fatalf("Apply(snapshot) error = %v", err)
	}

	update := &WSData[WSOrderBook]{
		Arg:    WSArg{Channel: WSChannelBooks, InstId: "BTC-USDT"},
		Action: "update",
		Data: []WSOrderBook{{
			Asks:      []OrderBookLevel{{Px: "101", Sz: "1", LiqOrd: "0", NumOrders: "1"}},
			Bids:      []OrderBookLevel{{Px: "100", Sz: "1", LiqOrd: "0", NumOrders: "1"}},
			TS:        2,
			Checksum:  1189976625,
			PrevSeqId: 9,
			SeqId:     11,
		}},
	}

	err := store.Apply(update)
	if err == nil {
		t.Fatalf("expected error")
	}
	var seqErr *WSOrderBookSequenceError
	if !errors.As(err, &seqErr) {
		t.Fatalf("error = %T, want *WSOrderBookSequenceError", err)
	}
	if store.Ready() {
		t.Fatalf("expected store reset after sequence mismatch")
	}
}

func TestWSOrderBookStore_ChecksumMismatch_Resets(t *testing.T) {
	store := NewWSOrderBookStore(WSChannelBooks, "BTC-USDT")

	snapshot := &WSData[WSOrderBook]{
		Arg:    WSArg{Channel: WSChannelBooks, InstId: "BTC-USDT"},
		Action: "snapshot",
		Data: []WSOrderBook{{
			Asks:      []OrderBookLevel{{Px: "101", Sz: "1", LiqOrd: "0", NumOrders: "1"}},
			Bids:      []OrderBookLevel{{Px: "100", Sz: "1", LiqOrd: "0", NumOrders: "1"}},
			TS:        1,
			Checksum:  1189976625,
			PrevSeqId: -1,
			SeqId:     10,
		}},
	}
	if err := store.Apply(snapshot); err != nil {
		t.Fatalf("Apply(snapshot) error = %v", err)
	}

	update := &WSData[WSOrderBook]{
		Arg:    WSArg{Channel: WSChannelBooks, InstId: "BTC-USDT"},
		Action: "update",
		Data: []WSOrderBook{{
			Asks:      []OrderBookLevel{{Px: "101", Sz: "1", LiqOrd: "0", NumOrders: "1"}},
			Bids:      []OrderBookLevel{{Px: "100", Sz: "1", LiqOrd: "0", NumOrders: "1"}},
			TS:        2,
			Checksum:  1,
			PrevSeqId: 10,
			SeqId:     11,
		}},
	}

	err := store.Apply(update)
	if err == nil {
		t.Fatalf("expected error")
	}
	var chkErr *WSOrderBookChecksumError
	if !errors.As(err, &chkErr) {
		t.Fatalf("error = %T, want *WSOrderBookChecksumError", err)
	}
	if store.Ready() {
		t.Fatalf("expected store reset after checksum mismatch")
	}
}

func TestWSOrderBookStore_Books5_Replace(t *testing.T) {
	store := NewWSOrderBookStore(WSChannelBooks5, "BTC-USDT")

	first := &WSData[WSOrderBook]{
		Arg:    WSArg{Channel: WSChannelBooks5, InstId: "BTC-USDT"},
		Action: "snapshot",
		Data: []WSOrderBook{{
			Asks: []OrderBookLevel{
				{Px: "11", Sz: "1", LiqOrd: "0", NumOrders: "1"},
				{Px: "12", Sz: "1", LiqOrd: "0", NumOrders: "1"},
			},
			Bids: []OrderBookLevel{
				{Px: "10", Sz: "1", LiqOrd: "0", NumOrders: "1"},
				{Px: "9", Sz: "1", LiqOrd: "0", NumOrders: "1"},
			},
			TS:       1,
			Checksum: -1041668651,
		}},
	}
	if err := store.Apply(first); err != nil {
		t.Fatalf("Apply(first) error = %v", err)
	}

	second := &WSData[WSOrderBook]{
		Arg:    WSArg{Channel: WSChannelBooks5, InstId: "BTC-USDT"},
		Action: "update",
		Data: []WSOrderBook{{
			Asks: []OrderBookLevel{
				{Px: "11", Sz: "3", LiqOrd: "0", NumOrders: "1"},
			},
			Bids: []OrderBookLevel{
				{Px: "10", Sz: "2", LiqOrd: "0", NumOrders: "1"},
			},
			TS:       2,
			Checksum: 66699475,
		}},
	}
	if err := store.Apply(second); err != nil {
		t.Fatalf("Apply(second) error = %v", err)
	}

	got := store.Snapshot()
	if len(got.Bids) != 1 || got.Bids[0].Px != "10" || got.Bids[0].Sz != "2" {
		t.Fatalf("bids = %#v, want only 10(sz=2)", got.Bids)
	}
	if len(got.Asks) != 1 || got.Asks[0].Px != "11" || got.Asks[0].Sz != "3" {
		t.Fatalf("asks = %#v, want only 11(sz=3)", got.Asks)
	}
}

func TestWSOrderBookStore_ApplyMessage_Filter(t *testing.T) {
	store := NewWSOrderBookStore(WSChannelBooks, "BTC-USDT")

	t.Run("non_order_book_message", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"tickers","instId":"BTC-USDT"},"data":[{"instId":"BTC-USDT","last":"1","bidPx":"1","bidSz":"1","askPx":"1","askSz":"1","ts":"1"}]}`)
		ok, err := store.ApplyMessage(msg)
		if err != nil {
			t.Fatalf("ApplyMessage() error = %v", err)
		}
		if ok {
			t.Fatalf("expected ok=false")
		}
	})

	t.Run("channel_mismatch", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"books5","instId":"BTC-USDT"},"action":"snapshot","data":[{"asks":[["11","1","0","1"]],"bids":[["10","1","0","1"]],"instId":"BTC-USDT","ts":"1","checksum":66699475}]}`)
		ok, err := store.ApplyMessage(msg)
		if err != nil {
			t.Fatalf("ApplyMessage() error = %v", err)
		}
		if ok {
			t.Fatalf("expected ok=false")
		}
	})

	t.Run("instId_mismatch", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"books","instId":"ETH-USDT"},"action":"snapshot","data":[{"asks":[["11","1","0","1"]],"bids":[["10","1","0","1"]],"instId":"ETH-USDT","ts":"1","checksum":66699475,"prevSeqId":-1,"seqId":1}]}`)
		ok, err := store.ApplyMessage(msg)
		if err != nil {
			t.Fatalf("ApplyMessage() error = %v", err)
		}
		if ok {
			t.Fatalf("expected ok=false")
		}
	})
}
