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
