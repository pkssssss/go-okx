package okx

import (
	"context"
	"net/http"
)

// SupportAnnouncementType 表示公告类型。
type SupportAnnouncementType struct {
	AnnType     string `json:"annType"`
	AnnTypeDesc string `json:"annTypeDesc"`
}

// SupportAnnouncementTypesService 获取公告类型列表。
type SupportAnnouncementTypesService struct {
	c *Client
}

// NewSupportAnnouncementTypesService 创建 SupportAnnouncementTypesService。
func (c *Client) NewSupportAnnouncementTypesService() *SupportAnnouncementTypesService {
	return &SupportAnnouncementTypesService{c: c}
}

// Do 获取公告类型（GET /api/v5/support/announcement-types）。
func (s *SupportAnnouncementTypesService) Do(ctx context.Context) ([]SupportAnnouncementType, error) {
	var data []SupportAnnouncementType
	if err := s.c.do(ctx, http.MethodGet, "/api/v5/support/announcement-types", nil, nil, false, &data); err != nil {
		return nil, err
	}
	return data, nil
}
