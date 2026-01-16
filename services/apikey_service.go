// Package services provides business logic for the MiniEye Intranet API.
package services

import (
	"github.com/hujia-team/intranet-sdk/client"
	"github.com/hujia-team/intranet-sdk/models"
	"github.com/hujia-team/intranet-sdk/utils"
)

// ApiKeyService defines the apikey service interface.
type ApiKeyService interface {
	// CreateApiKey creates a new API key.
	CreateApiKey(apiKey *models.ApiKeyInfo) (uint64, error)

	// UpdateApiKey updates an existing API key.
	UpdateApiKey(apiKey *models.ApiKeyInfo) error

	// DeleteApiKey deletes API keys by IDs.
	DeleteApiKey(ids []uint64) error

	// GetApiKeyList gets the API key list.
	GetApiKeyList(req *models.ApiKeyListReq) (*models.ApiKeyListResp, error)

	// GetApiKeyByID gets an API key by ID.
	GetApiKeyByID(id uint64) (*models.ApiKeyInfo, error)

	// GetSub2ApiKey gets the sub2api API key for current user.
	GetSub2ApiKey() (*models.ApiKeyInfo, error)

	// GetAvailableGroups gets available subscription groups for current user.
	GetAvailableGroups() (*models.GetAvailableGroupsResp, error)

	// GetCurrentGroup gets the current subscription group bound to the user's API key.
	GetCurrentGroup() (*models.CurrentGroupResp, error)

	// SwitchGroup switches the subscription group for an API key.
	SwitchGroup(req *models.SwitchGroupReq) (*models.CurrentGroupResp, error)
}

// apiKeyService implements the ApiKeyService interface.
type apiKeyService struct {
	httpClient *client.HTTPClient
}

// NewApiKeyService creates a new apikey service.
func NewApiKeyService(httpClient *client.HTTPClient) ApiKeyService {
	return &apiKeyService{
		httpClient: httpClient,
	}
}

// CreateApiKey implements the ApiKeyService.CreateApiKey method.
func (s *apiKeyService) CreateApiKey(apiKey *models.ApiKeyInfo) (uint64, error) {
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			ID uint64 `json:"id"`
		} `json:"data"`
	}

	utils.Debug("Creating API key: %s", apiKey.Name)
	err := s.httpClient.Post("/aiplorer/api_key/create", apiKey, &response)
	if err != nil {
		utils.Error("Failed to create API key: %v", err)
		return 0, utils.NewAPIError("failed to create API key", err)
	}

	if response.Code != 0 {
		utils.Error("API error: %s", response.Msg)
		return 0, utils.NewAPIError(response.Msg, nil)
	}

	utils.Debug("Created API key successfully, ID: %d", response.Data.ID)
	return response.Data.ID, nil
}

// UpdateApiKey implements the ApiKeyService.UpdateApiKey method.
func (s *apiKeyService) UpdateApiKey(apiKey *models.ApiKeyInfo) error {
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	utils.Debug("Updating API key ID: %d", apiKey.ID)
	err := s.httpClient.Post("/aiplorer/api_key/update", apiKey, &response)
	if err != nil {
		utils.Error("Failed to update API key: %v", err)
		return utils.NewAPIError("failed to update API key", err)
	}

	if response.Code != 0 {
		utils.Error("API error: %s", response.Msg)
		return utils.NewAPIError(response.Msg, nil)
	}

	utils.Debug("Updated API key successfully")
	return nil
}

// DeleteApiKey implements the ApiKeyService.DeleteApiKey method.
func (s *apiKeyService) DeleteApiKey(ids []uint64) error {
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	reqBody := struct {
		IDs []uint64 `json:"ids"`
	}{
		IDs: ids,
	}

	utils.Debug("Deleting API keys: %v", ids)
	err := s.httpClient.Post("/aiplorer/api_key/delete", reqBody, &response)
	if err != nil {
		utils.Error("Failed to delete API keys: %v", err)
		return utils.NewAPIError("failed to delete API keys", err)
	}

	if response.Code != 0 {
		utils.Error("API error: %s", response.Msg)
		return utils.NewAPIError(response.Msg, nil)
	}

	utils.Debug("Deleted API keys successfully")
	return nil
}

// GetApiKeyList implements the ApiKeyService.GetApiKeyList method.
func (s *apiKeyService) GetApiKeyList(req *models.ApiKeyListReq) (*models.ApiKeyListResp, error) {
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Total uint64              `json:"total"`
			Data  []models.ApiKeyInfo `json:"data"`
		} `json:"data"`
	}

	utils.Debug("Getting API key list, page: %d, page_size: %d", req.Page, req.PageSize)
	err := s.httpClient.Post("/aiplorer/api_key/list", req, &response)
	if err != nil {
		utils.Error("Failed to get API key list: %v", err)
		return nil, utils.NewAPIError("failed to get API key list", err)
	}

	if response.Code != 0 {
		utils.Error("API error: %s", response.Msg)
		return nil, utils.NewAPIError(response.Msg, nil)
	}

	result := &models.ApiKeyListResp{
		Total: response.Data.Total,
		List:  response.Data.Data,
	}

	utils.Debug("Got API key list successfully, total: %d", result.Total)
	return result, nil
}

