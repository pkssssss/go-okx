package okx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkssssss/go-okx/v5/internal/sign"
)

const (
	wsPublicURL       = "wss://ws.okx.com:8443/ws/v5/public"
	wsPrivateURL      = "wss://ws.okx.com:8443/ws/v5/private"
	wsBusinessURL     = "wss://ws.okx.com:8443/ws/v5/business"
	wsPublicDemoURL   = "wss://wspap.okx.com:8443/ws/v5/public"
	wsPrivateDemoURL  = "wss://wspap.okx.com:8443/ws/v5/private"
	wsBusinessDemoURL = "wss://wspap.okx.com:8443/ws/v5/business"
)

// WSMessageHandler 处理 WS 原始消息（text/binary 的 payload）。
type WSMessageHandler func(message []byte)

// WSErrorHandler 处理 WS 运行时错误（连接断开、登录失败等）。
type WSErrorHandler func(err error)

// WSEventHandler 处理 WS event 消息（subscribe/login/error/notice 等）。
type WSEventHandler func(event WSEvent)

// WSArg 表示 OKX WS 订阅参数。
// v0.1 仅覆盖通用字段，后续按频道需要扩展。
type WSArg struct {
	Channel     string `json:"channel"`
	InstId      string `json:"instId,omitempty"`
	InstType    string `json:"instType,omitempty"`
	InstFamily  string `json:"instFamily,omitempty"`
	SprdId      string `json:"sprdId,omitempty"`
	Uly         string `json:"uly,omitempty"`
	UID         string `json:"uid,omitempty"`
	Ccy         string `json:"ccy,omitempty"`
	ExtraParams string `json:"extraParams,omitempty"`
}

func (a WSArg) key() string {
	return a.Channel + "|" + a.InstId + "|" + a.InstType + "|" + a.InstFamily + "|" + a.SprdId + "|" + a.Uly + "|" + a.Ccy
}

type wsOpRequest struct {
	ID   string  `json:"id,omitempty"`
	Op   string  `json:"op"`
	Args []WSArg `json:"args"`
}

type wsLoginArg struct {
	APIKey     string `json:"apiKey"`
	Passphrase string `json:"passphrase"`
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
}

type wsLoginRequest struct {
	Op   string       `json:"op"`
	Args []wsLoginArg `json:"args"`
}

// WSOption 用于配置 WSClient。
type WSOption func(*WSClient)

// WithWSURL 覆盖 WS Endpoint（主要用于测试或自定义网关）。
func WithWSURL(url string) WSOption {
	return func(c *WSClient) {
		c.endpoint = url
	}
}

// WithWSHeartbeat 设置 WS 文本心跳间隔。
//
// OKX 建议：若 N 秒内未收到新消息，发送字符串 "ping"，期待收到字符串 "pong" 作为回应；
// 且 N 需小于 30 秒（否则可能被服务器断开）。
//
// 默认启用（25s）。传入 interval<=0 可禁用。
func WithWSHeartbeat(interval time.Duration) WSOption {
	return func(c *WSClient) {
		c.heartbeat = interval
	}
}

// WithWSTypedHandlerAsync 启用 typed handler 的异步分发（orders/fills/account/positions/balance_and_position/deposit-info/withdrawal-info/sprd-orders/sprd-trades/op reply）。
//
// 默认情况下，typed handler 会在 WS read goroutine 中执行；若 handler 逻辑较重可能阻塞读取导致断线/重连。
// 启用该选项后，SDK 会将 typed handler 的执行移动到独立 worker goroutine，并通过有界队列解耦。
//
// 注意：
// - 队列满时会丢弃该条 typed 回调，并通过 errHandler 回调报告错误（调用方需自行调大 buffer 或优化 handler）。
// - raw handler（Start 的 handler 参数）仍在 read goroutine 中执行，如需解耦请避免在 raw handler 中做重逻辑。
func WithWSTypedHandlerAsync(buffer int) WSOption {
	return func(c *WSClient) {
		if buffer <= 0 {
			buffer = 1024
		}
		c.typedAsync = true
		c.typedBuffer = buffer
	}
}

// WithWSHeader 追加/覆盖握手 header。
func WithWSHeader(header http.Header) WSOption {
	return func(c *WSClient) {
		c.header = header.Clone()
	}
}

// WithWSDialer 覆盖 websocket dialer。
func WithWSDialer(d *websocket.Dialer) WSOption {
	return func(c *WSClient) {
		c.dialer = d
	}
}

// WithWSEventHandler 设置 event 消息回调（subscribe/unsubscribe/login/error/notice 等）。
func WithWSEventHandler(handler WSEventHandler) WSOption {
	return func(c *WSClient) {
		c.eventHandler = handler
	}
}

