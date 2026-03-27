// Package services provides business logic for the MiniEye Intranet API.
package services

import (
	"encoding/json"
	"fmt"

	"github.com/hujia-team/intranet-sdk/client"
	"github.com/hujia-team/intranet-sdk/models"
	"github.com/hujia-team/intranet-sdk/utils"
)

// ArtifactService defines artifact-management operations.
type ArtifactService interface {
	CreateArtifact(artifact *models.ArtifactInfo) (*models.BaseMsgResp, error)
	UpdateArtifact(artifact *models.ArtifactInfo) (*models.BaseMsgResp, error)
	DeleteArtifacts(ids []uint64) (*models.BaseMsgResp, error)
	ListArtifacts(req *models.ArtifactListReq) (*models.ArtifactListResp, error)
	GetArtifactByID(id uint64) (*models.ArtifactInfo, error)
	GetArtifactByName(name string, lookup *models.ArtifactLookupOptions) (*models.ArtifactInfo, error)
	GetArtifactCommitDiff(artifactIDA, artifactIDB uint64) (*models.ArtifactCommitDiffInfo, error)
	GetArtifactTagSchema(version string) (*models.ArtifactTagSchemaInfo, error)
	GetArtifactTagSchemaJSON(version string) (map[string]any, error)
	GetJfrogToken(projectName string) (*models.JfrogTokenInfo, error)
	GetJfrogTokenByArtifactName(name string, lookup *models.ArtifactLookupOptions) (*models.JfrogTokenInfo, error)
	GetArtifactDownloadURL(artifactID uint64, downloadType string) (*models.ArtifactDownloadURLInfo, error)
	GetArtifactDownloadURLByName(name string, lookup *models.ArtifactLookupOptions, downloadType string) (*models.ArtifactDownloadURLInfo, error)
	GetParsedArtifactTags(artifactID uint64) (map[string]any, error)
	UpdateArtifactTags(artifactID uint64, tags map[string]any, tagSchemaVersion string) (*models.BaseMsgResp, error)
	ParseArtifactTags(tags string, schema any) (map[string]any, error)
}

type artifactService struct {
	httpClient *client.HTTPClient
}

// NewArtifactService creates a new artifact service.
func NewArtifactService(httpClient *client.HTTPClient) ArtifactService {
	return &artifactService{httpClient: httpClient}
}

func (s *artifactService) CreateArtifact(artifact *models.ArtifactInfo) (*models.BaseMsgResp, error) {
	var response models.BaseMsgResp
	utils.Debug("Creating artifact")
	if err := s.httpClient.Post("/aiplorer/artifact/create", artifact, &response); err != nil {
		return nil, utils.NewAPIError("failed to create artifact", err)
	}
	if response.Code != 0 {
		return nil, utils.NewAPIError(response.Msg, nil)
	}
	return &response, nil
}

func (s *artifactService) UpdateArtifact(artifact *models.ArtifactInfo) (*models.BaseMsgResp, error) {
	var response models.BaseMsgResp
	utils.Debug("Updating artifact")
	if err := s.httpClient.Post("/aiplorer/artifact/update", artifact, &response); err != nil {
		return nil, utils.NewAPIError("failed to update artifact", err)
	}
	if response.Code != 0 {
		return nil, utils.NewAPIError(response.Msg, nil)
	}
	return &response, nil
}

func (s *artifactService) DeleteArtifacts(ids []uint64) (*models.BaseMsgResp, error) {
	var response models.BaseMsgResp
	utils.Debug("Deleting artifacts: %v", ids)
	if err := s.httpClient.Post("/aiplorer/artifact/delete", &models.IDsReq{IDs: ids}, &response); err != nil {
		return nil, utils.NewAPIError("failed to delete artifacts", err)
	}
	if response.Code != 0 {
		return nil, utils.NewAPIError(response.Msg, nil)
	}
	return &response, nil
}

func (s *artifactService) ListArtifacts(req *models.ArtifactListReq) (*models.ArtifactListResp, error) {
	var response struct {
		Code int                     `json:"code"`
		Msg  string                  `json:"msg"`
		Data models.ArtifactListResp `json:"data"`
	}
	utils.Debug("Listing artifacts")
	if err := s.httpClient.Post("/aiplorer/artifact/list", req, &response); err != nil {
		return nil, utils.NewAPIError("failed to list artifacts", err)
	}
	if response.Code != 0 {
		return nil, utils.NewAPIError(response.Msg, nil)
	}
	return &response.Data, nil
}

