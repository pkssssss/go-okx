package okx

import "fmt"

// TradeAlgoBatchError 表示 OKX 策略委托类批量接口的部分失败（顶层 code=0，但 data[i].sCode!=0）。
type TradeAlgoBatchError struct {
	HTTPStatus  int
	Method      string
	RequestPath string
	RequestID   string

	Acks []TradeAlgoOrderAck
}

func (e *TradeAlgoBatchError) Error() string {
	if e == nil {
		return "<OKX TradeAlgoBatchError>"
	}

	failed := 0
	firstCode := ""
	firstMsg := ""
	for _, ack := range e.Acks {
		if ack.SCode != "0" {
			failed++
			if firstCode == "" {
				firstCode = ack.SCode
				if firstCode == "" {
					firstCode = "<empty>"
					if firstMsg == "" {
						firstMsg = "missing sCode"
					}
				}
				if firstMsg == "" {
					firstMsg = ack.SMsg
				}
			}
		}
	}

	if failed == 0 {
		requestIDPart := ""
		if e.RequestID != "" {
			requestIDPart = " requestId=" + e.RequestID
		}
		return fmt.Sprintf("<OKX TradeAlgoBatchError> http=%d method=%s path=%s%s", e.HTTPStatus, e.Method, e.RequestPath, requestIDPart)
	}
	requestIDPart := ""
	if e.RequestID != "" {
		requestIDPart = " requestId=" + e.RequestID
	}
	return fmt.Sprintf("<OKX TradeAlgoBatchError> http=%d failed=%d code=%s msg=%s method=%s path=%s%s", e.HTTPStatus, failed, firstCode, firstMsg, e.Method, e.RequestPath, requestIDPart)
}

func tradeCheckAlgoAcks(method, requestPath, requestID string, acks []TradeAlgoOrderAck) error {
	if len(acks) == 0 {
		return &TradeAlgoBatchError{
			HTTPStatus:  200,
			Method:      method,
			RequestPath: requestPath,
			RequestID:   requestID,
			Acks:        acks,
		}
	}

	for _, ack := range acks {
		if ack.SCode != "0" {
			return &TradeAlgoBatchError{
				HTTPStatus:  200,
				Method:      method,
				RequestPath: requestPath,
				RequestID:   requestID,
				Acks:        acks,
			}
		}
	}
	return nil
}
