package okx

import "testing"

func TestWSParseEvent(t *testing.T) {
	t.Run("subscribe_event", func(t *testing.T) {
		msg := []byte(`{"id":"1512","event":"subscribe","arg":{"channel":"orders","instType":"SPOT","instId":"BTC-USDT"},"connId":"a4d3ae55"}`)
		ev, ok, err := WSParseEvent(msg)
		if err != nil {
			t.Fatalf("WSParseEvent() error = %v", err)
		}
		if !ok {
			t.Fatalf("expected ok")
		}
		if ev.Event != "subscribe" {
			t.Fatalf("Event = %q, want %q", ev.Event, "subscribe")
		}
		if ev.Arg == nil || ev.Arg.Channel != "orders" {
			t.Fatalf("Arg = %#v, want channel orders", ev.Arg)
		}
		if ev.ConnID != "a4d3ae55" {
			t.Fatalf("ConnID = %q, want %q", ev.ConnID, "a4d3ae55")
		}
	})

	t.Run("not_event_is_data", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"orders","instType":"SPOT","instId":"BTC-USDT"},"data":[{"ordId":"1"}]}`)
		_, ok, err := WSParseEvent(msg)
		if err != nil {
			t.Fatalf("WSParseEvent() error = %v", err)
		}
		if ok {
			t.Fatalf("expected ok=false")
		}
	})
}

func TestWSParseData(t *testing.T) {
	t.Run("orders_data", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"orders","instType":"SPOT","instId":"BTC-USDT","uid":"614488474791936"},"data":[{"instType":"SPOT","instId":"BTC-USDT","ordId":"1","side":"buy","ordType":"limit","state":"live","px":"1","sz":"1","avgPx":"0","fillPx":"0","fillSz":"0","accFillSz":"0","uTime":"1597026383085","cTime":"1597026383085"}]}`)
		dm, ok, err := WSParseData[TradeOrder](msg)
		if err != nil {
			t.Fatalf("WSParseData() error = %v", err)
		}
		if !ok {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != "orders" {
			t.Fatalf("Arg.Channel = %q, want %q", dm.Arg.Channel, "orders")
		}
		if dm.Arg.UID != "614488474791936" {
			t.Fatalf("Arg.UID = %q, want %q", dm.Arg.UID, "614488474791936")
		}
		if len(dm.Data) != 1 {
			t.Fatalf("len(Data) = %d, want %d", len(dm.Data), 1)
		}
		if dm.Data[0].OrdId != "1" {
			t.Fatalf("Data[0].OrdId = %q, want %q", dm.Data[0].OrdId, "1")
		}
	})

	t.Run("fills_data", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"fills","instId":"BTC-USDT-SWAP","uid":"614488474791111"},"data":[{"instId":"BTC-USDT-SWAP","fillSz":"100","fillPx":"70000","side":"buy","ts":"1705449605015","ordId":"680800019749904384","clOrdId":"1234567890","tradeId":"12345","execType":"T","count":"10"}]}`)
		dm, ok, err := WSParseData[WSFill](msg)
		if err != nil {
			t.Fatalf("WSParseData() error = %v", err)
		}
		if !ok {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != "fills" {
			t.Fatalf("Arg.Channel = %q, want %q", dm.Arg.Channel, "fills")
		}
		if len(dm.Data) != 1 || dm.Data[0].TradeId != "12345" {
			t.Fatalf("Data = %#v, want tradeId 12345", dm.Data)
		}
	})

	t.Run("not_data_is_event", func(t *testing.T) {
		msg := []byte(`{"event":"subscribe","arg":{"channel":"orders","instType":"SPOT"},"connId":"a4d3ae55"}`)
		_, ok, err := WSParseData[TradeOrder](msg)
		if err != nil {
			t.Fatalf("WSParseData() error = %v", err)
		}
		if ok {
			t.Fatalf("expected ok=false")
		}
	})
}

func TestWSParseChannelHelpers(t *testing.T) {
	t.Run("orders", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"orders","instType":"SPOT","instId":"BTC-USDT","uid":"614488474791936"},"data":[{"instType":"SPOT","instId":"BTC-USDT","ordId":"1","side":"buy","ordType":"limit","state":"live","px":"1","sz":"1","avgPx":"0","fillPx":"0","fillSz":"0","accFillSz":"0","uTime":"1597026383085","cTime":"1597026383085"}]}`)
		dm, ok, err := WSParseOrders(msg)
		if err != nil {
			t.Fatalf("WSParseOrders() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != "orders" || dm.Data[0].OrdId != "1" {
			t.Fatalf("parsed = %#v", dm)
		}
	})

	t.Run("fills", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"fills","instId":"BTC-USDT-SWAP","uid":"614488474791111"},"data":[{"instId":"BTC-USDT-SWAP","fillSz":"100","fillPx":"70000","side":"buy","ts":"1705449605015","ordId":"680800019749904384","clOrdId":"1234567890","tradeId":"12345","execType":"T","count":"10"}]}`)
		dm, ok, err := WSParseFills(msg)
		if err != nil {
			t.Fatalf("WSParseFills() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != "fills" || dm.Data[0].TradeId != "12345" {
			t.Fatalf("parsed = %#v", dm)
		}
	})

	t.Run("channel_mismatch", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"fills"},"data":[{"instId":"BTC-USDT-SWAP"}]}`)
		_, ok, err := WSParseOrders(msg)
		if err != nil {
			t.Fatalf("WSParseOrders() error = %v", err)
		}
		if ok {
			t.Fatalf("expected ok=false")
		}
	})
}

