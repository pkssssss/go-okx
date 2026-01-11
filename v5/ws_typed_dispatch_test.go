package okx

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestWSClient_DispatchTyped_DropsWhenQueueFull(t *testing.T) {
	errCh := make(chan error, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask), // unbuffered + no receiver => full
		errHandler: func(err error) {
			select {
			case errCh <- err:
			default:
			}
		},
	}

	w.dispatchTyped(wsTypedTask{kind: wsTypedKindOrders, orders: []TradeOrder{{OrdId: "o1"}}})

	select {
	case err := <-errCh:
		if err == nil || !strings.Contains(err.Error(), "queue full") || !strings.Contains(err.Error(), "kind=orders") {
			t.Fatalf("err = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting error")
	}
}

func TestWSClient_TypedDispatchLoop_PanicRecovered(t *testing.T) {
	errCh := make(chan error, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
		errHandler: func(err error) {
			select {
			case errCh <- err:
			default:
			}
		},
	}

	w.OnOrders(func(order TradeOrder) {
		_ = order
		panic("boom")
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.dispatchTyped(wsTypedTask{kind: wsTypedKindOrders, orders: []TradeOrder{{OrdId: "o1"}}})

	select {
	case err := <-errCh:
		if err == nil || !strings.Contains(err.Error(), "panic") || !strings.Contains(err.Error(), "kind=orders") {
			t.Fatalf("err = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting error")
	}
}

func TestWSClient_onDataMessage_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan TradeOrder, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnOrders(func(order TradeOrder) {
		select {
		case gotCh <- order:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"orders"},"data":[{"instType":"SWAP","instId":"BTC-USDT-SWAP","ordId":"o1","clOrdId":"c1"}]}`))

	select {
	case o := <-gotCh:
		if o.OrdId != "o1" || o.ClOrdId != "c1" {
			t.Fatalf("order = %#v", o)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting order")
	}
}

func TestWSClient_onDataMessage_DepositInfo_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan WSDepositInfo, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnDepositInfo(func(info WSDepositInfo) {
		select {
		case gotCh <- info:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"deposit-info"},"data":[{"ccy":"USDT","chain":"USDT-TRC20","amt":"1","depId":"d1","ts":"1674103661123","pTime":"1674103661147","subAcct":"test","uid":"u1"}]}`))

	select {
	case info := <-gotCh:
		if info.DepId != "d1" || info.Ccy != "USDT" || info.TS != 1674103661123 || info.PTime != 1674103661147 || info.UID != "u1" {
			t.Fatalf("info = %#v", info)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting deposit info")
	}
}

func TestWSClient_onDataMessage_SprdOrders_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan SprdOrder, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnSprdOrders(func(order SprdOrder) {
		select {
		case gotCh <- order:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"sprd-orders"},"data":[{"sprdId":"BTC-USDT_BTC-USDT-SWAP","ordId":"o1","clOrdId":"c1","px":"1","sz":"1","side":"buy","ordType":"limit","state":"live","uTime":"1597026383085","cTime":"1597026383085"}]}`))

	select {
	case o := <-gotCh:
		if o.OrdId != "o1" || o.ClOrdId != "c1" || o.SprdId != "BTC-USDT_BTC-USDT-SWAP" || o.UTime != 1597026383085 {
			t.Fatalf("order = %#v", o)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting order")
	}
}

func TestWSClient_onDataMessage_SprdTrades_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan WSSprdTrade, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnSprdTrades(func(trade WSSprdTrade) {
		select {
		case gotCh <- trade:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"sprd-trades"},"data":[{"sprdId":"s1","tradeId":"t1","ordId":"o1","clOrdId":"c1","fillPx":"999","fillSz":"3","state":"filled","side":"buy","execType":"M","ts":"1597026383085","legs":[{"instId":"BTC-USDT-SWAP","px":"20000","sz":"3","szCont":"0.03","side":"buy","fillPnl":"","fee":"","feeCcy":"","tradeId":"lt1"}],"code":"","msg":""}]}`))

	select {
	case tr := <-gotCh:
		if tr.TradeId != "t1" || tr.OrdId != "o1" || tr.TS != 1597026383085 {
			t.Fatalf("trade = %#v", tr)
		}
		if len(tr.Legs) != 1 || tr.Legs[0].InstId != "BTC-USDT-SWAP" {
			t.Fatalf("legs = %#v", tr.Legs)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting trade")
	}
}

func TestWSClient_onDataMessage_Tickers_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan MarketTicker, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnTickers(func(ticker MarketTicker) {
		select {
		case gotCh <- ticker:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"tickers"},"data":[{"instId":"BTC-USDT","last":"1","bidPx":"0.9","bidSz":"1","askPx":"1.1","askSz":"2","ts":"1700000000000"}]}`))

	select {
	case tk := <-gotCh:
		if tk.InstId != "BTC-USDT" || tk.Last != "1" || tk.TS != 1700000000000 {
			t.Fatalf("ticker = %#v", tk)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting ticker")
	}
}

func TestWSClient_onDataMessage_Candles_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan WSCandle, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnCandles(func(candle WSCandle) {
		select {
		case gotCh <- candle:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"candle1m","instId":"BTC-USDT"},"data":[["1700000000000","1","2","0.5","1.5","100","10","15","1"]]}`))

	select {
	case c := <-gotCh:
		if c.Arg.InstId != "BTC-USDT" || c.Arg.Channel != "candle1m" {
			t.Fatalf("candle arg = %#v", c.Arg)
		}
		if c.Candle.TS != 1700000000000 || c.Candle.Close != "1.5" || c.Candle.Confirm != "1" {
			t.Fatalf("candle = %#v", c.Candle)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting candle")
	}
}

