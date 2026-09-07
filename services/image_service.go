// Package services provides business logic for the MiniEye Intranet API.
package services

import (
	"net/url"
	"strconv"

	"github.com/hujia-team/intranet-sdk/client"
	"github.com/hujia-team/intranet-sdk/models"
	"github.com/hujia-team/intranet-sdk/utils"
)

// ImageService defines OCI image metadata operations.
type ImageService interface {
	// ListImages lists OCI images using optional pagination and filters.
	ListImages(req *models.ImageListReq) (*models.ImageListResp, error)
}

type imageService struct {
	httpClient *client.HTTPClient
}

// NewImageService creates a new OCI image service.
func NewImageService(httpClient *client.HTTPClient) ImageService {
	return &imageService{httpClient: httpClient}
}

// ListImages implements ImageService.ListImages.
func (s *imageService) ListImages(req *models.ImageListReq) (*models.ImageListResp, error) {
	var response struct {
		Code int                   `json:"code"`
		Msg  string                `json:"msg"`
		Data *models.ImageListResp `json:"data"`
	}

	utils.Debug("Listing OCI images")
	if err := s.httpClient.GetWithQuery("/api/v1/images", imageListQuery(req), &response); err != nil {
		utils.Error("Failed to list OCI images: %v", err)
		return nil, utils.NewAPIError("failed to list images", err)
	}
	if response.Code != 0 {
		utils.Error("API error: %s", response.Msg)
		return nil, utils.NewAPIError(response.Msg, nil)
	}
	if response.Data == nil {
		return &models.ImageListResp{Items: []models.ImageInfo{}}, nil
	}
	if response.Data.Items == nil {
		response.Data.Items = []models.ImageInfo{}
	}
	return response.Data, nil
}

func imageListQuery(req *models.ImageListReq) url.Values {
	query := url.Values{}
	if req == nil {
		return query
	}
	if req.PageNum != nil {
		query.Set("page_num", strconv.FormatUint(*req.PageNum, 10))
	}
	if req.PageSize != nil {
		query.Set("page_size", strconv.FormatUint(*req.PageSize, 10))
	}
	if req.GitRepo != nil {
		query.Set("git_repo", *req.GitRepo)
	}
	if req.GitBranch != nil {
		query.Set("git_branch", *req.GitBranch)
	}
	if req.ImgFullname != nil {
		query.Set("img_fullname", *req.ImgFullname)
	}
	if req.Digest != nil {
		query.Set("digest", *req.Digest)
	}
	if req.GitCommitShortHash != nil {
		query.Set("git_commit_short_hash", *req.GitCommitShortHash)
	}
	if req.PipelineID != nil {
		query.Set("pipeline_id", strconv.FormatUint(*req.PipelineID, 10))
	}
	if req.JobID != nil {
		query.Set("job_id", strconv.FormatUint(*req.JobID, 10))
	}
	return query
}
