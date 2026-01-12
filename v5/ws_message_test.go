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

	t.Run("order_book_sprd_books5", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"sprd-books5","sprdId":"BTC-USDT_BTC-USDT-SWAP"},"action":"snapshot","data":[{"asks":[["1.9","1.1","3"]],"bids":[["1.8","0.165","1"]],"ts":"1724391380926","checksum":-1285595583,"prevSeqId":-1,"seqId":1724294007352168320}]}`)
		dm, ok, err := WSParseOrderBook(msg)
		if err != nil {
			t.Fatalf("WSParseOrderBook() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelSprdBooks5 || dm.Arg.SprdId != "BTC-USDT_BTC-USDT-SWAP" || dm.Action != "snapshot" {
			t.Fatalf("meta = %#v", dm)
		}
		if len(dm.Data) != 1 || dm.Data[0].TS != 1724391380926 || dm.Data[0].Checksum != -1285595583 || dm.Data[0].PrevSeqId != -1 || dm.Data[0].SeqId != 1724294007352168320 {
			t.Fatalf("data = %#v", dm.Data)
		}
		if len(dm.Data[0].Asks) != 1 || dm.Data[0].Asks[0].Px != "1.9" || dm.Data[0].Asks[0].Sz != "1.1" || dm.Data[0].Asks[0].NumOrders != "3" {
			t.Fatalf("asks = %#v", dm.Data[0].Asks)
		}
		if len(dm.Data[0].Bids) != 1 || dm.Data[0].Bids[0].Px != "1.8" || dm.Data[0].Bids[0].Sz != "0.165" || dm.Data[0].Bids[0].NumOrders != "1" {
			t.Fatalf("bids = %#v", dm.Data[0].Bids)
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

func TestWSParseDepositAndWithdrawalInfo(t *testing.T) {
	t.Run("deposit_info", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"deposit-info","uid":"289320000000"},"data":[{"actualDepBlkConfirm":"0","amt":"1","areaCodeFrom":"","ccy":"USDT","chain":"USDT-TRC20","depId":"88165462","from":"","fromWdId":"","pTime":"1674103661147","state":"0","subAcct":"test","to":"TEhFAqpuHa3LYxxxxx8ByNoGnrmexeGMw","ts":"1674103661123","txId":"bc5376","uid":"289320000000"}]}`)
		dm, ok, err := WSParseDepositInfo(msg)
		if err != nil {
			t.Fatalf("WSParseDepositInfo() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelDepositInfo || dm.Arg.UID != "289320000000" {
			t.Fatalf("arg = %#v", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].DepId != "88165462" || dm.Data[0].PTime != 1674103661147 || dm.Data[0].TS != 1674103661123 {
			t.Fatalf("data = %#v", dm.Data)
		}
	})

	t.Run("withdrawal_info", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"withdrawal-info","uid":"289333000000"},"data":[{"addrEx":null,"amt":"2","areaCodeFrom":"","areaCodeTo":"","ccy":"USDT","chain":"USDT-TRC20","clientId":"","fee":"0.8","feeCcy":"USDT","from":"","memo":"","nonTradableAsset":false,"note":"","pTime":"1674103268578","pmtId":"","state":"0","subAcct":"test","tag":"","to":"TN8CKTQMnpWfTxxxxxx8KipbJ24ErguhF","toAddrType":"1","ts":"1674103268472","txId":"","uid":"289333000000","wdId":"15447421"}]}`)
		dm, ok, err := WSParseWithdrawalInfo(msg)
		if err != nil {
			t.Fatalf("WSParseWithdrawalInfo() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelWithdrawalInfo || dm.Arg.UID != "289333000000" {
			t.Fatalf("arg = %#v", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].WdId != "15447421" || dm.Data[0].PTime != 1674103268578 || dm.Data[0].TS != 1674103268472 {
			t.Fatalf("data = %#v", dm.Data)
		}
		if dm.Data[0].Fee != "0.8" || dm.Data[0].FeeCcy != "USDT" {
			t.Fatalf("fee = %s %s", dm.Data[0].Fee, dm.Data[0].FeeCcy)
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

	t.Run("sprd_candle_channel_helper", func(t *testing.T) {
		if got, want := WSSprdCandleChannel("1m"), "sprd-candle1m"; got != want {
			t.Fatalf("WSSprdCandleChannel() = %q, want %q", got, want)
		}
		if got, want := WSSprdCandleChannel("sprd-candle1D"), "sprd-candle1D"; got != want {
			t.Fatalf("WSSprdCandleChannel() = %q, want %q", got, want)
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

	t.Run("sprd_candles", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"sprd-candle1D","sprdId":"BTC-USDT_BTC-USDT-SWAP"},"data":[["1597026383085","8533.02","8553.74","8527.17","8548.26","45247","0"]]}`)
		dm, ok, err := WSParseSprdCandles(msg)
		if err != nil {
			t.Fatalf("WSParseSprdCandles() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != "sprd-candle1D" || dm.Arg.SprdId != "BTC-USDT_BTC-USDT-SWAP" {
			t.Fatalf("arg = %#v, want sprd-candle1D sprdId=BTC-USDT_BTC-USDT-SWAP", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].TS != 1597026383085 || dm.Data[0].Open != "8533.02" || dm.Data[0].Confirm != "0" {
			t.Fatalf("data = %#v, want ts=1597026383085 open=8533.02 confirm=0", dm.Data)
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

	t.Run("sprd_tickers", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"sprd-tickers","sprdId":"BTC-USDT_BTC-USDT-SWAP"},"data":[{"sprdId":"BTC-USDT_BTC-USDT-SWAP","last":"4","lastSz":"0.01","askPx":"19.7","askSz":"5.79","bidPx":"5.9","bidSz":"5.79","open24h":"-7","high24h":"19.6","low24h":"-7","vol24h":"9.87","ts":"1715247061026"}]}`)
		dm, ok, err := WSParseSprdTickers(msg)
		if err != nil {
			t.Fatalf("WSParseSprdTickers() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelSprdTickers || dm.Arg.SprdId != "BTC-USDT_BTC-USDT-SWAP" {
			t.Fatalf("arg = %#v, want sprd-tickers sprdId=BTC-USDT_BTC-USDT-SWAP", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].SprdId != "BTC-USDT_BTC-USDT-SWAP" || dm.Data[0].Last != "4" || dm.Data[0].TS != 1715247061026 {
			t.Fatalf("data = %#v", dm.Data)
		}
	})

	t.Run("sprd_public_trades", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"sprd-public-trades","sprdId":"BTC-USDT_BTC-USDT-SWAP"},"data":[{"sprdId":"BTC-USDT_BTC-USDT-SWAP","side":"sell","sz":"0.1","px":"964.1","tradeId":"242720719","ts":"1654161641568"}]}`)
		dm, ok, err := WSParseSprdPublicTrades(msg)
		if err != nil {
			t.Fatalf("WSParseSprdPublicTrades() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelSprdPublicTrades || dm.Arg.SprdId != "BTC-USDT_BTC-USDT-SWAP" {
			t.Fatalf("arg = %#v, want sprd-public-trades sprdId=BTC-USDT_BTC-USDT-SWAP", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].SprdId != "BTC-USDT_BTC-USDT-SWAP" || dm.Data[0].TradeId != "242720719" || dm.Data[0].TS != 1654161641568 {
			t.Fatalf("data = %#v", dm.Data)
		}
		if dm.Data[0].Px != "964.1" || dm.Data[0].Sz != "0.1" || dm.Data[0].Side != "sell" {
			t.Fatalf("trade = %#v", dm.Data[0])
		}
	})
}

