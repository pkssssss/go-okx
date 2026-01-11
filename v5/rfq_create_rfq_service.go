package okx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const rfqMaxLegs = 15
const rfqMaxAcctAlloc = 10

type rfqCreateRFQRequest struct {
	Counterparties        []string       `json:"counterparties"`
	Anonymous             *bool          `json:"anonymous,omitempty"`
	AllowPartialExecution *bool          `json:"allowPartialExecution,omitempty"`
	ClRfqId               string         `json:"clRfqId,omitempty"`
	Tag                   string         `json:"tag,omitempty"`
	Legs                  []RFQLeg       `json:"legs"`
	AcctAlloc             []RFQAcctAlloc `json:"acctAlloc,omitempty"`
}

// RFQCreateRFQService 创建询价单。
type RFQCreateRFQService struct {
	c *Client

	counterparties        []string
	anonymous             *bool
	allowPartialExecution *bool
	clRfqId               string
	tag                   string

	legs      []RFQLeg
	acctAlloc []RFQAcctAlloc
}

// NewRFQCreateRFQService 创建 RFQCreateRFQService。
func (c *Client) NewRFQCreateRFQService() *RFQCreateRFQService {
	return &RFQCreateRFQService{c: c}
}

// Counterparties 设置希望收到询价的报价方列表（必填）。
func (s *RFQCreateRFQService) Counterparties(counterparties []string) *RFQCreateRFQService {
	s.counterparties = counterparties
	return s
}

// Anonymous 设置是否匿名询价（可选；默认 false）。
func (s *RFQCreateRFQService) Anonymous(anonymous bool) *RFQCreateRFQService {
	s.anonymous = &anonymous
	return s
}

// AllowPartialExecution 设置是否允许部分执行（可选；默认 false）。
func (s *RFQCreateRFQService) AllowPartialExecution(allow bool) *RFQCreateRFQService {
	s.allowPartialExecution = &allow
	return s
}

// ClRfqId 设置询价单自定义 ID（可选）。
func (s *RFQCreateRFQService) ClRfqId(clRfqId string) *RFQCreateRFQService {
	s.clRfqId = clRfqId
	return s
}

// Tag 设置询价单标签（可选）。
func (s *RFQCreateRFQService) Tag(tag string) *RFQCreateRFQService {
	s.tag = tag
	return s
}

// Legs 设置询价单腿列表（必填，最多 15 条）。
func (s *RFQCreateRFQService) Legs(legs []RFQLeg) *RFQCreateRFQService {
	s.legs = legs
	return s
}

// AcctAlloc 设置组合询价单账户分配（可选，最多 10 个账户）。
func (s *RFQCreateRFQService) AcctAlloc(acctAlloc []RFQAcctAlloc) *RFQCreateRFQService {
	s.acctAlloc = acctAlloc
	return s
}

var (
	errRFQCreateRFQMissingCounterparties = errors.New("okx: rfq create rfq requires at least one counterparty")
	errRFQCreateRFQMissingLegs           = errors.New("okx: rfq create rfq requires at least one leg")
	errRFQCreateRFQTooManyLegs           = errors.New("okx: rfq create rfq max 15 legs")
	errRFQCreateRFQTooManyAcctAlloc      = errors.New("okx: rfq create rfq max 10 acctAlloc")
	errEmptyRFQCreateRFQResponse         = errors.New("okx: empty create rfq response")
)

// Do 创建询价单（POST /api/v5/rfq/create-rfq）。
func (s *RFQCreateRFQService) Do(ctx context.Context) (*RFQ, error) {
	if len(s.counterparties) == 0 {
		return nil, errRFQCreateRFQMissingCounterparties
	}
	if len(s.legs) == 0 {
		return nil, errRFQCreateRFQMissingLegs
	}
	if len(s.legs) > rfqMaxLegs {
		return nil, errRFQCreateRFQTooManyLegs
	}
	if len(s.acctAlloc) > rfqMaxAcctAlloc {
		return nil, errRFQCreateRFQTooManyAcctAlloc
	}

	for i, leg := range s.legs {
		if leg.InstId == "" {
			return nil, fmt.Errorf("okx: rfq create rfq legs[%d] missing instId", i)
		}
		if leg.Sz == "" {
			return nil, fmt.Errorf("okx: rfq create rfq legs[%d] missing sz", i)
		}
		if leg.Side == "" {
			return nil, fmt.Errorf("okx: rfq create rfq legs[%d] missing side", i)
		}
	}

	for i, alloc := range s.acctAlloc {
		if alloc.Acct == "" {
			return nil, fmt.Errorf("okx: rfq create rfq acctAlloc[%d] missing acct", i)
		}
		if len(alloc.Legs) == 0 {
			return nil, fmt.Errorf("okx: rfq create rfq acctAlloc[%d] requires at least one leg", i)
		}
		for j, leg := range alloc.Legs {
			if leg.InstId == "" {
				return nil, fmt.Errorf("okx: rfq create rfq acctAlloc[%d].legs[%d] missing instId", i, j)
			}
			if leg.Sz == "" {
				return nil, fmt.Errorf("okx: rfq create rfq acctAlloc[%d].legs[%d] missing sz", i, j)
			}
		}
	}

	req := rfqCreateRFQRequest{
		Counterparties:        s.counterparties,
		Anonymous:             s.anonymous,
		AllowPartialExecution: s.allowPartialExecution,
		ClRfqId:               s.clRfqId,
		Tag:                   s.tag,
		Legs:                  s.legs,
		AcctAlloc:             s.acctAlloc,
	}

	var data []RFQ
	if err := s.c.do(ctx, http.MethodPost, "/api/v5/rfq/create-rfq", nil, req, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyRFQCreateRFQResponse
	}
	return &data[0], nil
}
