package okx

import "fmt"

// AccountMovePositionsError 表示 account/move-positions 的“部分失败/状态失败”。
//
// 约定：
// - OKX 顶层 code=0 仅表示请求被处理，不代表移仓全部成功；
// - 需要结合 ack.state 以及 legs[from/to].sCode 判定最终执行结果。
//
// 该错误会尽量携带 Ack 以便排障与重试决策。
type AccountMovePositionsError struct {
	HTTPStatus  int
	Method      string
	RequestPath string
	RequestID   string

	Ack *AccountMovePositionsAck
}

func (e *AccountMovePositionsError) Error() string {
	if e == nil {
		return "<OKX AccountMovePositionsError>"
	}

	requestIDPart := ""
	if e.RequestID != "" {
		requestIDPart = " requestId=" + e.RequestID
	}

	state := ""
	legsFailed := 0
	firstCode := ""
	firstMsg := ""
	if e.Ack != nil {
		state = e.Ack.State
		for _, leg := range e.Ack.Legs {
			if leg.From.SCode != "0" {
				legsFailed++
				if firstCode == "" {
					firstCode = leg.From.SCode
					firstMsg = leg.From.SMsg
				}
			}
			if leg.To.SCode != "0" {
				legsFailed++
				if firstCode == "" {
					firstCode = leg.To.SCode
					firstMsg = leg.To.SMsg
				}
			}
		}
	}

	if legsFailed == 0 && (state == "" || state == "filled") {
		return fmt.Sprintf("<OKX AccountMovePositionsError> http=%d method=%s path=%s%s", e.HTTPStatus, e.Method, e.RequestPath, requestIDPart)
	}

	statePart := ""
	if state != "" {
		statePart = " state=" + state
	}
	legsPart := ""
	if legsFailed > 0 {
		legsPart = fmt.Sprintf(" legsFailed=%d sCode=%s sMsg=%s", legsFailed, firstCode, firstMsg)
	}
	return fmt.Sprintf("<OKX AccountMovePositionsError> http=%d%s%s method=%s path=%s%s", e.HTTPStatus, statePart, legsPart, e.Method, e.RequestPath, requestIDPart)
}

func accountCheckMovePositionsAck(method, requestPath, requestID string, ack *AccountMovePositionsAck) error {
	if ack == nil {
		return &AccountMovePositionsError{
			HTTPStatus:  200,
			Method:      method,
			RequestPath: requestPath,
			RequestID:   requestID,
		}
	}
	if ack.State != "filled" {
		return &AccountMovePositionsError{
			HTTPStatus:  200,
			Method:      method,
			RequestPath: requestPath,
			RequestID:   requestID,
			Ack:         ack,
		}
	}

	for _, leg := range ack.Legs {
		if leg.From.SCode != "0" || leg.To.SCode != "0" {
			return &AccountMovePositionsError{
				HTTPStatus:  200,
				Method:      method,
				RequestPath: requestPath,
				RequestID:   requestID,
				Ack:         ack,
			}
		}
	}
	return nil
}
