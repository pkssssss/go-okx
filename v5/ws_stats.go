package okx

import (
	"time"
)

// WSLastError 表示 WS 客户端最近一次记录的运行错误（用于监控/排障）。
type WSLastError struct {
	Time    time.Time
	Message string
}

// WSStats 是 WSClient 的运行状态快照（可用于监控/告警/自检）。
type WSStats struct {
	Endpoint  string
	Kind      string
	NeedLogin bool

	Started   bool
	Connected bool

	LastRecv  time.Time
	LastPing  time.Time
	LastError WSLastError

	DialAttempts uint64
	Connects     uint64
	Reconnects   uint64

	SubscribeOK      uint64
	SubscribeError   uint64
	UnsubscribeOK    uint64
	UnsubscribeError uint64

	DesiredSubscriptions int
	Waiters              int
	OpWaiters            int

	TypedQueueLen int
	TypedQueueCap int
	RawQueueLen   int
	RawQueueCap   int

	Backoff time.Duration
}

// Stats 返回 WSClient 的运行状态快照（并发安全）。
func (w *WSClient) Stats() WSStats {
	var s WSStats
	if w == nil {
		return s
	}

	s.Endpoint = w.endpoint
	s.NeedLogin = w.needLogin
	s.Started = w.started.Load()

	switch w.kind {
	case wsKindPublic:
		s.Kind = "public"
	case wsKindPrivate:
		s.Kind = "private"
	case wsKindBusiness:
		s.Kind = "business"
	default:
		s.Kind = "unknown"
	}

	w.mu.Lock()
	s.Connected = w.conn != nil
	s.DesiredSubscriptions = len(w.desired)
	s.Backoff = w.backoff
	w.mu.Unlock()

	if ns := w.lastRecv.Load(); ns != 0 {
		s.LastRecv = time.Unix(0, ns)
	}
	if ns := w.lastPing.Load(); ns != 0 {
		s.LastPing = time.Unix(0, ns)
	}

	if v := w.lastError.Load(); v != nil {
		if le, ok := v.(WSLastError); ok {
			s.LastError = le
		}
	}

	s.DialAttempts = w.dialAttempts.Load()
	s.Connects = w.connects.Load()
	s.Reconnects = w.reconnects.Load()
	s.SubscribeOK = w.subscribeOK.Load()
	s.SubscribeError = w.subscribeErr.Load()
	s.UnsubscribeOK = w.unsubscribeOK.Load()
	s.UnsubscribeError = w.unsubscribeErr.Load()

	if q := w.typedQueue; q != nil {
		s.TypedQueueLen = len(q)
		s.TypedQueueCap = cap(q)
	}
	if q := w.rawQueue; q != nil {
		s.RawQueueLen = len(q)
		s.RawQueueCap = cap(q)
	}

	w.waitMu.Lock()
	s.Waiters = len(w.waiters)
	w.waitMu.Unlock()

	w.opWaitMu.Lock()
	s.OpWaiters = len(w.opWaiters)
	w.opWaitMu.Unlock()

	return s
}
