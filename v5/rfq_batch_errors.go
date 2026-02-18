package okx

import "fmt"

// RFQCancelBatchRFQsError 表示 RFQ 批量取消询价单的部分失败（顶层 code=0，但 data[i].sCode!=0）。
type RFQCancelBatchRFQsError struct {
	HTTPStatus  int
	Method      string
	RequestPath string
	RequestID   string
	Expected    int

	Acks []RFQCancelAck
}

func (e *RFQCancelBatchRFQsError) Error() string {
	if e == nil {
		return "<OKX RFQCancelBatchRFQsError>"
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
		return fmt.Sprintf("<OKX RFQCancelBatchRFQsError> http=%d expected=%d actual=%d method=%s path=%s%s", e.HTTPStatus, e.Expected, len(e.Acks), e.Method, e.RequestPath, requestIDPart)
	}
	if failed == 0 {
		return fmt.Sprintf("<OKX RFQCancelBatchRFQsError> http=%d method=%s path=%s%s", e.HTTPStatus, e.Method, e.RequestPath, requestIDPart)
	}
	return fmt.Sprintf("<OKX RFQCancelBatchRFQsError> http=%d failed=%d code=%s msg=%s method=%s path=%s%s", e.HTTPStatus, failed, firstCode, firstMsg, e.Method, e.RequestPath, requestIDPart)
}

func rfqCheckCancelBatchRFQs(method, requestPath, requestID string, expectedCount int, acks []RFQCancelAck) error {
	if len(acks) == 0 || (expectedCount > 0 && len(acks) != expectedCount) {
		return &RFQCancelBatchRFQsError{
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
			return &RFQCancelBatchRFQsError{
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

// RFQCancelBatchQuotesError 表示 RFQ 批量取消报价单的部分失败（顶层 code=0，但 data[i].sCode!=0）。
type RFQCancelBatchQuotesError struct {
	HTTPStatus  int
	Method      string
	RequestPath string
	RequestID   string
	Expected    int

	Acks []RFQCancelQuoteAck
}

func (e *RFQCancelBatchQuotesError) Error() string {
	if e == nil {
		return "<OKX RFQCancelBatchQuotesError>"
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
		return fmt.Sprintf("<OKX RFQCancelBatchQuotesError> http=%d expected=%d actual=%d method=%s path=%s%s", e.HTTPStatus, e.Expected, len(e.Acks), e.Method, e.RequestPath, requestIDPart)
	}
	if failed == 0 {
		return fmt.Sprintf("<OKX RFQCancelBatchQuotesError> http=%d method=%s path=%s%s", e.HTTPStatus, e.Method, e.RequestPath, requestIDPart)
	}
	return fmt.Sprintf("<OKX RFQCancelBatchQuotesError> http=%d failed=%d code=%s msg=%s method=%s path=%s%s", e.HTTPStatus, failed, firstCode, firstMsg, e.Method, e.RequestPath, requestIDPart)
}

func rfqCheckCancelBatchQuotes(method, requestPath, requestID string, expectedCount int, acks []RFQCancelQuoteAck) error {
	if len(acks) == 0 || (expectedCount > 0 && len(acks) != expectedCount) {
		return &RFQCancelBatchQuotesError{
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
			return &RFQCancelBatchQuotesError{
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