func TestWSParseAccountPositionsAndBalanceAndPosition(t *testing.T) {
	t.Run("account_data_with_paging", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"account","uid":"44*********584"},"eventType":"snapshot","curPage":1,"lastPage":true,"data":[{"uTime":"1597026383085","totalEq":"1","adjEq":"1","availEq":"1","details":[{"ccy":"BTC","eq":"1","cashBal":"1","availEq":"1"}]}]}`)
		dm, ok, err := WSParseAccount(msg)
		if err != nil {
			t.Fatalf("WSParseAccount() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.EventType != "snapshot" || dm.CurPage != 1 || dm.LastPage != true {
			t.Fatalf("meta = %#v, want snapshot page=1 last=true", dm)
		}
		if len(dm.Data) != 1 || dm.Data[0].TotalEq != "1" || dm.Data[0].AvailEq != "1" {
			t.Fatalf("data = %#v, want totalEq=1 availEq=1", dm.Data)
		}
		if len(dm.Data[0].Details) != 1 || dm.Data[0].Details[0].Ccy != "BTC" || dm.Data[0].Details[0].AvailEq != "1" {
			t.Fatalf("details = %#v, want BTC availEq=1", dm.Data[0].Details)
		}
	})

	t.Run("positions_data", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"positions","instType":"SWAP","uid":"44*********584"},"eventType":"snapshot","curPage":1,"lastPage":true,"data":[{"instType":"SWAP","instId":"BTC-USDT-SWAP","posId":"1","tradeId":"2","posSide":"long","pos":"10","mgnMode":"cross","ccy":"BTC","uTime":"1597026383085"}]}`)
		dm, ok, err := WSParsePositions(msg)
		if err != nil {
			t.Fatalf("WSParsePositions() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelPositions {
			t.Fatalf("Arg.Channel = %q, want %q", dm.Arg.Channel, WSChannelPositions)
		}
		if len(dm.Data) != 1 || dm.Data[0].InstId != "BTC-USDT-SWAP" || dm.Data[0].TradeId != "2" || dm.Data[0].UTime != 1597026383085 {
			t.Fatalf("data = %#v, want instId BTC-USDT-SWAP tradeId=2 uTime=1597026383085", dm.Data)
		}
	})

	t.Run("balance_and_position_data", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"balance_and_position","uid":"77982378738415879"},"data":[{"pTime":"1597026383085","eventType":"snapshot","balData":[{"ccy":"BTC","cashBal":"1","uTime":"1597026383085"}],"posData":[{"posId":"1111111111","tradeId":"2","instId":"BTC-USD-191018","instType":"FUTURES","mgnMode":"cross","posSide":"long","pos":"10","ccy":"BTC","posCcy":"","avgPx":"3320","uTime":"1597026383085"}]}]}`)
		dm, ok, err := WSParseBalanceAndPosition(msg)
		if err != nil {
			t.Fatalf("WSParseBalanceAndPosition() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if len(dm.Data) != 1 || dm.Data[0].PTime != 1597026383085 || dm.Data[0].EventType != "snapshot" {
			t.Fatalf("data = %#v, want pTime=1597026383085 eventType=snapshot", dm.Data)
		}
		if len(dm.Data[0].BalData) != 1 || dm.Data[0].BalData[0].Ccy != "BTC" || dm.Data[0].BalData[0].CashBal != "1" {
			t.Fatalf("balData = %#v, want BTC cashBal=1", dm.Data[0].BalData)
		}
		if len(dm.Data[0].PosData) != 1 || dm.Data[0].PosData[0].InstType != "FUTURES" || dm.Data[0].PosData[0].PosId != "1111111111" {
			t.Fatalf("posData = %#v, want FUTURES posId=1111111111", dm.Data[0].PosData)
		}
	})
}

func TestWSParsePublicTickersTradesAndOrderBook(t *testing.T) {
	t.Run("tickers", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"tickers","instId":"BTC-USDT"},"data":[{"instType":"SPOT","instId":"BTC-USDT","last":"9999.99","bidPx":"8888.88","bidSz":"5","askPx":"9999.99","askSz":"11","ts":"1597026383085"}]}`)
		dm, ok, err := WSParseTickers(msg)
		if err != nil {
			t.Fatalf("WSParseTickers() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelTickers || dm.Arg.InstId != "BTC-USDT" {
			t.Fatalf("arg = %#v, want tickers BTC-USDT", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].InstId != "BTC-USDT" || dm.Data[0].Last != "9999.99" || dm.Data[0].TS != 1597026383085 {
			t.Fatalf("data = %#v, want instId BTC-USDT last=9999.99 ts=1597026383085", dm.Data)
		}
	})

	t.Run("trades", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"trades","instId":"BTC-USDT"},"data":[{"instId":"BTC-USDT","tradeId":"1","px":"100","sz":"1","side":"buy","ts":"1597026383085","count":"3","source":"0","seqId":1234}]}`)
		dm, ok, err := WSParseTrades(msg)
		if err != nil {
			t.Fatalf("WSParseTrades() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelTrades || dm.Arg.InstId != "BTC-USDT" {
			t.Fatalf("arg = %#v, want trades BTC-USDT", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].TradeId != "1" || dm.Data[0].Count != "3" || dm.Data[0].Source != "0" || dm.Data[0].SeqId != 1234 {
			t.Fatalf("data = %#v, want tradeId=1 count=3 source=0 seqId=1234", dm.Data)
		}
	})

	t.Run("order_book_books5", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"books5","instId":"BTC-USDT"},"action":"snapshot","data":[{"asks":[["8476.98","415","0","13"]],"bids":[["8476.97","10","0","1"]],"instId":"BTC-USDT","ts":"1597026383085","checksum":-855196043,"prevSeqId":-1,"seqId":10}]}`)
		dm, ok, err := WSParseOrderBook(msg)
		if err != nil {
			t.Fatalf("WSParseOrderBook() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelBooks5 || dm.Action != "snapshot" {
			t.Fatalf("meta = %#v, want channel=books5 action=snapshot", dm)
		}
		if len(dm.Data) != 1 || dm.Data[0].Checksum != -855196043 || dm.Data[0].PrevSeqId != -1 || dm.Data[0].SeqId != 10 || dm.Data[0].TS != 1597026383085 {
			t.Fatalf("data = %#v, want checksum=-855196043 prevSeqId=-1 seqId=10 ts=1597026383085", dm.Data)
		}
		if len(dm.Data[0].Asks) != 1 || dm.Data[0].Asks[0].Px != "8476.98" || dm.Data[0].Asks[0].Sz != "415" {
			t.Fatalf("asks = %#v, want px=8476.98 sz=415", dm.Data[0].Asks)
		}
		if len(dm.Data[0].Bids) != 1 || dm.Data[0].Bids[0].Px != "8476.97" || dm.Data[0].Bids[0].Sz != "10" {
			t.Fatalf("bids = %#v, want px=8476.97 sz=10", dm.Data[0].Bids)
		}
	})

	t.Run("order_book_channel_mismatch", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"tickers","instId":"BTC-USDT"},"action":"snapshot","data":[{"asks":[["1","2","0","1"]],"bids":[["1","2","0","1"]],"instId":"BTC-USDT","ts":"1597026383085","checksum":1,"prevSeqId":-1,"seqId":10}]}`)
		_, ok, err := WSParseOrderBook(msg)
		if err != nil {
			t.Fatalf("WSParseOrderBook() error = %v", err)
		}
		if ok {
			t.Fatalf("expected ok=false")
		}
	})
}

func TestWSParseBusinessCandlesAndTradesAll(t *testing.T) {
	t.Run("candle_channel_helper", func(t *testing.T) {
		if got, want := WSCandleChannel("1m"), "candle1m"; got != want {
			t.Fatalf("WSCandleChannel() = %q, want %q", got, want)
		}
		if got, want := WSCandleChannel("candle1D"), "candle1D"; got != want {
			t.Fatalf("WSCandleChannel() = %q, want %q", got, want)
		}
	})

	t.Run("candles", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"candle1D","instId":"BTC-USDT"},"data":[["1629993600000","42500","48199.9","41006.1","41006.1","3587.41204591","166741046.22583129","166741046.22583129","0"]]}`)
		dm, ok, err := WSParseCandles(msg)
		if err != nil {
			t.Fatalf("WSParseCandles() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != "candle1D" || dm.Arg.InstId != "BTC-USDT" {
			t.Fatalf("arg = %#v, want candle1D BTC-USDT", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].TS != 1629993600000 || dm.Data[0].Open != "42500" || dm.Data[0].Confirm != "0" {
			t.Fatalf("data = %#v, want ts=1629993600000 open=42500 confirm=0", dm.Data)
		}
	})

	t.Run("trades_all", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"trades-all","instId":"BTC-USDT"},"data":[{"instId":"BTC-USDT","tradeId":"1","px":"100","sz":"1","side":"buy","ts":"1597026383085","source":"0"}]}`)
		dm, ok, err := WSParseTradesAll(msg)
		if err != nil {
			t.Fatalf("WSParseTradesAll() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelTradesAll || dm.Arg.InstId != "BTC-USDT" {
			t.Fatalf("arg = %#v, want trades-all BTC-USDT", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].TradeId != "1" || dm.Data[0].Source != "0" {
			t.Fatalf("data = %#v, want tradeId=1 source=0", dm.Data)
		}
	})
}

func TestWSParsePublicOptionTradesAndCallAuctionDetails(t *testing.T) {
	t.Run("option_trades", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"option-trades","instType":"OPTION","instFamily":"BTC-USD"},"data":[{"instFamily":"BTC-USD","instId":"BTC-USD-230224-18000-C","markPx":"0.04690107010619562","optType":"C","px":"0.045","side":"sell","sz":"2","tradeId":"38","fillVol":"0.1","fwdPx":"17000","idxPx":"16537.2","ts":"1672286551080"}]}`)
		dm, ok, err := WSParseOptionTrades(msg)
		if err != nil {
			t.Fatalf("WSParseOptionTrades() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelOptionTrades {
			t.Fatalf("arg = %#v, want channel option-trades", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].InstFamily != "BTC-USD" || dm.Data[0].OptType != "C" || dm.Data[0].TS != 1672286551080 {
			t.Fatalf("data = %#v, want instFamily=BTC-USD optType=C ts=1672286551080", dm.Data)
		}
	})

	t.Run("call_auction_details", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"call-auction-details","instId":"ONDO-USDC"},"data":[{"instId":"ONDO-USDC","eqPx":"0.6","matchedSz":"44978","unmatchedSz":"123","state":"continuous_trading","auctionEndTime":"1726542000000","ts":"1726542000007"}]}`)
		dm, ok, err := WSParseCallAuctionDetails(msg)
		if err != nil {
			t.Fatalf("WSParseCallAuctionDetails() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelCallAuctionDetails || dm.Arg.InstId != "ONDO-USDC" {
			t.Fatalf("arg = %#v, want call-auction-details ONDO-USDC", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].EqPx != "0.6" || dm.Data[0].MatchedSz != "44978" || dm.Data[0].AuctionEndTime != 1726542000000 {
			t.Fatalf("data = %#v, want eqPx=0.6 matchedSz=44978 auctionEndTime=1726542000000", dm.Data)
		}
	})
}
