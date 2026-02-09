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

	// ListUsers lists users with pagination and filters.
	ListUsers(req *models.UserListReq) (*models.UserListRsp, error)

	// GetUserById gets user information by UUID.
	GetUserById(uuid string) (*models.UserInfo, error)
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

// ListUsers implements the UserService.ListUsers method.
func (s *userService) ListUsers(req *models.UserListReq) (*models.UserListRsp, error) {
	// 使用嵌套结构体直接解析响应
	var response struct {
		Code int                 `json:"code"`
		Msg  string              `json:"msg"`
		Data models.UserListRsp  `json:"data"`
	}

	utils.Debug("Listing users with page=%d, pageSize=%d", req.Page, req.PageSize)
	err := s.httpClient.Post("/user/list", req, &response)
	if err != nil {
		utils.Error("Failed to list users: %v", err)
		return nil, utils.NewAPIError("failed to list users", err)
	}

	if response.Code != 0 {
		utils.Error("API error: %s", response.Msg)
		return nil, utils.NewAPIError(response.Msg, nil)
	}
	utils.Debug("Listed %d users successfully (total: %d)", len(response.Data.Data), response.Data.Total)
	return &response.Data, nil
}

// GetUserById implements the UserService.GetUserById method.
func (s *userService) GetUserById(uuid string) (*models.UserInfo, error) {
	// 构造请求参数
	req := models.UUIDReq{
		Id: uuid,
	}

	// 使用嵌套结构体直接解析响应
	var response struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data models.UserInfo `json:"data"`
	}

	utils.Debug("Getting user info by UUID: %s", uuid)
	err := s.httpClient.Post("/user", req, &response)
	if err != nil {
		utils.Error("Failed to get user by id: %v", err)
		return nil, utils.NewAPIError("failed to get user by id", err)
	}

	if response.Code != 0 {
		utils.Error("API error: %s", response.Msg)
		return nil, utils.NewAPIError(response.Msg, nil)
	}
	utils.Debug("Got user info successfully for UUID: %s", uuid)
	return &response.Data, nil
}