func TestWSClient_onDataMessage_PriceCandles_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan WSPriceCandle, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnPriceCandles(func(candle WSPriceCandle) {
		select {
		case gotCh <- candle:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"mark-price-candle1m","instId":"BTC-USDT-SWAP"},"data":[["1700000000000","1","2","0.5","1.5","1"]]}`))

	select {
	case c := <-gotCh:
		if c.Arg.InstId != "BTC-USDT-SWAP" || c.Arg.Channel != "mark-price-candle1m" {
			t.Fatalf("price candle arg = %#v", c.Arg)
		}
		if c.Candle.TS != 1700000000000 || c.Candle.Close != "1.5" || c.Candle.Confirm != "1" {
			t.Fatalf("price candle = %#v", c.Candle)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting price candle")
	}
}

func TestWSClient_onDataMessage_SprdTickers_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan MarketSprdTicker, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnSprdTickers(func(ticker MarketSprdTicker) {
		select {
		case gotCh <- ticker:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"sprd-tickers"},"data":[{"sprdId":"s1","last":"1","bidPx":"0.9","bidSz":"1","askPx":"1.1","askSz":"2","ts":"1700000000000"}]}`))

	select {
	case tk := <-gotCh:
		if tk.SprdId != "s1" || tk.Last != "1" || tk.TS != 1700000000000 {
			t.Fatalf("ticker = %#v", tk)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting sprd ticker")
	}
}

func TestWSClient_onDataMessage_OrderBook_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan WSData[WSOrderBook], 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnOrderBook(func(data WSData[WSOrderBook]) {
		select {
		case gotCh <- data:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"books5","instId":"BTC-USDT"},"action":"snapshot","data":[{"asks":[["1","2","0","1"]],"bids":[["1","2","0","1"]],"instId":"BTC-USDT","ts":"1700000000000","checksum":1,"prevSeqId":-1,"seqId":10}]}`))

	select {
	case dm := <-gotCh:
		if dm.Arg.Channel != "books5" || dm.Arg.InstId != "BTC-USDT" || dm.Action != "snapshot" {
			t.Fatalf("dm = %#v", dm)
		}
		if len(dm.Data) != 1 || dm.Data[0].InstId != "BTC-USDT" || dm.Data[0].TS != 1700000000000 {
			t.Fatalf("data = %#v", dm.Data)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting order book")
	}
}

func TestWSClient_onDataMessage_LiquidationWarning_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan WSLiquidationWarning, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnLiquidationWarning(func(warning WSLiquidationWarning) {
		select {
		case gotCh <- warning:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"liquidation-warning","uid":"u1","instType":"FUTURES"},"data":[{"instType":"FUTURES","mgnMode":"isolated","posId":"p1","posSide":"long","pos":"1","posCcy":"","instId":"ETH-USD-210430","lever":"10","markPx":"2353.849","mgnRatio":"11.731726509588816","ccy":"ETH","cTime":"1619507758793","uTime":"1619507761462","pTime":"1619507761462"}]}`))

	select {
	case warning := <-gotCh:
		if warning.PosId != "p1" || warning.InstId != "ETH-USD-210430" || warning.MarkPx != "2353.849" || warning.CTime != 1619507758793 {
			t.Fatalf("warning = %#v", warning)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting liquidation warning")
	}
}