// WSClient 是 OKX WebSocket 客户端（支持 public/private/business）。
// v0.1：实现 ping/pong、private 登录、订阅发送、断线重连与重订阅的基础骨架。
type WSClient struct {
	c         *Client
	endpoint  string
	header    http.Header
	dialer    *websocket.Dialer
	needLogin bool

	heartbeat time.Duration
	lastRecv  atomic.Int64
	lastPing  atomic.Int64

	handler      WSMessageHandler
	errHandler   WSErrorHandler
	eventHandler WSEventHandler

	typedMu                   sync.RWMutex
	ordersHandler             func(order TradeOrder)
	fillsHandler              func(fill WSFill)
	accountHandler            func(balance AccountBalance)
	positionsHandler          func(position AccountPosition)
	balanceAndPositionHandler func(data WSBalanceAndPosition)
	liquidationWarningHandler func(warning WSLiquidationWarning)
	accountGreeksHandler      func(greeks AccountGreeks)
	depositInfoHandler        func(info WSDepositInfo)
	withdrawalInfoHandler     func(info WSWithdrawalInfo)
	sprdOrdersHandler         func(order SprdOrder)
	sprdTradesHandler         func(trade WSSprdTrade)
	tickersHandler            func(ticker MarketTicker)
	tradesHandler             func(trade MarketTrade)
	tradesAllHandler          func(trade MarketTrade)
	orderBookHandler          func(data WSData[WSOrderBook])
	statusHandler             func(status SystemStatus)
	openInterestHandler       func(oi OpenInterest)
	fundingRateHandler        func(rate FundingRate)
	markPriceHandler          func(price MarkPrice)
	indexTickersHandler       func(ticker IndexTicker)
	priceLimitHandler         func(limit PriceLimit)
	optSummaryHandler         func(summary OptSummary)
	liquidationOrdersHandler  func(order LiquidationOrder)
	optionTradesHandler       func(trade WSOptionTrade)
	callAuctionDetailsHandler func(detail WSCallAuctionDetails)
	candlesHandler            func(candle WSCandle)
	priceCandlesHandler       func(candle WSPriceCandle)
	sprdPublicTradesHandler   func(trade WSSprdPublicTrade)
	sprdTickersHandler        func(ticker MarketSprdTicker)
	opReplyHandler            func(reply WSOpReply, raw []byte)

	typedAsync  bool
	typedBuffer int
	typedQueue  chan wsTypedTask

	started atomic.Bool
	cancel  context.CancelFunc
	done    chan struct{}

	opSeq atomic.Uint64

	mu      sync.Mutex
	writeMu sync.Mutex
	conn    *websocket.Conn
	connCh  chan struct{}
	desired map[string]WSArg
	backoff time.Duration

	waitMu  sync.Mutex
	waiters map[string]*wsOpWaiter

	opWaitMu  sync.Mutex
	opWaiters map[string]*wsOpRespWaiter
}