func (s *artifactService) GetArtifactByID(id uint64) (*models.ArtifactInfo, error) {
	var response struct {
		Code int                 `json:"code"`
		Msg  string              `json:"msg"`
		Data models.ArtifactInfo `json:"data"`
	}
	utils.Debug("Getting artifact by ID: %d", id)
	if err := s.httpClient.Post("/aiplorer/artifact", &models.IDReq{ID: id}, &response); err != nil {
		return nil, utils.NewAPIError("failed to get artifact by id", err)
	}
	if response.Code != 0 {
		return nil, utils.NewAPIError(response.Msg, nil)
	}
	return &response.Data, nil
}

func (s *artifactService) GetArtifactByName(name string, lookup *models.ArtifactLookupOptions) (*models.ArtifactInfo, error) {
	req := &models.ArtifactListReq{Page: 1, PageSize: 100, Name: &name}
	if lookup != nil {
		if lookup.ModulePath != "" {
			req.ModulePath = &lookup.ModulePath
		}
		if lookup.ArtifactType != "" {
			req.Type = &lookup.ArtifactType
		}
		if lookup.SemanticVersion != "" {
			req.SemanticVersion = &lookup.SemanticVersion
		}
		req.IsVirtual = lookup.IncludeVirtual
	}
	result, err := s.ListArtifacts(req)
	if err != nil {
		return nil, err
	}
	var matched []models.ArtifactInfo
	for _, item := range result.Data {
		if item.Name != nil && *item.Name == name {
			matched = append(matched, item)
		}
	}
	if len(matched) == 0 {
		return nil, utils.NewAPIError(fmt.Sprintf("artifact not found by name: %s", name), nil)
	}
	if len(matched) > 1 {
		return nil, utils.NewAPIError(fmt.Sprintf("multiple artifacts found by name: %s", name), nil)
	}
	if matched[0].ID == nil {
		return nil, utils.NewAPIError(fmt.Sprintf("artifact id missing for artifact: %s", name), nil)
	}
	return s.GetArtifactByID(*matched[0].ID)
}

func (s *artifactService) GetArtifactCommitDiff(artifactIDA, artifactIDB uint64) (*models.ArtifactCommitDiffInfo, error) {
	var response struct {
		Code int                           `json:"code"`
		Msg  string                        `json:"msg"`
		Data models.ArtifactCommitDiffInfo `json:"data"`
	}
	if err := s.httpClient.Post("/aiplorer/artifact/commit-diff", &models.ArtifactCommitDiffReq{
		ArtifactIDA: artifactIDA,
		ArtifactIDB: artifactIDB,
	}, &response); err != nil {
		return nil, utils.NewAPIError("failed to get artifact commit diff", err)
	}
	if response.Code != 0 {
		return nil, utils.NewAPIError(response.Msg, nil)
	}
	return &response.Data, nil
}

func (s *artifactService) GetArtifactTagSchema(version string) (*models.ArtifactTagSchemaInfo, error) {
	var response struct {
		Code int                          `json:"code"`
		Msg  string                       `json:"msg"`
		Data models.ArtifactTagSchemaInfo `json:"data"`
	}
	req := &models.ArtifactTagSchemaReq{}
	if version != "" {
		req.Version = version
	}
	if err := s.httpClient.Post("/aiplorer/artifact/tag-schema", req, &response); err != nil {
		return nil, utils.NewAPIError("failed to get artifact tag schema", err)
	}
	if response.Code != 0 {
		return nil, utils.NewAPIError(response.Msg, nil)
	}
	return &response.Data, nil
}

func (s *artifactService) GetArtifactTagSchemaJSON(version string) (map[string]any, error) {
	schema, err := s.GetArtifactTagSchema(version)
	if err != nil {
		return nil, err
	}
	return models.ParseJSON(schema.Content)
}

func (s *artifactService) GetJfrogToken(projectName string) (*models.JfrogTokenInfo, error) {
	var response struct {
		Code int                   `json:"code"`
		Msg  string                `json:"msg"`
		Data models.JfrogTokenInfo `json:"data"`
	}
	if err := s.httpClient.Post("/aiplorer/jfrog/token", &models.JfrogTokenReq{ProjectName: projectName}, &response); err != nil {
		return nil, utils.NewAPIError("failed to get jfrog token", err)
	}
	if response.Code != 0 {
		return nil, utils.NewAPIError(response.Msg, nil)
	}
	return &response.Data, nil
}

func (s *artifactService) GetJfrogTokenByArtifactName(name string, lookup *models.ArtifactLookupOptions) (*models.JfrogTokenInfo, error) {
	artifact, err := s.GetArtifactByName(name, lookup)
	if err != nil {
		return nil, err
	}
	if artifact.ProjectName == nil || *artifact.ProjectName == "" {
		return nil, utils.NewAPIError(fmt.Sprintf("artifact project_name is empty: %s", name), nil)
	}
	return s.GetJfrogToken(*artifact.ProjectName)
}