func TestWSClient_onDataMessage_AccountGreeks_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan AccountGreeks, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnAccountGreeks(func(greeks AccountGreeks) {
		select {
		case gotCh <- greeks:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"account-greeks","ccy":"BTC","uid":"u1"},"data":[{"ccy":"BTC","deltaBS":"1.1","deltaPA":"2.2","gammaBS":"0","gammaPA":"0.1","thetaBS":"0","thetaPA":"0","vegaBS":"0","vegaPA":"0","ts":"1597026383085"}]}`))

	select {
	case greeks := <-gotCh:
		if greeks.Ccy != "BTC" || greeks.DeltaBS != "1.1" || greeks.TS != 1597026383085 {
			t.Fatalf("greeks = %#v", greeks)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting account greeks")
	}
}

func TestWSClient_onDataMessage_Status_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan SystemStatus, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnStatus(func(status SystemStatus) {
		select {
		case gotCh <- status:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"status"},"data":[{"title":"Trading account WebSocket system upgrade","state":"completed","begin":"1672823400000","end":"1672825980000","href":"","preOpenBegin":"","scheDesc":"","serviceType":"0","system":"unified","maintType":"1","env":"1","ts":"1672826038470"}]}`))

	select {
	case status := <-gotCh:
		if status.Title == "" || status.Begin != 1672823400000 || status.TS != 1672826038470 {
			t.Fatalf("status = %#v", status)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting status")
	}
}

func TestWSClient_onDataMessage_OrdersAlgo_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan TradeAlgoOrder, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnOrdersAlgo(func(order TradeAlgoOrder) {
		select {
		case gotCh <- order:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"orders-algo","instType":"SPOT","instId":"BTC-USDC"},"data":[{"instType":"SPOT","instId":"BTC-USDC","algoId":"581878926302093312","cTime":"1685002746818","uTime":"1708679675245"}]}`))

	select {
	case order := <-gotCh:
		if order.AlgoId != "581878926302093312" || order.InstId != "BTC-USDC" || order.CTime != 1685002746818 {
			t.Fatalf("order = %#v", order)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting orders-algo")
	}
}

func TestWSClient_onDataMessage_AlgoAdvance_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan TradeAlgoOrder, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnAlgoAdvance(func(order TradeAlgoOrder) {
		select {
		case gotCh <- order:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"algo-advance","instType":"SPOT","instId":"BTC-USDT"},"data":[{"instType":"SPOT","instId":"BTC-USDT","algoId":"a1","cTime":"1685002746818","uTime":"1708679675245"}]}`))

	select {
	case order := <-gotCh:
		if order.AlgoId != "a1" || order.InstId != "BTC-USDT" || order.CTime != 1685002746818 {
			t.Fatalf("order = %#v", order)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting algo-advance")
	}
}

func TestWSClient_onDataMessage_GridOrdersSpot_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan WSGridOrder, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnGridOrdersSpot(func(order WSGridOrder) {
		select {
		case gotCh <- order:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"grid-orders-spot","instType":"SPOT"},"data":[{"algoId":"a1","algoOrdType":"grid","instType":"SPOT","instId":"BTC-USDT","state":"running","cTime":"1681700496249","uTime":"1681700496250"}]}`))

	select {
	case order := <-gotCh:
		if order.AlgoId != "a1" || order.InstId != "BTC-USDT" || order.AlgoOrdType != "grid" || order.CTime != 1681700496249 {
			t.Fatalf("order = %#v", order)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting grid-orders-spot")
	}
}

func TestWSClient_onDataMessage_GridOrdersContract_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan WSGridOrder, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnGridOrdersContract(func(order WSGridOrder) {
		select {
		case gotCh <- order:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"grid-orders-contract","instType":"ANY"},"data":[{"algoId":"a1","algoOrdType":"contract_grid","instType":"SWAP","instId":"BTC-USDT-SWAP","state":"running","cTime":"1682418514204","uTime":"1682418514205"}]}`))

	select {
	case order := <-gotCh:
		if order.AlgoId != "a1" || order.InstId != "BTC-USDT-SWAP" || order.AlgoOrdType != "contract_grid" || order.CTime != 1682418514204 {
			t.Fatalf("order = %#v", order)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting grid-orders-contract")
	}
}

