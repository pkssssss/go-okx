package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// AffiliateInviteeDetail 表示被邀请人返佣信息。
// 数值字段保持为 string（无损）。
type AffiliateInviteeDetail struct {
	InviteeLevel      string    `json:"inviteeLevel"`
	JoinTime          UnixMilli `json:"joinTime"`
	InviteeRebateRate string    `json:"inviteeRebateRate"`
	TotalCommission   string    `json:"totalCommission"`

	FirstTradeTime UnixMilli `json:"firstTradeTime"`
	Level          string    `json:"level"`
	DepAmt         string    `json:"depAmt"`
	VolMonth       string    `json:"volMonth"`
	AccFee         string    `json:"accFee"`
	KycTime        UnixMilli `json:"kycTime"`
	Region         string    `json:"region"`
	AffiliateCode  string    `json:"affiliateCode"`
}

// AffiliateInviteeDetailService 获取被邀请人返佣信息。
type AffiliateInviteeDetailService struct {
	c   *Client
	uid string
}

// NewAffiliateInviteeDetailService 创建 AffiliateInviteeDetailService。
func (c *Client) NewAffiliateInviteeDetailService() *AffiliateInviteeDetailService {
	return &AffiliateInviteeDetailService{c: c}
}

// UID 设置被邀请人 UID（必填）。
func (s *AffiliateInviteeDetailService) UID(uid string) *AffiliateInviteeDetailService {
	s.uid = uid
	return s
}

var (
	errAffiliateInviteeDetailMissingUID = errors.New("okx: affiliate invitee detail requires uid")
	errEmptyAffiliateInviteeDetail      = errors.New("okx: empty affiliate invitee detail response")
)

// Do 获取被邀请人返佣信息（GET /api/v5/affiliate/invitee/detail）。
func (s *AffiliateInviteeDetailService) Do(ctx context.Context) (*AffiliateInviteeDetail, error) {
	if s.uid == "" {
		return nil, errAffiliateInviteeDetailMissingUID
	}

	q := url.Values{}
	q.Set("uid", s.uid)

	var data []AffiliateInviteeDetail
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/affiliate/invitee/detail", q, nil, true, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptyAffiliateInviteeDetail
	}
	return &data[0], nil
}