// NewWSPublic 创建 public WS 客户端。
func (c *Client) NewWSPublic(opts ...WSOption) *WSClient {
	endpoint := wsPublicURL
	if c.demo {
		endpoint = wsPublicDemoURL
	}
	w := &WSClient{
		c:         c,
		endpoint:  endpoint,
		connCh:    make(chan struct{}),
		desired:   map[string]WSArg{},
		backoff:   250 * time.Millisecond,
		heartbeat: 25 * time.Second,
		waiters:   map[string]*wsOpWaiter{},
		opWaiters: map[string]*wsOpRespWaiter{},
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// NewWSPrivate 创建 private WS 客户端（需要登录）。
func (c *Client) NewWSPrivate(opts ...WSOption) *WSClient {
	endpoint := wsPrivateURL
	if c.demo {
		endpoint = wsPrivateDemoURL
	}
	w := &WSClient{
		c:         c,
		endpoint:  endpoint,
		needLogin: true,
		connCh:    make(chan struct{}),
		desired:   map[string]WSArg{},
		backoff:   250 * time.Millisecond,
		heartbeat: 25 * time.Second,
		waiters:   map[string]*wsOpWaiter{},
		opWaiters: map[string]*wsOpRespWaiter{},
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// NewWSBusiness 创建 business WS 客户端（是否需要登录由具体频道决定；v0.1 默认不强制登录）。
func (c *Client) NewWSBusiness(opts ...WSOption) *WSClient {
	endpoint := wsBusinessURL
	if c.demo {
		endpoint = wsBusinessDemoURL
	}
	w := &WSClient{
		c:         c,
		endpoint:  endpoint,
		connCh:    make(chan struct{}),
		desired:   map[string]WSArg{},
		backoff:   250 * time.Millisecond,
		heartbeat: 25 * time.Second,
		waiters:   map[string]*wsOpWaiter{},
		opWaiters: map[string]*wsOpRespWaiter{},
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// NewWSBusinessPrivate 创建 business WS 客户端（需要登录）。
//
// 适用场景：资金账户相关推送（如 deposit-info/withdrawal-info）等要求登录的 business 频道。
func (c *Client) NewWSBusinessPrivate(opts ...WSOption) *WSClient {
	endpoint := wsBusinessURL
	if c.demo {
		endpoint = wsBusinessDemoURL
	}
	w := &WSClient{
		c:         c,
		endpoint:  endpoint,
		needLogin: true,
		connCh:    make(chan struct{}),
		desired:   map[string]WSArg{},
		backoff:   250 * time.Millisecond,
		heartbeat: 25 * time.Second,
		waiters:   map[string]*wsOpWaiter{},
		opWaiters: map[string]*wsOpRespWaiter{},
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Start 启动 WS 客户端（后台 goroutine）。
func (w *WSClient) Start(ctx context.Context, handler WSMessageHandler, errHandler WSErrorHandler) error {
	if !w.started.CompareAndSwap(false, true) {
		return errors.New("okx: ws client already started")
	}

	w.handler = handler
	w.errHandler = errHandler
	w.done = make(chan struct{})
	w.touchRecv()

	runCtx, cancel := context.WithCancel(ctx)
	w.cancel = cancel

	if w.typedAsync && w.typedQueue == nil {
		buf := w.typedBuffer
		if buf <= 0 {
			buf = 1024
		}
		w.typedQueue = make(chan wsTypedTask, buf)
		go w.typedDispatchLoop(runCtx)
	}

	if w.heartbeat > 0 {
		go w.heartbeatLoop(runCtx)
	}

	go w.run(runCtx)
	return nil
}

// Close 主动关闭 WS 客户端。
func (w *WSClient) Close() {
	if w.cancel != nil {
		w.cancel()
	}
	// 主动关闭连接以中断 ReadMessage 阻塞，确保 Done() 可及时返回。
	w.closeConn()
}

// Done 返回 WS 客户端停止后的信号通道。
func (w *WSClient) Done() <-chan struct{} {
	if w.done == nil {
		ch := make(chan struct{})
		close(ch)
		return ch
	}
	return w.done
}

// Subscribe 记录订阅并在连接可用时发送（断线后会自动重订阅）。
func (w *WSClient) Subscribe(args ...WSArg) error {
	w.mu.Lock()
	send := make([]WSArg, 0, len(args))
	for _, a := range args {
		if a.Channel == "" {
			w.mu.Unlock()
			return errors.New("okx: ws subscribe requires channel")
		}
		w.desired[a.key()] = a
		send = append(send, a)
	}
	conn := w.conn
	w.mu.Unlock()

	if conn == nil {
		return nil
	}
	return w.writeJSON(conn, wsOpRequest{ID: w.nextOpID(), Op: "subscribe", Args: send})
}

type wsOpWaiter struct {
	op        string
	remaining map[string]struct{}
	done      chan error
}

type wsOpRespResult struct {
	reply *WSOpReply
	raw   []byte
	err   error
}

type wsOpRespWaiter struct {
	op   string
	done chan wsOpRespResult
}

// doOpAndWaitRaw 发送业务 op 请求并等待对应响应（用于 WS 下单/撤单/改单等）。
func (w *WSClient) doOpAndWaitRaw(ctx context.Context, op string, args any) (*WSOpReply, []byte, error) {
	if !w.started.Load() {
		return nil, nil, errors.New("okx: ws client not started")
	}
	if op == "" {
		return nil, nil, errors.New("okx: ws op requires op")
	}
	if args == nil {
		return nil, nil, errors.New("okx: ws op requires args")
	}

	conn, err := w.waitConn(ctx)
	if err != nil {
		return nil, nil, err
	}

	id := w.nextOpID()
	waiter := w.registerOpWaiter(id, op)

	req := struct {
		ID   string `json:"id"`
		Op   string `json:"op"`
		Args any    `json:"args"`
	}{
		ID:   id,
		Op:   op,
		Args: args,
	}

	if err := w.writeJSON(conn, req); err != nil {
		w.removeOpWaiter(id)
		return nil, nil, err
	}

	select {
	case res := <-waiter.done:
		if res.err != nil {
			return nil, nil, res.err
		}
		if res.reply == nil {
			return nil, nil, errors.New("okx: ws op empty reply")
		}
		return res.reply, res.raw, nil
	case <-ctx.Done():
		w.removeOpWaiter(id)
		return nil, nil, ctx.Err()
	}
}

// SubscribeAndWait 发送订阅请求并等待 subscribe/error event 返回（推荐用于判定订阅是否成功）。
func (w *WSClient) SubscribeAndWait(ctx context.Context, args ...WSArg) error {
	if !w.started.Load() {
		return errors.New("okx: ws client not started")
	}

	send := make([]WSArg, 0, len(args))
	for _, a := range args {
		if a.Channel == "" {
			return errors.New("okx: ws subscribe requires channel")
		}
		send = append(send, a)
	}

	conn, err := w.waitConn(ctx)
	if err != nil {
		return err
	}

	id := w.nextOpID()
	waiter := w.registerWaiter(id, "subscribe", send)
	if err := w.writeJSON(conn, wsOpRequest{ID: id, Op: "subscribe", Args: send}); err != nil {
		w.removeWaiter(id)
		return err
	}

	select {
	case err := <-waiter.done:
		if err != nil {
			return err
		}
		w.mu.Lock()
		for _, a := range send {
			w.desired[a.key()] = a
		}
		w.mu.Unlock()
		return nil
	case <-ctx.Done():
		w.removeWaiter(id)
		return ctx.Err()
	}
}

// Unsubscribe 取消订阅并更新本地期望订阅集合（断线重连时不会重订阅）。
func (w *WSClient) Unsubscribe(args ...WSArg) error {
	w.mu.Lock()
	send := make([]WSArg, 0, len(args))
	for _, a := range args {
		if a.Channel == "" {
			w.mu.Unlock()
			return errors.New("okx: ws unsubscribe requires channel")
		}
		delete(w.desired, a.key())
		send = append(send, a)
	}
	conn := w.conn
	w.mu.Unlock()

	if conn == nil {
		return nil
	}
	return w.writeJSON(conn, wsOpRequest{ID: w.nextOpID(), Op: "unsubscribe", Args: send})
}

// UnsubscribeAndWait 发送取消订阅请求并等待 unsubscribe/error event 返回（推荐用于判定取消订阅是否成功）。
func (w *WSClient) UnsubscribeAndWait(ctx context.Context, args ...WSArg) error {
	if !w.started.Load() {
		return errors.New("okx: ws client not started")
	}

	send := make([]WSArg, 0, len(args))
	for _, a := range args {
		if a.Channel == "" {
			return errors.New("okx: ws unsubscribe requires channel")
		}
		send = append(send, a)
	}

	conn, err := w.waitConn(ctx)
	if err != nil {
		return err
	}

	id := w.nextOpID()
	waiter := w.registerWaiter(id, "unsubscribe", send)
	if err := w.writeJSON(conn, wsOpRequest{ID: id, Op: "unsubscribe", Args: send}); err != nil {
		w.removeWaiter(id)
		return err
	}

	select {
	case err := <-waiter.done:
		if err != nil {
			return err
		}
		w.mu.Lock()
		for _, a := range send {
			delete(w.desired, a.key())
		}
		w.mu.Unlock()
		return nil
	case <-ctx.Done():
		w.removeWaiter(id)
		return ctx.Err()
	}
}

func (w *WSClient) run(ctx context.Context) {
	defer close(w.done)

	for {
		if err := ctx.Err(); err != nil {
			w.closeConn()
			return
		}

		conn, err := w.dial(ctx)
		if err != nil {
			w.onError(err)
			w.sleepBackoff(ctx)
			continue
		}

		conn.SetReadLimit(1024 * 1024)

		conn.SetPingHandler(func(appData string) error {
			return w.writeControl(conn, websocket.PongMessage, []byte(appData), 5*time.Second)
		})

		if w.needLogin {
			if err := w.login(ctx, conn); err != nil {
				w.onError(err)
				_ = conn.Close()
				w.sleepBackoff(ctx)
				continue
			}
		}

		w.setConn(conn)

		if args := w.snapshotDesired(); len(args) > 0 {
			_ = w.writeJSON(conn, wsOpRequest{ID: w.nextOpID(), Op: "subscribe", Args: args})
		}

		if err := w.readLoop(ctx, conn); err != nil {
			w.onError(err)
		}

		w.closeConn()
		w.sleepBackoff(ctx)
	}
}

func (w *WSClient) dial(ctx context.Context) (*websocket.Conn, error) {
	d := w.dialer
	if d == nil {
		d = &websocket.Dialer{
			HandshakeTimeout:  45 * time.Second,
			EnableCompression: true,
		}
	}

	header := http.Header{}
	if w.header != nil {
		header = w.header.Clone()
	}

	conn, _, err := d.DialContext(ctx, w.endpoint, header)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (w *WSClient) login(ctx context.Context, conn *websocket.Conn) error {
	if w.c.creds == nil || w.c.creds.APIKey == "" || w.c.creds.SecretKey == "" || w.c.creds.Passphrase == "" {
		return errMissingCredentials
	}

	tm := w.c.now().Add(-w.c.TimeOffset())
	timestamp := sign.TimestampUnixSeconds(tm)
	prehash := sign.PrehashWSLogin(timestamp)
	sig := sign.SignHMACSHA256Base64(w.c.creds.SecretKey, prehash)

	req := wsLoginRequest{
		Op: "login",
		Args: []wsLoginArg{{
			APIKey:     w.c.creds.APIKey,
			Passphrase: w.c.creds.Passphrase,
			Timestamp:  timestamp,
			Sign:       sig,
		}},
	}
	if err := w.writeJSON(conn, req); err != nil {
		return err
	}

	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		w.touchRecv()
		if w.handleHeartbeatMessage(conn, msg) {
			continue
		}

		if w.handler != nil {
			w.safeRawHandlerCall(msg)
		}

		ev, ok, err := WSParseEvent(msg)
		if err != nil || !ok {
			continue
		}
		w.onEvent(*ev)

		switch ev.Event {
		case "login":
			if ev.Code == "0" {
				return nil
			}
			return fmt.Errorf("okx: ws login failed code=%s msg=%s", ev.Code, ev.Msg)
		case "error":
			return fmt.Errorf("okx: ws error code=%s msg=%s", ev.Code, ev.Msg)
		}
	}
}

func (w *WSClient) readLoop(ctx context.Context, conn *websocket.Conn) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		w.touchRecv()
		if w.handleHeartbeatMessage(conn, msg) {
			continue
		}

		if w.handler != nil {
			w.safeRawHandlerCall(msg)
		}

		ev, ok, err := WSParseEvent(msg)
		if err != nil || !ok {
			if r, ok2, err2 := WSParseOpReply(msg); err2 == nil && ok2 {
				w.onOpReply(*r, msg)
			} else {
				w.onDataMessage(msg)
			}
			continue
		}
		w.onEvent(*ev)
		if ev.Event == "notice" && ev.Code == "64008" {
			return errors.New("okx: ws notice 64008 reconnect")
		}
	}
}

func (w *WSClient) onEvent(ev WSEvent) {
	w.notifyWaiter(ev)
	w.notifyOpWaiterError(ev)
	if w.eventHandler != nil && ev.Event != "" {
		w.safeEventHandlerCall(ev)
	}
}

func (w *WSClient) onOpReply(reply WSOpReply, raw []byte) {
	w.notifyOpWaiter(reply, raw)

	w.typedMu.RLock()
	h := w.opReplyHandler
	w.typedMu.RUnlock()
	if h != nil {
		rawCopy := raw
		if w.typedAsync {
			rawCopy = append([]byte(nil), raw...)
		}
		w.dispatchTyped(wsTypedTask{
			kind:  wsTypedKindOpReply,
			op:    reply,
			opRaw: rawCopy,
		})
	}
}

func (w *WSClient) onDataMessage(message []byte) {
	w.typedMu.RLock()
	ordersH := w.ordersHandler
	fillsH := w.fillsHandler
	accountH := w.accountHandler
	positionsH := w.positionsHandler
	balPosH := w.balanceAndPositionHandler
	liqWarningH := w.liquidationWarningHandler
	accountGreeksH := w.accountGreeksHandler
	depInfoH := w.depositInfoHandler
	wdInfoH := w.withdrawalInfoHandler
	sprdOrdersH := w.sprdOrdersHandler
	sprdTradesH := w.sprdTradesHandler
	tickersH := w.tickersHandler
	tradesH := w.tradesHandler
	tradesAllH := w.tradesAllHandler
	orderBookH := w.orderBookHandler
	statusH := w.statusHandler
	openInterestH := w.openInterestHandler
	fundingRateH := w.fundingRateHandler
	markPriceH := w.markPriceHandler
	indexTickersH := w.indexTickersHandler
	priceLimitH := w.priceLimitHandler
	optSummaryH := w.optSummaryHandler
	liquidationOrdersH := w.liquidationOrdersHandler
	optionTradesH := w.optionTradesHandler
	callAuctionDetailsH := w.callAuctionDetailsHandler
	candlesH := w.candlesHandler
	priceCandlesH := w.priceCandlesHandler
	sprdPublicTradesH := w.sprdPublicTradesHandler
	sprdTickersH := w.sprdTickersHandler
	w.typedMu.RUnlock()

	if ordersH == nil &&
		fillsH == nil &&
		accountH == nil &&
		positionsH == nil &&
		balPosH == nil &&
		liqWarningH == nil &&
		accountGreeksH == nil &&
		depInfoH == nil &&
		wdInfoH == nil &&
		sprdOrdersH == nil &&
		sprdTradesH == nil &&
		tickersH == nil &&
		tradesH == nil &&
		tradesAllH == nil &&
		orderBookH == nil &&
		statusH == nil &&
		openInterestH == nil &&
		fundingRateH == nil &&
		markPriceH == nil &&
		indexTickersH == nil &&
		priceLimitH == nil &&
		optSummaryH == nil &&
		liquidationOrdersH == nil &&
		optionTradesH == nil &&
		callAuctionDetailsH == nil &&
		candlesH == nil &&
		priceCandlesH == nil &&
		sprdPublicTradesH == nil &&
		sprdTickersH == nil {
		return
	}

	var probe struct {
		Arg WSArg `json:"arg"`
	}
	if err := json.Unmarshal(message, &probe); err != nil {
		return
	}
	if probe.Arg.Channel == "" {
		return
	}

	switch probe.Arg.Channel {
	case WSChannelOrders:
		if ordersH == nil {
			return
		}
		dm, ok, err := WSParseOrders(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindOrders, orders: dm.Data})
	case WSChannelFills:
		if fillsH == nil {
			return
		}
		dm, ok, err := WSParseFills(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindFills, fills: dm.Data})
	case WSChannelAccount:
		if accountH == nil {
			return
		}
		dm, ok, err := WSParseAccount(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindAccount, balances: dm.Data})
	case WSChannelPositions:
		if positionsH == nil {
			return
		}
		dm, ok, err := WSParsePositions(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindPositions, positions: dm.Data})
	case WSChannelBalanceAndPosition:
		if balPosH == nil {
			return
		}
		dm, ok, err := WSParseBalanceAndPosition(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindBalanceAndPosition, balPos: dm.Data})
	case WSChannelLiquidationWarning:
		if liqWarningH == nil {
			return
		}
		dm, ok, err := WSParseLiquidationWarning(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindLiquidationWarning, liquidationWarnings: dm.Data})
	case WSChannelAccountGreeks:
		if accountGreeksH == nil {
			return
		}
		dm, ok, err := WSParseAccountGreeks(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindAccountGreeks, accountGreeks: dm.Data})
	case WSChannelDepositInfo:
		if depInfoH == nil {
			return
		}
		dm, ok, err := WSParseDepositInfo(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindDepositInfo, depositInfo: dm.Data})
	case WSChannelWithdrawalInfo:
		if wdInfoH == nil {
			return
		}
		dm, ok, err := WSParseWithdrawalInfo(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindWithdrawalInfo, withdrawalInfo: dm.Data})
	case WSChannelSprdOrders:
		if sprdOrdersH == nil {
			return
		}
		dm, ok, err := WSParseSprdOrders(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindSprdOrders, sprdOrders: dm.Data})
	case WSChannelSprdTrades:
		if sprdTradesH == nil {
			return
		}
		dm, ok, err := WSParseSprdTrades(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindSprdTrades, sprdTrades: dm.Data})
	case WSChannelTickers:
		if tickersH == nil {
			return
		}
		dm, ok, err := WSParseTickers(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindTickers, tickers: dm.Data})
	case WSChannelTrades:
		if tradesH == nil {
			return
		}
		dm, ok, err := WSParseTrades(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindTrades, trades: dm.Data})
	case WSChannelTradesAll:
		if tradesAllH == nil {
			return
		}
		dm, ok, err := WSParseTradesAll(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindTradesAll, tradesAll: dm.Data})
	case WSChannelStatus:
		if statusH == nil {
			return
		}
		dm, ok, err := WSParseStatus(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindStatus, statuses: dm.Data})
	case WSChannelOpenInterest:
		if openInterestH == nil {
			return
		}
		dm, ok, err := WSParseOpenInterest(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindOpenInterest, openInterests: dm.Data})
	case WSChannelFundingRate:
		if fundingRateH == nil {
			return
		}
		dm, ok, err := WSParseFundingRate(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindFundingRate, fundingRates: dm.Data})
	case WSChannelMarkPrice:
		if markPriceH == nil {
			return
		}
		dm, ok, err := WSParseMarkPrice(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindMarkPrice, markPrices: dm.Data})
	case WSChannelIndexTickers:
		if indexTickersH == nil {
			return
		}
		dm, ok, err := WSParseIndexTickers(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindIndexTickers, indexTickers: dm.Data})
	case WSChannelPriceLimit:
		if priceLimitH == nil {
			return
		}
		dm, ok, err := WSParsePriceLimit(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindPriceLimit, priceLimits: dm.Data})
	case WSChannelOptSummary:
		if optSummaryH == nil {
			return
		}
		dm, ok, err := WSParseOptSummary(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindOptSummary, optSummaries: dm.Data})
	case WSChannelLiquidationOrders:
		if liquidationOrdersH == nil {
			return
		}
		dm, ok, err := WSParseLiquidationOrders(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindLiquidationOrders, liquidationOrders: dm.Data})
	case WSChannelOptionTrades:
		if optionTradesH == nil {
			return
		}
		dm, ok, err := WSParseOptionTrades(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindOptionTrades, optionTrades: dm.Data})
	case WSChannelCallAuctionDetails:
		if callAuctionDetailsH == nil {
			return
		}
		dm, ok, err := WSParseCallAuctionDetails(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindCallAuctionDetails, callAuctionDetails: dm.Data})
	case WSChannelSprdPublicTrades:
		if sprdPublicTradesH == nil {
			return
		}
		dm, ok, err := WSParseSprdPublicTrades(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindSprdPublicTrades, sprdPublicTrades: dm.Data})
	case WSChannelSprdTickers:
		if sprdTickersH == nil {
			return
		}
		dm, ok, err := WSParseSprdTickers(message)
		if err != nil || !ok || len(dm.Data) == 0 {
			return
		}
		w.dispatchTyped(wsTypedTask{kind: wsTypedKindSprdTickers, sprdTickers: dm.Data})
	default:
		channel := probe.Arg.Channel

		if orderBookH != nil && isOrderBookChannel(channel) {
			dm, ok, err := WSParseOrderBook(message)
			if err != nil || !ok || len(dm.Data) == 0 {
				return
			}
			w.dispatchTyped(wsTypedTask{kind: wsTypedKindOrderBook, orderBooks: []WSData[WSOrderBook]{*dm}})
			return
		}

		if candlesH != nil {
			switch {
			case isCandleChannel(channel):
				dm, ok, err := WSParseCandles(message)
				if err != nil || !ok || len(dm.Data) == 0 {
					return
				}
				out := make([]WSCandle, 0, len(dm.Data))
				for _, c := range dm.Data {
					cc := c
					out = append(out, WSCandle{Arg: dm.Arg, Candle: cc})
				}
				w.dispatchTyped(wsTypedTask{kind: wsTypedKindCandles, candles: out})
				return
			case isSprdCandleChannel(channel):
				dm, ok, err := WSParseSprdCandles(message)
				if err != nil || !ok || len(dm.Data) == 0 {
					return
				}
				out := make([]WSCandle, 0, len(dm.Data))
				for _, c := range dm.Data {
					cc := c
					out = append(out, WSCandle{Arg: dm.Arg, Candle: cc})
				}
				w.dispatchTyped(wsTypedTask{kind: wsTypedKindCandles, candles: out})
				return
			default:
				// continue
			}
		}

		if priceCandlesH != nil {
			switch {
			case isMarkPriceCandleChannel(channel):
				dm, ok, err := WSParseMarkPriceCandles(message)
				if err != nil || !ok || len(dm.Data) == 0 {
					return
				}
				out := make([]WSPriceCandle, 0, len(dm.Data))
				for _, c := range dm.Data {
					cc := c
					out = append(out, WSPriceCandle{Arg: dm.Arg, Candle: cc})
				}
				w.dispatchTyped(wsTypedTask{kind: wsTypedKindPriceCandles, priceCandles: out})
				return
			case isIndexCandleChannel(channel):
				dm, ok, err := WSParseIndexCandles(message)
				if err != nil || !ok || len(dm.Data) == 0 {
					return
				}
				out := make([]WSPriceCandle, 0, len(dm.Data))
				for _, c := range dm.Data {
					cc := c
					out = append(out, WSPriceCandle{Arg: dm.Arg, Candle: cc})
				}
				w.dispatchTyped(wsTypedTask{kind: wsTypedKindPriceCandles, priceCandles: out})
				return
			default:
				// continue
			}
		}

		return
	}
}

