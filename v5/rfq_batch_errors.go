package okx

import "fmt"

// RFQCancelBatchRFQsError 表示 RFQ 批量取消询价单的部分失败（顶层 code=0，但 data[i].sCode!=0）。
type RFQCancelBatchRFQsError struct {
	HTTPStatus  int
	Method      string
	RequestPath string

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
		if ack.SCode != "" && ack.SCode != "0" {
			failed++
			if firstCode == "" {
				firstCode = ack.SCode
				firstMsg = ack.SMsg
			}
		}
	}
	if failed == 0 {
		return fmt.Sprintf("<OKX RFQCancelBatchRFQsError> http=%d method=%s path=%s", e.HTTPStatus, e.Method, e.RequestPath)
	}
	return fmt.Sprintf("<OKX RFQCancelBatchRFQsError> http=%d failed=%d code=%s msg=%s method=%s path=%s", e.HTTPStatus, failed, firstCode, firstMsg, e.Method, e.RequestPath)
}

func rfqCheckCancelBatchRFQs(method, requestPath string, acks []RFQCancelAck) error {
	for _, ack := range acks {
		if ack.SCode != "" && ack.SCode != "0" {
			return &RFQCancelBatchRFQsError{
				HTTPStatus:  200,
				Method:      method,
				RequestPath: requestPath,
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
		if ack.SCode != "" && ack.SCode != "0" {
			failed++
			if firstCode == "" {
				firstCode = ack.SCode
				firstMsg = ack.SMsg
			}
		}
	}
	if failed == 0 {
		return fmt.Sprintf("<OKX RFQCancelBatchQuotesError> http=%d method=%s path=%s", e.HTTPStatus, e.Method, e.RequestPath)
	}
	return fmt.Sprintf("<OKX RFQCancelBatchQuotesError> http=%d failed=%d code=%s msg=%s method=%s path=%s", e.HTTPStatus, failed, firstCode, firstMsg, e.Method, e.RequestPath)
}

func rfqCheckCancelBatchQuotes(method, requestPath string, acks []RFQCancelQuoteAck) error {
	for _, ack := range acks {
		if ack.SCode != "" && ack.SCode != "0" {
			return &RFQCancelBatchQuotesError{
				HTTPStatus:  200,
				Method:      method,
				RequestPath: requestPath,
				Acks:        acks,
			}
		}
	}
	return nil
}