func TestWSParseBusinessPublicBlockTradingChannels(t *testing.T) {
	t.Run("public_struc_block_trades", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"public-struc-block-trades"},"data":[{"cTime":"1608267227834","blockTdId":"1802896","groupId":"","legs":[{"px":"0.323","sz":"25.0","instId":"BTC-USD-20220114-13250-C","side":"sell","tradeId":"15102"}]}]}`)
		dm, ok, err := WSParsePublicStrucBlockTrades(msg)
		if err != nil {
			t.Fatalf("WSParsePublicStrucBlockTrades() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelPublicStrucBlockTrades {
			t.Fatalf("arg = %#v, want channel public-struc-block-trades", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].BlockTdId != "1802896" || dm.Data[0].CTime != 1608267227834 {
			t.Fatalf("data = %#v", dm.Data)
		}
		if len(dm.Data[0].Legs) != 1 || dm.Data[0].Legs[0].InstId != "BTC-USD-20220114-13250-C" || dm.Data[0].Legs[0].TradeId != "15102" {
			t.Fatalf("legs = %#v", dm.Data[0].Legs)
		}
	})

	t.Run("public_block_trades", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"public-block-trades","instId":"BTC-USD-231020-5000-P"},"data":[{"fillVol":"5","fwdPx":"26808.16","groupId":"","idxPx":"27222.5","instId":"BTC-USD-231020-5000-P","markPx":"0.0022406326071111","px":"0.0048","side":"buy","sz":"1","tradeId":"633971452580106242","ts":"1697422572972"}]}`)
		dm, ok, err := WSParsePublicBlockTrades(msg)
		if err != nil {
			t.Fatalf("WSParsePublicBlockTrades() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelPublicBlockTrades || dm.Arg.InstId != "BTC-USD-231020-5000-P" {
			t.Fatalf("arg = %#v, want public-block-trades BTC-USD-231020-5000-P", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].InstId != "BTC-USD-231020-5000-P" || dm.Data[0].TradeId != "633971452580106242" || dm.Data[0].TS != 1697422572972 {
			t.Fatalf("data = %#v", dm.Data)
		}
		if dm.Data[0].Px != "0.0048" || dm.Data[0].Sz != "1" || dm.Data[0].Side != "buy" {
			t.Fatalf("trade = %#v", dm.Data[0])
		}
	})

	t.Run("block_tickers", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"block-tickers"},"data":[{"instType":"SWAP","instId":"LTC-USD-SWAP","volCcy24h":"0","vol24h":"0","ts":"1597026383085"}]}`)
		dm, ok, err := WSParseBlockTickers(msg)
		if err != nil {
			t.Fatalf("WSParseBlockTickers() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelBlockTickers {
			t.Fatalf("arg = %#v, want channel block-tickers", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].InstId != "LTC-USD-SWAP" || dm.Data[0].TS != 1597026383085 {
			t.Fatalf("data = %#v", dm.Data)
		}
	})
}

func TestWSParseBusinessSprdOrdersAndTrades(t *testing.T) {
	t.Run("sprd_orders", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"sprd-orders","sprdId":"BTC-USDT_BTC-USDT-SWAP","uid":"614488474791936"},"data":[{"sprdId":"BTC-USDT_BTC-UST-SWAP","ordId":"312269865356374016","clOrdId":"b1","tag":"","px":"999","sz":"3","ordType":"limit","side":"buy","fillSz":"0","fillPx":"","tradeId":"","accFillSz":"0","pendingFillSz":"2","pendingSettleSz":"1","canceledSz":"1","state":"live","avgPx":"0","cancelSource":"","uTime":"1597026383085","cTime":"1597026383085","code":"0","msg":"","reqId":"","amendResult":""}]}`)
		dm, ok, err := WSParseSprdOrders(msg)
		if err != nil {
			t.Fatalf("WSParseSprdOrders() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelSprdOrders || dm.Arg.SprdId != "BTC-USDT_BTC-USDT-SWAP" || dm.Arg.UID != "614488474791936" {
			t.Fatalf("arg = %#v", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].OrdId != "312269865356374016" || dm.Data[0].State != "live" || dm.Data[0].UTime != 1597026383085 {
			t.Fatalf("data = %#v", dm.Data)
		}
	})

	t.Run("sprd_trades", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"sprd-trades","sprdId":"BTC-USDT_BTC-USDT-SWAP","uid":"614488474791936"},"data":[{"sprdId":"BTC-USDT-SWAP_BTC-USDT-200329","tradeId":"123","ordId":"123445","clOrdId":"b16","tag":"","fillPx":"999","fillSz":"3","state":"filled","side":"buy","execType":"M","ts":"1597026383085","legs":[{"instId":"BTC-USDT-SWAP","px":"20000","sz":"3","szCont":"0.03","side":"buy","fillPnl":"","fee":"","feeCcy":"","tradeId":"1232342342"},{"instId":"BTC-USDT-200329","px":"21000","sz":"3","szCont":"0.03","side":"sell","fillPnl":"","fee":"","feeCcy":"","tradeId":"5345646634"}],"code":"","msg":""}]}`)
		dm, ok, err := WSParseSprdTrades(msg)
		if err != nil {
			t.Fatalf("WSParseSprdTrades() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelSprdTrades || dm.Arg.SprdId != "BTC-USDT_BTC-USDT-SWAP" || dm.Arg.UID != "614488474791936" {
			t.Fatalf("arg = %#v", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].TradeId != "123" || dm.Data[0].TS != 1597026383085 {
			t.Fatalf("data = %#v", dm.Data)
		}
		if len(dm.Data[0].Legs) != 2 || dm.Data[0].Legs[1].InstId != "BTC-USDT-200329" || dm.Data[0].Legs[1].Side != "sell" {
			t.Fatalf("legs = %#v", dm.Data[0].Legs)
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

func TestWSParsePublicDataChannels(t *testing.T) {
	t.Run("open_interest", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"open-interest","instId":"BTC-USDT-SWAP"},"data":[{"instId":"BTC-USDT-SWAP","instType":"SWAP","oi":"2216113.01000000309","oiCcy":"22161.1301000000309","oiUsd":"1939251795.54769270396321","ts":"1743041250440"}]}`)
		dm, ok, err := WSParseOpenInterest(msg)
		if err != nil {
			t.Fatalf("WSParseOpenInterest() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelOpenInterest || dm.Arg.InstId != "BTC-USDT-SWAP" {
			t.Fatalf("arg = %#v, want open-interest BTC-USDT-SWAP", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].InstType != "SWAP" || dm.Data[0].TS != 1743041250440 {
			t.Fatalf("data = %#v, want instType=SWAP ts=1743041250440", dm.Data)
		}
	})

	t.Run("funding_rate", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"funding-rate","instId":"BTC-USD-SWAP"},"data":[{"fundingRate":"0.0001875391284828","fundingTime":"1700726400000","instId":"BTC-USD-SWAP","instType":"SWAP","method":"current_period","maxFundingRate":"0.00375","minFundingRate":"-0.00375","nextFundingRate":"","nextFundingTime":"1700755200000","premium":"0.0001233824646391","settFundingRate":"0.0001699799259033","settState":"settled","ts":"1700724675402"}]}`)
		dm, ok, err := WSParseFundingRate(msg)
		if err != nil {
			t.Fatalf("WSParseFundingRate() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelFundingRate || dm.Arg.InstId != "BTC-USD-SWAP" {
			t.Fatalf("arg = %#v, want funding-rate BTC-USD-SWAP", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].FundingRate != "0.0001875391284828" || dm.Data[0].SettFundingRate != "0.0001699799259033" {
			t.Fatalf("data = %#v, want fundingRate/settFundingRate", dm.Data)
		}
	})

	t.Run("price_limit", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"price-limit","instId":"LTC-USD-190628"},"data":[{"instId":"LTC-USD-190628","buyLmt":"200","sellLmt":"300","ts":"1597026383085","enabled":true}]}`)
		dm, ok, err := WSParsePriceLimit(msg)
		if err != nil {
			t.Fatalf("WSParsePriceLimit() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelPriceLimit || dm.Arg.InstId != "LTC-USD-190628" {
			t.Fatalf("arg = %#v, want price-limit LTC-USD-190628", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].BuyLmt != "200" || !dm.Data[0].Enabled {
			t.Fatalf("data = %#v, want buyLmt=200 enabled=true", dm.Data)
		}
	})

	t.Run("mark_price", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"mark-price","instId":"BTC-USDT"},"data":[{"instType":"MARGIN","instId":"BTC-USDT","markPx":"42310.6","ts":"1630049139746"}]}`)
		dm, ok, err := WSParseMarkPrice(msg)
		if err != nil {
			t.Fatalf("WSParseMarkPrice() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelMarkPrice || dm.Arg.InstId != "BTC-USDT" {
			t.Fatalf("arg = %#v, want mark-price BTC-USDT", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].MarkPx != "42310.6" || dm.Data[0].TS != 1630049139746 {
			t.Fatalf("data = %#v, want markPx=42310.6 ts=1630049139746", dm.Data)
		}
	})

	t.Run("index_tickers", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"index-tickers","instId":"BTC-USDT"},"data":[{"instId":"BTC-USDT","idxPx":"0.1","high24h":"0.5","low24h":"0.1","open24h":"0.1","sodUtc0":"0.1","sodUtc8":"0.1","ts":"1597026383085"}]}`)
		dm, ok, err := WSParseIndexTickers(msg)
		if err != nil {
			t.Fatalf("WSParseIndexTickers() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelIndexTickers || dm.Arg.InstId != "BTC-USDT" {
			t.Fatalf("arg = %#v, want index-tickers BTC-USDT", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].IdxPx != "0.1" || dm.Data[0].TS != 1597026383085 {
			t.Fatalf("data = %#v, want idxPx=0.1 ts=1597026383085", dm.Data)
		}
	})

	t.Run("opt_summary", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"opt-summary","instFamily":"BTC-USD"},"data":[{"instType":"OPTION","instId":"BTC-USD-230224-18000-C","uly":"BTC-USD","askVol":"1","bidVol":"2","markVol":"3","realVol":"4","delta":"0.1","gamma":"0.2","theta":"0.3","vega":"0.4","volLv":"0.5","fwdPx":"17000","distance":"0.9","ts":"1672286551080"}]}`)
		dm, ok, err := WSParseOptSummary(msg)
		if err != nil {
			t.Fatalf("WSParseOptSummary() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelOptSummary || dm.Arg.InstFamily != "BTC-USD" {
			t.Fatalf("arg = %#v, want opt-summary BTC-USD", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].InstId != "BTC-USD-230224-18000-C" || dm.Data[0].TS != 1672286551080 {
			t.Fatalf("data = %#v, want instId=... ts=1672286551080", dm.Data)
		}
	})

	t.Run("liquidation_orders", func(t *testing.T) {
		msg := []byte(`{"id":"1512","arg":{"channel":"liquidation-orders","instType":"SWAP"},"data":[{"details":[{"bkLoss":"0","bkPx":"0.007831","ccy":"","posSide":"short","side":"buy","sz":"13","ts":"1692266434010"}],"instFamily":"IOST-USDT","instId":"IOST-USDT-SWAP","instType":"SWAP","uly":"IOST-USDT"}]}`)
		dm, ok, err := WSParseLiquidationOrders(msg)
		if err != nil {
			t.Fatalf("WSParseLiquidationOrders() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != WSChannelLiquidationOrders || dm.Arg.InstType != "SWAP" {
			t.Fatalf("arg = %#v, want liquidation-orders SWAP", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].InstId != "IOST-USDT-SWAP" || len(dm.Data[0].Details) != 1 || dm.Data[0].Details[0].TS != 1692266434010 {
			t.Fatalf("data = %#v, want instId=IOST-USDT-SWAP details ts=1692266434010", dm.Data)
		}
	})
}

func TestWSParseMarkPriceAndIndexCandles(t *testing.T) {
	t.Run("channel_helpers", func(t *testing.T) {
		if got, want := WSMarkPriceCandleChannel("1D"), "mark-price-candle1D"; got != want {
			t.Fatalf("WSMarkPriceCandleChannel() = %q, want %q", got, want)
		}
		if got, want := WSMarkPriceCandleChannel("mark-price-candle1m"), "mark-price-candle1m"; got != want {
			t.Fatalf("WSMarkPriceCandleChannel() = %q, want %q", got, want)
		}
		if got, want := WSIndexCandleChannel("30m"), "index-candle30m"; got != want {
			t.Fatalf("WSIndexCandleChannel() = %q, want %q", got, want)
		}
		if got, want := WSIndexCandleChannel("index-candle1Dutc"), "index-candle1Dutc"; got != want {
			t.Fatalf("WSIndexCandleChannel() = %q, want %q", got, want)
		}
	})

	t.Run("mark_price_candles", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"mark-price-candle1D","instId":"BTC-USD-190628"},"data":[["1597026383085","3.721","3.743","3.677","3.708","0"]]}`)
		dm, ok, err := WSParseMarkPriceCandles(msg)
		if err != nil {
			t.Fatalf("WSParseMarkPriceCandles() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != "mark-price-candle1D" || dm.Arg.InstId != "BTC-USD-190628" {
			t.Fatalf("arg = %#v, want mark-price-candle1D BTC-USD-190628", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].TS != 1597026383085 || dm.Data[0].Open != "3.721" || dm.Data[0].Confirm != "0" {
			t.Fatalf("data = %#v, want ts=1597026383085 open=3.721 confirm=0", dm.Data)
		}
	})

	t.Run("index_candles", func(t *testing.T) {
		msg := []byte(`{"arg":{"channel":"index-candle30m","instId":"BTC-USD"},"data":[["1597026383085","3.721","3.743","3.677","3.708","1"]]}`)
		dm, ok, err := WSParseIndexCandles(msg)
		if err != nil {
			t.Fatalf("WSParseIndexCandles() error = %v", err)
		}
		if !ok || dm == nil {
			t.Fatalf("expected ok")
		}
		if dm.Arg.Channel != "index-candle30m" || dm.Arg.InstId != "BTC-USD" {
			t.Fatalf("arg = %#v, want index-candle30m BTC-USD", dm.Arg)
		}
		if len(dm.Data) != 1 || dm.Data[0].TS != 1597026383085 || dm.Data[0].Close != "3.708" || dm.Data[0].Confirm != "1" {
			t.Fatalf("data = %#v, want ts=1597026383085 close=3.708 confirm=1", dm.Data)
		}
	})
}
