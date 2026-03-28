package services

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/hujia-team/intranet-sdk/client"
	"github.com/hujia-team/intranet-sdk/models"
	"github.com/hujia-team/intranet-sdk/utils"
)

type ClawSkillService interface {
	UploadLocalSkill(rawURL string, archiveName string, archive []byte, version string, uploadToken string, headers map[string]string) (*models.LocalSkillUploadResult, error)
	ResetLocalSkillUploadToken(rawURL string, slug string, headers map[string]string) (*models.LocalSkillTokenResetResult, error)
}

type clawSkillService struct {
	httpClient *client.HTTPClient
}

func NewClawSkillService(httpClient *client.HTTPClient) ClawSkillService {
	return &clawSkillService{httpClient: httpClient}
}

func (s *clawSkillService) UploadLocalSkill(rawURL string, archiveName string, archive []byte, version string, uploadToken string, headers map[string]string) (*models.LocalSkillUploadResult, error) {
	fields := map[string]string{}
	if strings.TrimSpace(version) != "" {
		fields["version"] = strings.TrimSpace(version)
	}
	reqHeaders := map[string]string{}
	for key, value := range headers {
		reqHeaders[key] = value
	}
	if strings.TrimSpace(uploadToken) != "" {
		reqHeaders["X-Skill-Token"] = strings.TrimSpace(uploadToken)
	}
	body, contentType, err := client.BuildMultipartBody(fields, "file", filepath.Base(strings.TrimSpace(archiveName)), archive)
	if err != nil {
		return nil, utils.NewInternalError("failed to build upload body", err)
	}

	rawResp, err := s.httpClient.PostMultipartRawURL(rawURL, body, contentType, reqHeaders)
	result := &models.LocalSkillUploadResult{}
	if rawResp != nil {
		result.StatusCode = rawResp.StatusCode
		result.BodyText = strings.TrimSpace(string(rawResp.Body))
		if len(rawResp.Body) > 0 {
			var parsed models.LocalSkillUploadResponse
			if parseErr := json.Unmarshal(rawResp.Body, &parsed); parseErr == nil {
				result.Parsed = &parsed
			} else {
				result.ParseError = parseErr.Error()
			}
		}
	}
	if err != nil {
		return result, utils.NewAPIError("failed to upload local skill", err)
	}
	return result, nil
}

func (s *clawSkillService) ResetLocalSkillUploadToken(rawURL string, slug string, headers map[string]string) (*models.LocalSkillTokenResetResult, error) {
	req := &models.LocalSkillTokenResetRequest{Slug: strings.TrimSpace(slug)}
	rawResp, err := s.httpClient.PostRawURL(rawURL, req, headers)
	result := &models.LocalSkillTokenResetResult{}
	if rawResp != nil {
		result.StatusCode = rawResp.StatusCode
		result.BodyText = strings.TrimSpace(string(rawResp.Body))
		if len(rawResp.Body) > 0 {
			var parsed models.LocalSkillTokenResetResponse
			if parseErr := json.Unmarshal(rawResp.Body, &parsed); parseErr == nil {
				result.Parsed = &parsed
			} else {
				result.ParseError = parseErr.Error()
			}
		}
	}
	if err != nil {
		return result, utils.NewAPIError("failed to reset local skill upload token", err)
	}
	return result, nil
}
