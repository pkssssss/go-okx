package okx

import (
	"context"
)

func (w *WSClient) rawDispatchLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-w.rawQueue:
			if w.handler == nil {
				continue
			}
			w.safeRawHandlerCall(msg)
		}
	}
}

func (w *WSClient) dispatchRaw(message []byte) {
	if w == nil || w.handler == nil {
		return
	}

	if !w.rawAsync || w.rawQueue == nil {
		w.safeRawHandlerCall(message)
		return
	}

	select {
	case w.rawQueue <- message:
		return
	default:
	}

	policy := w.rawQueueFullPolicy
	w.warnQueueFull(&w.rawQueueFullWarnAt, &WSQueueFullError{
		Queue:    "raw",
		Policy:   policy,
		QueueLen: len(w.rawQueue),
		QueueCap: cap(w.rawQueue),
	})

	switch policy {
	case WSQueueFullDrop:
		w.rawDropped.Add(1)
		return
	case WSQueueFullDisconnect:
		w.rawDropped.Add(1)
		w.closeConn()
		return
	default: // WSQueueFullBlock
		done := w.ctxDone
		if done == nil {
			w.rawQueue <- message
			return
		}
		select {
		case w.rawQueue <- message:
			return
		case <-done:
			return
		}
	}
}