func (w *WSClient) writeJSON(conn *websocket.Conn, v any) error {
	w.writeMu.Lock()
	defer w.writeMu.Unlock()
	return conn.WriteJSON(v)
}

func (w *WSClient) writeText(conn *websocket.Conn, message string) error {
	w.writeMu.Lock()
	defer w.writeMu.Unlock()
	return conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (w *WSClient) writeControl(conn *websocket.Conn, messageType int, data []byte, timeout time.Duration) error {
	w.writeMu.Lock()
	defer w.writeMu.Unlock()
	return conn.WriteControl(messageType, data, time.Now().Add(timeout))
}

func (w *WSClient) onError(err error) {
	if w.errHandler != nil && err != nil && !errors.Is(err, context.Canceled) {
		defer func() {
			_ = recover() // 避免用户 errHandler panic 影响 WS 主循环
		}()
		w.errHandler(err)
	}
}

func (w *WSClient) safeRawHandlerCall(message []byte) {
	defer func() {
		if r := recover(); r != nil {
			w.onError(fmt.Errorf("okx: ws raw handler panic: %v", r))
		}
	}()
	w.handler(message)
}

func (w *WSClient) safeEventHandlerCall(event WSEvent) {
	defer func() {
		if r := recover(); r != nil {
			w.onError(fmt.Errorf("okx: ws event handler panic: %v", r))
		}
	}()
	w.eventHandler(event)
}

func (w *WSClient) touchRecv() {
	if w == nil {
		return
	}
	w.lastRecv.Store(time.Now().UnixNano())
}

func isWSPongMessage(message []byte) bool {
	m := bytes.TrimSpace(message)
	return bytes.Equal(m, []byte("pong")) || bytes.Equal(m, []byte(`"pong"`))
}

func isWSPingMessage(message []byte) bool {
	m := bytes.TrimSpace(message)
	return bytes.Equal(m, []byte("ping")) || bytes.Equal(m, []byte(`"ping"`))
}

// handleHeartbeatMessage 处理 OKX 文本 ping/pong 心跳消息，返回 true 表示已处理且应跳过后续 JSON 解析/回调。
func (w *WSClient) handleHeartbeatMessage(conn *websocket.Conn, message []byte) bool {
	if isWSPongMessage(message) {
		return true
	}
	if isWSPingMessage(message) {
		if err := w.writeText(conn, "pong"); err != nil {
			w.onError(err)
			w.closeConn()
		}
		return true
	}
	return false
}

func (w *WSClient) heartbeatLoop(ctx context.Context) {
	interval := w.heartbeat
	if interval <= 0 {
		return
	}

	pongTimeout := 10 * time.Second
	if interval < pongTimeout {
		pongTimeout = interval
	}

	checkInterval := interval / 2
	if checkInterval < time.Second {
		checkInterval = time.Second
	}

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		now := time.Now()

		lastRecvNano := w.lastRecv.Load()
		if lastRecvNano == 0 {
			w.lastRecv.Store(now.UnixNano())
			continue
		}

		lastRecv := time.Unix(0, lastRecvNano)
		if now.Sub(lastRecv) < interval {
			continue
		}

		lastPingNano := w.lastPing.Load()
		if lastPingNano != 0 && lastRecvNano < lastPingNano {
			lastPing := time.Unix(0, lastPingNano)
			if now.Sub(lastPing) > pongTimeout {
				w.onError(errors.New("okx: ws heartbeat pong timeout"))
				w.closeConn()
			}
			continue
		}

		w.mu.Lock()
		conn := w.conn
		w.mu.Unlock()
		if conn == nil {
			continue
		}

		// 二次判定：避免刚收到消息就发 ping
		if now.Sub(time.Unix(0, w.lastRecv.Load())) < interval {
			continue
		}

		if err := w.writeText(conn, "ping"); err != nil {
			w.onError(err)
			w.closeConn()
			continue
		}
		w.lastPing.Store(now.UnixNano())
	}
}