func TestWSClient_onDataMessage_GridSubOrders_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan WSGridSubOrder, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnGridSubOrders(func(order WSGridSubOrder) {
		select {
		case gotCh <- order:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"grid-sub-orders","algoId":"449327675342323712"},"data":[{"algoId":"449327675342323712","algoOrdType":"contract_grid","instType":"SWAP","instId":"BTC-USDT-SWAP","ordId":"449518234142904321","side":"buy","ordType":"limit","state":"live","cTime":"1653445498664","uTime":"1653445498674","pTime":"1653486524502"}]}`))

	select {
	case order := <-gotCh:
		if order.AlgoId != "449327675342323712" || order.OrdId != "449518234142904321" || order.InstId != "BTC-USDT-SWAP" || order.PTime != 1653486524502 {
			t.Fatalf("order = %#v", order)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting grid-sub-orders")
	}
}

func TestWSClient_onDataMessage_GridPositions_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan WSGridPosition, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnGridPositions(func(position WSGridPosition) {
		select {
		case gotCh <- position:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"grid-positions","algoId":"449327675342323712"},"data":[{"algoId":"449327675342323712","adl":"1","instType":"SWAP","instId":"BTC-USDT-SWAP","pos":"35","mgnMode":"cross","posSide":"net","avgPx":"29181.4638888888888895","cTime":"1653400065917","uTime":"1653445498682","pTime":"1653536068723"}]}`))

	select {
	case position := <-gotCh:
		if position.AlgoId != "449327675342323712" || position.InstId != "BTC-USDT-SWAP" || position.Pos != "35" || position.PTime != 1653536068723 {
			t.Fatalf("position = %#v", position)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting grid-positions")
	}
}

func TestWSClient_onDataMessage_AlgoRecurringBuy_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan WSRecurringBuyOrder, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnAlgoRecurringBuy(func(order WSRecurringBuyOrder) {
		select {
		case gotCh <- order:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"algo-recurring-buy","instType":"SPOT"},"data":[{"algoId":"a1","algoOrdType":"recurring","instType":"SPOT","investmentCcy":"USDC","period":"hourly","state":"running","recurringList":[{"ccy":"BTC","ratio":"0.2","px":"36482","avgPx":"0","profit":"0","totalAmt":"0"}],"cTime":"1699932133373","uTime":"1699932136249"}]}`))

	select {
	case order := <-gotCh:
		if order.AlgoId != "a1" || order.InvestmentCcy != "USDC" || order.Period != "hourly" || order.CTime != 1699932133373 {
			t.Fatalf("order = %#v", order)
		}
		if len(order.RecurringList) != 1 || order.RecurringList[0].Ccy != "BTC" {
			t.Fatalf("recurringList = %#v", order.RecurringList)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting algo-recurring-buy")
	}
}

func TestWSClient_onDataMessage_CopyTradingLeadNotification_TypedAsync_Dispatches(t *testing.T) {
	gotCh := make(chan WSCopyTradingLeadNotification, 1)

	w := &WSClient{
		typedAsync: true,
		typedQueue: make(chan wsTypedTask, 1),
	}

	w.OnCopyTradingLeadNotification(func(note WSCopyTradingLeadNotification) {
		select {
		case gotCh <- note:
		default:
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go w.typedDispatchLoop(ctx)

	w.onDataMessage([]byte(`{"arg":{"channel":"copytrading-lead-notification","instType":"SWAP"},"data":[{"infoType":"2","instId":"","instType":"SWAP","maxLeadTraderNum":"3","minLeadEq":"","posSide":"","side":"","subPosId":"667695035433385984","uniqueCode":"3AF72F63E3EAD701"}]}`))

	select {
	case note := <-gotCh:
		if note.InfoType != "2" || note.InstType != "SWAP" || note.SubPosId != "667695035433385984" || note.UniqueCode != "3AF72F63E3EAD701" {
			t.Fatalf("note = %#v", note)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting copytrading-lead-notification")
	}
}
