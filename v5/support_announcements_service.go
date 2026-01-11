package okx

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// SupportAnnouncement 表示公告详情项。
type SupportAnnouncement struct {
	AnnType       string    `json:"annType"`
	Title         string    `json:"title"`
	URL           string    `json:"url"`
	PTime         UnixMilli `json:"pTime"`
	BusinessPTime UnixMilli `json:"businessPTime"`
}

// SupportAnnouncementsPage 表示公告分页结果（每页默认 20 条）。
type SupportAnnouncementsPage struct {
	Details   []SupportAnnouncement `json:"details"`
	TotalPage string                `json:"totalPage"`
}

// SupportAnnouncementsService 获取公告列表。
//
// 注意：OKX 文档说明该接口鉴权“可选”：
// - 若请求头包含 OK-ACCESS-KEY，则被视为私有接口且必须鉴权；
// - 否则视为公共接口，不需要鉴权（按 IP 限制）。
type SupportAnnouncementsService struct {
	c              *Client
	annType        string
	page           string
	acceptLanguage string
	signed         bool
}

// NewSupportAnnouncementsService 创建 SupportAnnouncementsService。
func (c *Client) NewSupportAnnouncementsService() *SupportAnnouncementsService {
	return &SupportAnnouncementsService{c: c}
}

// AnnType 设置公告类型（可选）。
func (s *SupportAnnouncementsService) AnnType(annType string) *SupportAnnouncementsService {
	s.annType = annType
	return s
}

// Page 设置查询页数（可选，默认 1）。
func (s *SupportAnnouncementsService) Page(page string) *SupportAnnouncementsService {
	s.page = page
	return s
}

// AcceptLanguage 设置请求头 Accept-Language（可选：en-US/zh-CN）。
func (s *SupportAnnouncementsService) AcceptLanguage(lang string) *SupportAnnouncementsService {
	s.acceptLanguage = lang
	return s
}

// Signed 指定是否按私有接口签名（默认 false）。
func (s *SupportAnnouncementsService) Signed(signed bool) *SupportAnnouncementsService {
	s.signed = signed
	return s
}

var errEmptySupportAnnouncements = errors.New("okx: empty support announcements response")

// Do 获取公告列表（GET /api/v5/support/announcements）。
func (s *SupportAnnouncementsService) Do(ctx context.Context) (*SupportAnnouncementsPage, error) {
	q := url.Values{}
	if s.annType != "" {
		q.Set("annType", s.annType)
	}
	if s.page != "" {
		q.Set("page", s.page)
	}
	if len(q) == 0 {
		q = nil
	}

	var header http.Header
	if s.acceptLanguage != "" {
		header = make(http.Header)
		header.Set("Accept-Language", s.acceptLanguage)
	}

	var data []SupportAnnouncementsPage
	if err := s.c.doWithHeaders(ctx, http.MethodGet, "/api/v5/support/announcements", q, nil, s.signed, header, &data); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errEmptySupportAnnouncements
	}
	return &data[0], nil
}