func (w *WSClient) sleepBackoff(ctx context.Context) {
	d := w.backoff
	if d <= 0 {
		d = 250 * time.Millisecond
	}
	if d > 10*time.Second {
		d = 10 * time.Second
	}
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return
	case <-timer.C:
	}

	if w.backoff < 10*time.Second {
		w.backoff *= 2
	}
}

func (w *WSClient) setConn(conn *websocket.Conn) {
	w.mu.Lock()
	w.conn = conn
	w.backoff = 250 * time.Millisecond
	w.notifyConnChangeLocked()
	w.mu.Unlock()

	w.lastPing.Store(0)
	w.touchRecv()
}

func (w *WSClient) closeConn() {
	w.mu.Lock()
	conn := w.conn
	w.conn = nil
	w.notifyConnChangeLocked()
	w.mu.Unlock()
	if conn != nil {
		_ = conn.Close()
	}
	w.lastPing.Store(0)
	w.failOpWaiters(errors.New("okx: ws disconnected"))
	w.failWaiters(errors.New("okx: ws disconnected"))
}

func (w *WSClient) snapshotDesired() []WSArg {
	w.mu.Lock()
	defer w.mu.Unlock()
	out := make([]WSArg, 0, len(w.desired))
	for _, a := range w.desired {
		out = append(out, a)
	}
	return out
}

