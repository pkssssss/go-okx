package okx

import (
	"context"
	"fmt"
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
		w.rawDropped.Add(1)
		w.onError(fmt.Errorf("okx: ws raw handler queue full; dropping"))
		return
	}
}
