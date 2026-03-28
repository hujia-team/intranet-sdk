package services

import (
	"github.com/hujia-team/intranet-sdk/client"
	"github.com/hujia-team/intranet-sdk/models"
	"github.com/hujia-team/intranet-sdk/utils"
)

type MultiRepoMergeSetService interface {
	Create(req *models.CreateMultiRepoMergeSetReq) (uint64, error)
	List(req *models.MultiRepoMergeSetListReq) (*models.MultiRepoMergeSetListResp, error)
	Get(id uint64) (*models.MultiRepoMergeSetInfo, error)
	AddItem(req *models.AddMultiRepoMergeSetItemReq) error
	RemoveItem(req *models.RemoveMultiRepoMergeSetItemReq) error
	Delete(id uint64) error
	UpsertPipeline(req *models.UpsertMultiRepoMergeSetPipelineReq) error
}

type multiRepoMergeSetService struct {
	httpClient *client.HTTPClient
}

func NewMultiRepoMergeSetService(httpClient *client.HTTPClient) MultiRepoMergeSetService {
	return &multiRepoMergeSetService{httpClient: httpClient}
}

func (s *multiRepoMergeSetService) Create(req *models.CreateMultiRepoMergeSetReq) (uint64, error) {
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			ID uint64 `json:"id"`
		} `json:"data"`
	}

	if err := s.httpClient.Post("/aiplorer/multi-repo-merge-set/create", req, &response); err != nil {
		return 0, utils.NewAPIError("failed to create multi repo merge set", err)
	}
	if response.Code != 0 {
		return 0, utils.NewAPIError(response.Msg, nil)
	}
	return response.Data.ID, nil
}

func (s *multiRepoMergeSetService) List(req *models.MultiRepoMergeSetListReq) (*models.MultiRepoMergeSetListResp, error) {
	var response struct {
		Code int                              `json:"code"`
		Msg  string                           `json:"msg"`
		Data models.MultiRepoMergeSetListResp `json:"data"`
	}

	if err := s.httpClient.Post("/aiplorer/multi-repo-merge-set/list", req, &response); err != nil {
		return nil, utils.NewAPIError("failed to list multi repo merge sets", err)
	}
	if response.Code != 0 {
		return nil, utils.NewAPIError(response.Msg, nil)
	}
	return &response.Data, nil
}

func (s *multiRepoMergeSetService) Get(id uint64) (*models.MultiRepoMergeSetInfo, error) {
	var response struct {
		Code int                          `json:"code"`
		Msg  string                       `json:"msg"`
		Data models.MultiRepoMergeSetInfo `json:"data"`
	}

	if err := s.httpClient.Post("/aiplorer/multi-repo-merge-set", &models.IDReq{ID: id}, &response); err != nil {
		return nil, utils.NewAPIError("failed to get multi repo merge set", err)
	}
	if response.Code != 0 {
		return nil, utils.NewAPIError(response.Msg, nil)
	}
	return &response.Data, nil
}

func (s *multiRepoMergeSetService) AddItem(req *models.AddMultiRepoMergeSetItemReq) error {
	return s.postBase("/aiplorer/multi-repo-merge-set/add-item", req, "failed to add multi repo merge set item")
}

func (s *multiRepoMergeSetService) RemoveItem(req *models.RemoveMultiRepoMergeSetItemReq) error {
	return s.postBase("/aiplorer/multi-repo-merge-set/remove-item", req, "failed to remove multi repo merge set item")
}

func (s *multiRepoMergeSetService) Delete(id uint64) error {
	return s.postBase("/aiplorer/multi-repo-merge-set/delete", &models.IDReq{ID: id}, "failed to delete multi repo merge set")
}

func (s *multiRepoMergeSetService) UpsertPipeline(req *models.UpsertMultiRepoMergeSetPipelineReq) error {
	return s.postBase("/aiplorer/multi-repo-merge-set/pipeline/upsert", req, "failed to upsert multi repo merge set pipeline")
}

func (s *multiRepoMergeSetService) postBase(endpoint string, req interface{}, message string) error {
	var response models.BaseMsgResp
	if err := s.httpClient.Post(endpoint, req, &response); err != nil {
		return utils.NewAPIError(message, err)
	}
	if response.Code != 0 {
		return utils.NewAPIError(response.Msg, nil)
	}
	return nil
}
