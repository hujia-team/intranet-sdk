package models

type MultiRepoMergeSetItemInfo struct {
	ID         *uint64 `json:"id,optional"`
	CreatedAt  *int64  `json:"createdAt,optional"`
	UpdatedAt  *int64  `json:"updatedAt,optional"`
	MergeSetID *uint64 `json:"mergeSetId,optional"`
	URL        *string `json:"url,optional"`
	CommitID   *string `json:"commitId,optional"`
}

type MultiRepoMergeSetInfo struct {
	ID        *uint64                     `json:"id,optional"`
	CreatedAt *int64                      `json:"createdAt,optional"`
	UpdatedAt *int64                      `json:"updatedAt,optional"`
	Name      *string                     `json:"name,optional"`
	Project   *string                     `json:"project,optional"`
	Branch    *string                     `json:"branch,optional"`
	Items     []MultiRepoMergeSetItemInfo `json:"items,optional"`
}

type CreateMultiRepoMergeSetReq struct {
	Name    string `json:"name"`
	Project string `json:"project"`
	Branch  string `json:"branch"`
}

type AddMultiRepoMergeSetItemReq struct {
	MergeSetID uint64 `json:"mergeSetId"`
	URL        string `json:"url"`
	CommitID   string `json:"commitId"`
}

type RemoveMultiRepoMergeSetItemReq struct {
	MergeSetID uint64 `json:"mergeSetId"`
	URL        string `json:"url"`
}

type MultiRepoMergeSetListReq struct {
	Project string `json:"project,optional"`
	Branch  string `json:"branch,optional"`
}

type MultiRepoMergeSetListResp struct {
	Total uint64                  `json:"total"`
	Data  []MultiRepoMergeSetInfo `json:"data"`
}

type UpsertMultiRepoMergeSetPipelineReq struct {
	Project      string `json:"project"`
	Branch       string `json:"branch"`
	PipelineID   string `json:"pipelineId"`
	SourceURL    string `json:"sourceUrl"`
	SourceCommit string `json:"sourceCommit"`
}
