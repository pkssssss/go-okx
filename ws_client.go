package okx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkssssss/go-okx/internal/sign"
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

// WSArg 表示 OKX WS 订阅参数。
// v0.1 仅覆盖通用字段，后续按频道需要扩展。
type WSArg struct {
	Channel    string `json:"channel"`
	InstId     string `json:"instId,omitempty"`
	InstType   string `json:"instType,omitempty"`
	InstFamily string `json:"instFamily,omitempty"`
	Uly        string `json:"uly,omitempty"`
	UID        string `json:"uid,omitempty"`
}

func (a WSArg) key() string {
	return a.Channel + "|" + a.InstId + "|" + a.InstType + "|" + a.InstFamily + "|" + a.Uly
}

type wsOpRequest struct {
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

type wsEventEnvelope struct {
	Event string `json:"event"`
	Code  string `json:"code"`
	Msg   string `json:"msg"`
}

// WSOption 用于配置 WSClient。
type WSOption func(*WSClient)

// WithWSURL 覆盖 WS Endpoint（主要用于测试或自定义网关）。
func WithWSURL(url string) WSOption {
	return func(c *WSClient) {
		c.endpoint = url
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

// WSClient 是 OKX WebSocket 客户端（支持 public/private/business）。
// v0.1：实现 ping/pong、private 登录、订阅发送、断线重连与重订阅的基础骨架。
type WSClient struct {
	c         *Client
	endpoint  string
	header    http.Header
	dialer    *websocket.Dialer
	needLogin bool

	handler    WSMessageHandler
	errHandler WSErrorHandler

	started atomic.Bool
	cancel  context.CancelFunc
	done    chan struct{}

	mu      sync.Mutex
	writeMu sync.Mutex
	conn    *websocket.Conn
	desired map[string]WSArg
	backoff time.Duration
}

// NewWSPublic 创建 public WS 客户端。
func (c *Client) NewWSPublic(opts ...WSOption) *WSClient {
	endpoint := wsPublicURL
	if c.demo {
		endpoint = wsPublicDemoURL
	}
	w := &WSClient{
		c:        c,
		endpoint: endpoint,
		desired:  map[string]WSArg{},
		backoff:  250 * time.Millisecond,
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
		desired:   map[string]WSArg{},
		backoff:   250 * time.Millisecond,
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
		c:        c,
		endpoint: endpoint,
		desired:  map[string]WSArg{},
		backoff:  250 * time.Millisecond,
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

	runCtx, cancel := context.WithCancel(ctx)
	w.cancel = cancel

	go w.run(runCtx)
	return nil
}

// Close 主动关闭 WS 客户端。
func (w *WSClient) Close() {
	if w.cancel != nil {
		w.cancel()
	}
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
	return w.writeJSON(conn, wsOpRequest{Op: "subscribe", Args: send})
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

		w.setConn(conn)
		conn.SetReadLimit(1024 * 1024)

		conn.SetPingHandler(func(appData string) error {
			return w.writeControl(conn, websocket.PongMessage, []byte(appData), 5*time.Second)
		})

		if w.needLogin {
			if err := w.login(ctx, conn); err != nil {
				w.onError(err)
				w.closeConn()
				w.sleepBackoff(ctx)
				continue
			}
		}

		if args := w.snapshotDesired(); len(args) > 0 {
			_ = w.writeJSON(conn, wsOpRequest{Op: "subscribe", Args: args})
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

		if w.handler != nil {
			w.handler(msg)
		}

		var ev wsEventEnvelope
		if err := json.Unmarshal(msg, &ev); err != nil || ev.Event == "" {
			continue
		}
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

		if w.handler != nil {
			w.handler(msg)
		}

		var ev wsEventEnvelope
		if err := json.Unmarshal(msg, &ev); err == nil {
			if ev.Event == "notice" && ev.Code == "64008" {
				return errors.New("okx: ws notice 64008 reconnect")
			}
		}
	}
}

func (w *WSClient) writeJSON(conn *websocket.Conn, v any) error {
	w.writeMu.Lock()
	defer w.writeMu.Unlock()
	return conn.WriteJSON(v)
}

func (w *WSClient) writeControl(conn *websocket.Conn, messageType int, data []byte, timeout time.Duration) error {
	w.writeMu.Lock()
	defer w.writeMu.Unlock()
	return conn.WriteControl(messageType, data, time.Now().Add(timeout))
}

func (w *WSClient) onError(err error) {
	if w.errHandler != nil && err != nil && !errors.Is(err, context.Canceled) {
		w.errHandler(err)
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
	w.mu.Unlock()
}

func (w *WSClient) closeConn() {
	w.mu.Lock()
	conn := w.conn
	w.conn = nil
	w.mu.Unlock()
	if conn != nil {
		_ = conn.Close()
	}
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