func (w *WSClient) notifyConnChangeLocked() {
	close(w.connCh)
	w.connCh = make(chan struct{})
}

func (w *WSClient) waitConn(ctx context.Context) (*websocket.Conn, error) {
	for {
		w.mu.Lock()
		conn := w.conn
		ch := w.connCh
		w.mu.Unlock()

		if conn != nil {
			return conn, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ch:
		}
	}
}

func (w *WSClient) nextOpID() string {
	seq := w.opSeq.Add(1)
	return strconv.FormatUint(seq, 10)
}

func (w *WSClient) registerWaiter(id string, op string, args []WSArg) *wsOpWaiter {
	waiter := &wsOpWaiter{
		op:        op,
		remaining: make(map[string]struct{}, len(args)),
		done:      make(chan error, 1),
	}
	for _, a := range args {
		waiter.remaining[a.key()] = struct{}{}
	}

	w.waitMu.Lock()
	w.waiters[id] = waiter
	w.waitMu.Unlock()
	return waiter
}

func (w *WSClient) removeWaiter(id string) {
	w.waitMu.Lock()
	delete(w.waiters, id)
	w.waitMu.Unlock()
}

func (w *WSClient) notifyWaiter(ev WSEvent) {
	if ev.ID == "" {
		return
	}

	w.waitMu.Lock()
	waiter := w.waiters[ev.ID]
	if waiter == nil {
		w.waitMu.Unlock()
		return
	}

	if ev.Event == "error" {
		delete(w.waiters, ev.ID)
		w.waitMu.Unlock()
		waiter.done <- fmt.Errorf("okx: ws op=%s id=%s code=%s msg=%s", waiter.op, ev.ID, ev.Code, ev.Msg)
		return
	}

	if ev.Event != waiter.op || ev.Arg == nil {
		w.waitMu.Unlock()
		return
	}

	delete(waiter.remaining, ev.Arg.key())
	if len(waiter.remaining) != 0 {
		w.waitMu.Unlock()
		return
	}
	delete(w.waiters, ev.ID)
	w.waitMu.Unlock()

	waiter.done <- nil
}

