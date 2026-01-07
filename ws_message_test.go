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
