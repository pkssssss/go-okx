package okx

import "fmt"

// TradingBotBatchError 表示 TradingBot 批量接口的部分失败（顶层 code=0，但 data[i].sCode!=0）。
type TradingBotBatchError struct {
	HTTPStatus  int
	Method      string
	RequestPath string
	RequestID   string
	Expected    int

	Acks []TradingBotOrderAck
}

func (e *TradingBotBatchError) Error() string {
	if e == nil {
		return "<OKX TradingBotBatchError>"
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

	requestIDPart := ""
	if e.RequestID != "" {
		requestIDPart = " requestId=" + e.RequestID
	}
	if e.Expected > 0 && len(e.Acks) != e.Expected {
		return fmt.Sprintf("<OKX TradingBotBatchError> http=%d expected=%d actual=%d method=%s path=%s%s", e.HTTPStatus, e.Expected, len(e.Acks), e.Method, e.RequestPath, requestIDPart)
	}
	if failed == 0 {
		return fmt.Sprintf("<OKX TradingBotBatchError> http=%d method=%s path=%s%s", e.HTTPStatus, e.Method, e.RequestPath, requestIDPart)
	}
	return fmt.Sprintf("<OKX TradingBotBatchError> http=%d failed=%d code=%s msg=%s method=%s path=%s%s", e.HTTPStatus, failed, firstCode, firstMsg, e.Method, e.RequestPath, requestIDPart)
}

func tradingBotCheckBatchAcks(method, requestPath, requestID string, expectedCount int, acks []TradingBotOrderAck) error {
	if len(acks) == 0 || (expectedCount > 0 && len(acks) != expectedCount) {
		return &TradingBotBatchError{
			HTTPStatus:  200,
			Method:      method,
			RequestPath: requestPath,
			RequestID:   requestID,
			Expected:    expectedCount,
			Acks:        acks,
		}
	}

	for _, ack := range acks {
		if ack.SCode != "0" {
			return &TradingBotBatchError{
				HTTPStatus:  200,
				Method:      method,
				RequestPath: requestPath,
				RequestID:   requestID,
				Expected:    expectedCount,
				Acks:        acks,
			}
		}
	}
	return nil
}