func (w *WSClient) failWaiters(err error) {
	w.waitMu.Lock()
	waiters := w.waiters
	w.waiters = map[string]*wsOpWaiter{}
	w.waitMu.Unlock()

	for id, waiter := range waiters {
		_ = id
		select {
		case waiter.done <- err:
		default:
		}
	}
}

func (w *WSClient) registerOpWaiter(id string, op string) *wsOpRespWaiter {
	waiter := &wsOpRespWaiter{
		op:   op,
		done: make(chan wsOpRespResult, 1),
	}
	w.opWaitMu.Lock()
	w.opWaiters[id] = waiter
	w.opWaitMu.Unlock()
	return waiter
}

func (w *WSClient) removeOpWaiter(id string) {
	w.opWaitMu.Lock()
	delete(w.opWaiters, id)
	w.opWaitMu.Unlock()
}

func (w *WSClient) notifyOpWaiter(reply WSOpReply, raw []byte) {
	if reply.ID == "" {
		return
	}

	w.opWaitMu.Lock()
	waiter := w.opWaiters[reply.ID]
	if waiter == nil {
		w.opWaitMu.Unlock()
		return
	}
	delete(w.opWaiters, reply.ID)
	w.opWaitMu.Unlock()

	waiter.done <- wsOpRespResult{reply: &reply, raw: raw}
}

func (w *WSClient) notifyOpWaiterError(ev WSEvent) {
	if ev.Event != "error" || ev.ID == "" {
		return
	}

	w.opWaitMu.Lock()
	waiter := w.opWaiters[ev.ID]
	if waiter == nil {
		w.opWaitMu.Unlock()
		return
	}
	delete(w.opWaiters, ev.ID)
	w.opWaitMu.Unlock()

	waiter.done <- wsOpRespResult{err: fmt.Errorf("okx: ws op=%s id=%s code=%s msg=%s", waiter.op, ev.ID, ev.Code, ev.Msg)}
}

func (w *WSClient) failOpWaiters(err error) {
	w.opWaitMu.Lock()
	waiters := w.opWaiters
	w.opWaiters = map[string]*wsOpRespWaiter{}
	w.opWaitMu.Unlock()

	for _, waiter := range waiters {
		select {
		case waiter.done <- wsOpRespResult{err: err}:
		default:
		}
	}
}
