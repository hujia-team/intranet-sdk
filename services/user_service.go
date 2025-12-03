// Package services provides business logic for the MiniEye Intranet API.
package services

import (
	"github.com/hujia-team/intranet-sdk/client"
	"github.com/hujia-team/intranet-sdk/models"
	"github.com/hujia-team/intranet-sdk/utils"
)

// UserService defines the user service interface.
type UserService interface {
	// GetUserInfo gets the current user's information.
	GetUserInfo() (*models.UserInfo, error)
}

// userService implements the UserService interface.
type userService struct {
	httpClient *client.HTTPClient
}

// NewUserService creates a new user service.
func NewUserService(httpClient *client.HTTPClient) UserService {
	return &userService{
		httpClient: httpClient,
	}
}

// GetUserInfo implements the UserService.GetUserInfo method.
func (s *userService) GetUserInfo() (*models.UserInfo, error) {
	// 使用嵌套结构体直接解析响应
	var response struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data models.UserInfo `json:"data"`
	}

	utils.Debug("Getting current user info")
	err := s.httpClient.Get("/user/info", &response)
	if err != nil {
		utils.Error("Failed to get user info: %v", err)
		return nil, utils.NewAPIError("failed to get user info", err)
	}

	if response.Code != 0 {
		utils.Error("API error: %s", response.Msg)
		return nil, utils.NewAPIError(response.Msg, nil)
	}
	utils.Debug("Got user info successfully")
	return &response.Data, nil
}