// GetApiKeyByID implements the ApiKeyService.GetApiKeyByID method.
func (s *apiKeyService) GetApiKeyByID(id uint64) (*models.ApiKeyInfo, error) {
	var response struct {
		Code int               `json:"code"`
		Msg  string            `json:"msg"`
		Data models.ApiKeyInfo `json:"data"`
	}

	reqBody := struct {
		ID uint64 `json:"id"`
	}{
		ID: id,
	}

	utils.Debug("Getting API key by ID: %d", id)
	err := s.httpClient.Post("/aiplorer/api_key", reqBody, &response)
	if err != nil {
		utils.Error("Failed to get API key: %v", err)
		return nil, utils.NewAPIError("failed to get API key", err)
	}

	if response.Code != 0 {
		utils.Error("API error: %s", response.Msg)
		return nil, utils.NewAPIError(response.Msg, nil)
	}

	utils.Debug("Got API key successfully")
	return &response.Data, nil
}

// GetSub2ApiKey implements the ApiKeyService.GetSub2ApiKey method.
func (s *apiKeyService) GetSub2ApiKey() (*models.ApiKeyInfo, error) {
	var response struct {
		Code int               `json:"code"`
		Msg  string            `json:"msg"`
		Data models.ApiKeyInfo `json:"data"`
	}

	utils.Debug("Getting sub2api API key for current user")
	err := s.httpClient.Post("/aiplorer/sub2api/api_key", nil, &response)
	if err != nil {
		utils.Error("Failed to get sub2api API key: %v", err)
		return nil, utils.NewAPIError("failed to get sub2api API key", err)
	}

	if response.Code != 0 {
		utils.Error("API error: %s", response.Msg)
		return nil, utils.NewAPIError(response.Msg, nil)
	}

	utils.Debug("Got sub2api API key successfully")
	return &response.Data, nil
}

// GetAvailableGroups implements the ApiKeyService.GetAvailableGroups method.
func (s *apiKeyService) GetAvailableGroups() (*models.GetAvailableGroupsResp, error) {
	var response struct {
		Code int                       `json:"code"`
		Msg  string                    `json:"msg"`
		Data []models.Sub2ApiGroupInfo `json:"data"`
	}

	utils.Debug("Getting available subscription groups")
	err := s.httpClient.Post("/aiplorer/sub2api/group/available", nil, &response)
	if err != nil {
		utils.Error("Failed to get available groups: %v", err)
		return nil, utils.NewAPIError("failed to get available groups", err)
	}

	if response.Code != 0 {
		utils.Error("API error: %s", response.Msg)
		return nil, utils.NewAPIError(response.Msg, nil)
	}

	utils.Debug("Got available groups successfully, count: %d", len(response.Data))
	return &models.GetAvailableGroupsResp{
		Data: response.Data,
	}, nil
}

// GetCurrentGroup implements the ApiKeyService.GetCurrentGroup method.
func (s *apiKeyService) GetCurrentGroup() (*models.CurrentGroupResp, error) {
	var response models.CurrentGroupResp

	utils.Debug("Getting current group for user's API key")
	err := s.httpClient.Post("/aiplorer/sub2api/group/current", nil, &response)
	if err != nil {
		utils.Error("Failed to get current group: %v", err)
		return nil, utils.NewAPIError("failed to get current group", err)
	}

	if response.Code != 0 {
		utils.Error("API error: %s", response.Msg)
		return nil, utils.NewAPIError(response.Msg, nil)
	}

	utils.Debug("Got current group successfully")
	return &response, nil
}

// SwitchGroup implements the ApiKeyService.SwitchGroup method.
func (s *apiKeyService) SwitchGroup(req *models.SwitchGroupReq) (*models.CurrentGroupResp, error) {
	var response models.CurrentGroupResp
	utils.Debug("Switching group to group ID: %d", req.GroupID)
	err := s.httpClient.Post("/aiplorer/sub2api/group/switch", req, &response)
	if err != nil {
		utils.Error("Failed to switch group: %v", err)
		return nil, utils.NewAPIError("failed to switch group", err)
	}

	if response.Code != 0 {
		utils.Error("API error: %s", response.Msg)
		return nil, utils.NewAPIError(response.Msg, nil)
	}

	utils.Debug("Switched group successfully")
	return &response, nil
}
