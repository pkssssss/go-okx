package okx

import "fmt"

// TradeAlgoBatchError 表示 OKX 策略委托类批量接口的部分失败（顶层 code=0，但 data[i].sCode!=0）。
type TradeAlgoBatchError struct {
	HTTPStatus  int
	Method      string
	RequestPath string

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
		if ack.SCode != "" && ack.SCode != "0" {
			failed++
			if firstCode == "" {
				firstCode = ack.SCode
				firstMsg = ack.SMsg
			}
		}
	}

	if failed == 0 {
		return fmt.Sprintf("<OKX TradeAlgoBatchError> http=%d method=%s path=%s", e.HTTPStatus, e.Method, e.RequestPath)
	}
	return fmt.Sprintf("<OKX TradeAlgoBatchError> http=%d failed=%d code=%s msg=%s method=%s path=%s", e.HTTPStatus, failed, firstCode, firstMsg, e.Method, e.RequestPath)
}

func tradeCheckAlgoAcks(method, requestPath string, acks []TradeAlgoOrderAck) error {
	for _, ack := range acks {
		if ack.SCode != "" && ack.SCode != "0" {
			return &TradeAlgoBatchError{
				HTTPStatus:  200,
				Method:      method,
				RequestPath: requestPath,
				Acks:        acks,
			}
		}
	}
	return nil
}