func (s *artifactService) GetArtifactDownloadURL(artifactID uint64, downloadType string) (*models.ArtifactDownloadURLInfo, error) {
	if downloadType == "" {
		downloadType = "artifact"
	}
	var response struct {
		Code int                            `json:"code"`
		Msg  string                         `json:"msg"`
		Data models.ArtifactDownloadURLInfo `json:"data"`
	}
	if err := s.httpClient.Post("/aiplorer/artifact/download-url", &models.ArtifactDownloadURLReq{
		ArtifactID:   artifactID,
		DownloadType: downloadType,
	}, &response); err != nil {
		return nil, utils.NewAPIError("failed to get artifact download url", err)
	}
	if response.Code != 0 {
		return nil, utils.NewAPIError(response.Msg, nil)
	}
	return &response.Data, nil
}

func (s *artifactService) GetArtifactDownloadURLByName(name string, lookup *models.ArtifactLookupOptions, downloadType string) (*models.ArtifactDownloadURLInfo, error) {
	artifact, err := s.GetArtifactByName(name, lookup)
	if err != nil {
		return nil, err
	}
	if artifact.ID == nil {
		return nil, utils.NewAPIError(fmt.Sprintf("artifact id missing for artifact: %s", name), nil)
	}
	return s.GetArtifactDownloadURL(*artifact.ID, downloadType)
}

func (s *artifactService) GetParsedArtifactTags(artifactID uint64) (map[string]any, error) {
	artifact, err := s.GetArtifactByID(artifactID)
	if err != nil {
		return nil, err
	}
	if artifact.Tags == nil || *artifact.Tags == "" {
		return map[string]any{}, nil
	}
	version := ""
	if artifact.TagSchemaVersion != nil {
		version = *artifact.TagSchemaVersion
	}
	schema, err := s.GetArtifactTagSchema(version)
	if err != nil {
		return nil, err
	}
	return s.ParseArtifactTags(*artifact.Tags, schema)
}

func (s *artifactService) UpdateArtifactTags(artifactID uint64, tags map[string]any, tagSchemaVersion string) (*models.BaseMsgResp, error) {
	if tagSchemaVersion == "" {
		if rawVersion, ok := tags["schema_version"].(string); ok {
			tagSchemaVersion = rawVersion
		}
	}
	if tagSchemaVersion == "" {
		artifact, err := s.GetArtifactByID(artifactID)
		if err != nil {
			return nil, err
		}
		if artifact.TagSchemaVersion != nil {
			tagSchemaVersion = *artifact.TagSchemaVersion
		}
	}
	schema, err := s.GetArtifactTagSchema(tagSchemaVersion)
	if err != nil {
		return nil, err
	}
	if _, err := s.ParseArtifactTags(mustJSON(tags), schema); err != nil {
		return nil, err
	}
	return s.UpdateArtifact(&models.ArtifactInfo{
		ID:               &artifactID,
		Tags:             stringPtr(mustJSON(tags)),
		TagSchemaVersion: stringPtr(tagSchemaVersion),
	})
}

func (s *artifactService) ParseArtifactTags(tags string, schema any) (map[string]any, error) {
	parsedTags, err := models.ParseJSON(tags)
	if err != nil {
		return nil, utils.NewAPIError("failed to decode artifact tags", err)
	}
	var parsedSchema map[string]any
	switch v := schema.(type) {
	case *models.ArtifactTagSchemaInfo:
		parsedSchema, err = models.ParseJSON(v.Content)
	case models.ArtifactTagSchemaInfo:
		parsedSchema, err = models.ParseJSON(v.Content)
	case string:
		parsedSchema, err = models.ParseJSON(v)
	case []byte:
		parsedSchema, err = models.ParseJSON(string(v))
	case map[string]any:
		parsedSchema = v
	default:
		return nil, utils.NewAPIError("unsupported artifact tag schema type", nil)
	}
	if err != nil {
		return nil, utils.NewAPIError("failed to decode artifact tag schema", err)
	}
	schemaVersion, hasSchemaVersion := parsedSchema["version"].(string)
	tagSchemaVersion, hasTagSchemaVersion := parsedTags["schema_version"].(string)
	if hasSchemaVersion && hasTagSchemaVersion && schemaVersion != tagSchemaVersion {
		return nil, utils.NewAPIError("artifact tag schema version does not match tag schema_version", nil)
	}
	return parsedTags, nil
}

func mustJSON(v any) string {
	buf, _ := json.Marshal(v)
	return string(buf)
}

func stringPtr(v string) *string {
	return &v
}
