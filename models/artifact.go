// Package models defines the data structures used in the MiniEye Intranet API.
package models

import "encoding/json"

// CommitInfo describes a git commit associated with an artifact.
type CommitInfo struct {
	ID             *uint64 `json:"id,omitempty"`
	CreatedAt      *int64  `json:"createdAt,omitempty"`
	UpdatedAt      *int64  `json:"updatedAt,omitempty"`
	RepositoryID   *uint64 `json:"repositoryId,omitempty"`
	RepositoryName *string `json:"repositoryName,omitempty"`
	RepositoryPath *string `json:"repositoryPath,omitempty"`
	CommitHash     *string `json:"commitHash,omitempty"`
	ShortHash      *string `json:"shortHash,omitempty"`
	Branch         *string `json:"branch,omitempty"`
	Author         *string `json:"author,omitempty"`
	AuthorEmail    *string `json:"authorEmail,omitempty"`
	Message        *string `json:"message,omitempty"`
	CommittedAt    *int64  `json:"committedAt,omitempty"`
	CommitTitle    *string `json:"commitTitle,omitempty"`
}

// ArtifactDependencyInfo describes a dependency or dependent artifact.
type ArtifactDependencyInfo struct {
	ID          *uint64      `json:"id,omitempty"`
	CreatedAt   *int64       `json:"createdAt,omitempty"`
	UpdatedAt   *int64       `json:"updatedAt,omitempty"`
	Name        *string      `json:"name,omitempty"`
	Type        *string      `json:"type,omitempty"`
	IsVirtual   *bool        `json:"isVirtual,omitempty"`
	CommitHash  *string      `json:"commitHash,omitempty"`
	ModulePath  *string      `json:"modulePath,omitempty"`
	PipelineURL *string      `json:"pipelineUrl,omitempty"`
	ParentID    *uint64      `json:"parentId,omitempty"`
	Commits     []CommitInfo `json:"commits,omitempty"`
}

// ArtifactInfo describes artifact metadata and lineage information.
type ArtifactInfo struct {
	ID               *uint64                  `json:"id,omitempty"`
	CreatedAt        *int64                   `json:"createdAt,omitempty"`
	UpdatedAt        *int64                   `json:"updatedAt,omitempty"`
	FullPath         *string                  `json:"fullPath,omitempty"`
	MetadataPath     *string                  `json:"metadataPath,omitempty"`
	FileHash         *string                  `json:"fileHash,omitempty"`
	CommitHash       *string                  `json:"commitHash,omitempty"`
	Name             *string                  `json:"name,omitempty"`
	Type             *string                  `json:"type,omitempty"`
	ProjectName      *string                  `json:"projectName,omitempty"`
	IsVirtual        *bool                    `json:"isVirtual,omitempty"`
	ModulePath       *string                  `json:"modulePath,omitempty"`
	SemanticVersion  *string                  `json:"semanticVersion,omitempty"`
	Tags             *string                  `json:"tags,omitempty"`
	TagSchemaVersion *string                  `json:"tagSchemaVersion,omitempty"`
	Extra            *string                  `json:"extra,omitempty"`
	PipelineID       *string                  `json:"pipelineId,omitempty"`
	PipelineURL      *string                  `json:"pipelineUrl,omitempty"`
	BuildDate        *int64                   `json:"buildDate,omitempty"`
	Commits          []CommitInfo             `json:"commits,omitempty"`
	Dependencies     []ArtifactDependencyInfo `json:"dependencies,omitempty"`
	Dependents       []ArtifactDependencyInfo `json:"dependents,omitempty"`
}

// ArtifactListReq lists artifacts with filters and pagination.
type ArtifactListReq struct {
	Page             uint64  `json:"page"`
	PageSize         uint64  `json:"pageSize"`
	CreatedAt        *int64  `json:"createdAt,omitempty"`
	UpdatedAt        *int64  `json:"updatedAt,omitempty"`
	FullPath         *string `json:"fullPath,omitempty"`
	FileHash         *string `json:"fileHash,omitempty"`
	CommitHash       *string `json:"commitHash,omitempty"`
	Name             *string `json:"name,omitempty"`
	ProjectName      *string `json:"projectName,omitempty"`
	Type             *string `json:"type,omitempty"`
	IsVirtual        *bool   `json:"isVirtual,omitempty"`
	ModulePath       *string `json:"modulePath,omitempty"`
	SemanticVersion  *string `json:"semanticVersion,omitempty"`
	Tags             *string `json:"tags,omitempty"`
	TagSchemaVersion *string `json:"tagSchemaVersion,omitempty"`
	Extra            *string `json:"extra,omitempty"`
	PipelineID       *string `json:"pipelineId,omitempty"`
	PipelineURL      *string `json:"pipelineUrl,omitempty"`
	BuildDate        *int64  `json:"buildDate,omitempty"`
	ExactCommitHash  *bool   `json:"exactCommitHash,omitempty"`
}

// ArtifactListResp is the paged artifact list response body.
type ArtifactListResp struct {
	Total uint64         `json:"total"`
	Data  []ArtifactInfo `json:"data"`
}

// ArtifactLookupOptions holds reusable filters for name-based lookups.
type ArtifactLookupOptions struct {
	ModulePath      string
	ArtifactType    string
	SemanticVersion string
	IncludeVirtual  *bool
	ProjectName     string
}

// ArtifactChildHashInfo describes one dependency node extracted from artifact lineage.
type ArtifactChildHashInfo struct {
	ID         *uint64 `json:"id,omitempty"`
	ParentID   *uint64 `json:"parentId,omitempty"`
	Name       *string `json:"name,omitempty"`
	Type       *string `json:"type,omitempty"`
	CommitHash *string `json:"commitHash,omitempty"`
	ModulePath *string `json:"modulePath,omitempty"`
}

// ArtifactChildHashesInfo is a helper view built from artifact detail responses.
type ArtifactChildHashesInfo struct {
	RootArtifactID   *uint64                 `json:"rootArtifactId,omitempty"`
	RootArtifactName *string                 `json:"rootArtifactName,omitempty"`
	RootArtifactType *string                 `json:"rootArtifactType,omitempty"`
	RootCommitHash   *string                 `json:"rootCommitHash,omitempty"`
	ChildHashes      []ArtifactChildHashInfo `json:"childHashes,omitempty"`
}

// ArtifactTagSchemaInfo contains schema metadata.
type ArtifactTagSchemaInfo struct {
	Version string `json:"version"`
	Content string `json:"content"`
}

// JfrogTokenInfo contains a JFrog access token response.
type JfrogTokenInfo struct {
	TokenID     string `json:"token_id"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	URL         string `json:"url"`
}

// ArtifactDownloadURLInfo contains a signed download URL.
type ArtifactDownloadURLInfo struct {
	DownloadURL string `json:"downloadUrl"`
	ExpireTime  string `json:"expireTime"`
	FileName    string `json:"fileName"`
	FilePath    string `json:"filePath"`
}

// ArtifactExistenceInfo describes whether an artifact exists for a lookup target.
type ArtifactExistenceInfo struct {
	CommitHash string  `json:"commitHash"`
	Exists     bool    `json:"exists"`
	ArtifactID *uint64 `json:"artifactId,omitempty"`
	Name       *string `json:"name,omitempty"`
}

// BatchArtifactExistenceResp is the batch existence check response body.
type BatchArtifactExistenceResp struct {
	Data []ArtifactExistenceInfo `json:"data"`
}

// GetArtifactByCommitHashReq looks up an artifact by exact commit hash.
type GetArtifactByCommitHashReq struct {
	CommitHash      string  `json:"commitHash"`
	ModulePath      *string `json:"modulePath,omitempty"`
	ArtifactType    *string `json:"type,omitempty"`
	SemanticVersion *string `json:"semanticVersion,omitempty"`
	IsVirtual       *bool   `json:"isVirtual,omitempty"`
	ProjectName     *string `json:"projectName,omitempty"`
}

// BatchCheckArtifactsExistReq performs a batch existence check by exact commit hashes.
type BatchCheckArtifactsExistReq struct {
	CommitHashes    []string `json:"commitHashes"`
	ModulePath      *string  `json:"modulePath,omitempty"`
	ArtifactType    *string  `json:"type,omitempty"`
	SemanticVersion *string  `json:"semanticVersion,omitempty"`
	IsVirtual       *bool    `json:"isVirtual,omitempty"`
	ProjectName     *string  `json:"projectName,omitempty"`
}

// ArtifactVersionMetadataInfo contains version metadata fetched by commit hash.
type ArtifactVersionMetadataInfo struct {
	ArtifactID       *uint64        `json:"artifactId,omitempty"`
	CommitHash       *string        `json:"commitHash,omitempty"`
	Name             *string        `json:"name,omitempty"`
	MetadataPath     *string        `json:"metadataPath,omitempty"`
	MetadataFileName *string        `json:"metadataFileName,omitempty"`
	RawContent       *string        `json:"rawContent,omitempty"`
	Parsed           map[string]any `json:"-"`
}

// ArtifactDownloadPlan describes a resolved download plan for one artifact.
type ArtifactDownloadPlan struct {
	Artifact        *ArtifactInfo            `json:"artifact,omitempty"`
	Token           *JfrogTokenInfo          `json:"token,omitempty"`
	DownloadURL     *ArtifactDownloadURLInfo `json:"downloadUrl,omitempty"`
	TargetPath      string                   `json:"targetPath"`
	Checksum        string                   `json:"checksum,omitempty"`
	SkippedExisting bool                     `json:"skippedExisting,omitempty"`
}

// RepoDiff groups artifact commit differences by repository.
type RepoDiff struct {
	RepositoryID   uint64       `json:"repositoryId"`
	RepositoryName *string      `json:"repositoryName,omitempty"`
	RepositoryPath *string      `json:"repositoryPath,omitempty"`
	OlderCommit    *CommitInfo  `json:"olderCommit,omitempty"`
	NewerCommit    *CommitInfo  `json:"newerCommit,omitempty"`
	Commits        []CommitInfo `json:"commits,omitempty"`
}

// ArtifactCommitDiffInfo describes the commit diff between two artifacts.
type ArtifactCommitDiffInfo struct {
	OlderArtifactID  uint64     `json:"olderArtifactId"`
	NewerArtifactID  uint64     `json:"newerArtifactId"`
	ChangedRepoCount uint64     `json:"changedRepoCount"`
	RepoDiffs        []RepoDiff `json:"repoDiffs"`
}

// IDsReq is a numeric ID list request.
type IDsReq struct {
	IDs []uint64 `json:"ids"`
}

// IDReq is a numeric ID request.
type IDReq struct {
	ID uint64 `json:"id"`
}

// ArtifactCommitDiffReq compares two artifacts.
type ArtifactCommitDiffReq struct {
	ArtifactIDA uint64 `json:"artifactIdA"`
	ArtifactIDB uint64 `json:"artifactIdB"`
}

// ArtifactTagSchemaReq queries a tag schema by version.
type ArtifactTagSchemaReq struct {
	Version string `json:"version,omitempty"`
}

// JfrogTokenReq requests a JFrog token.
type JfrogTokenReq struct {
	ProjectName string `json:"projectName"`
}

// ArtifactDownloadURLReq requests an artifact download URL.
type ArtifactDownloadURLReq struct {
	ArtifactID   uint64 `json:"artifactId"`
	DownloadType string `json:"downloadType"`
}

// ArtifactTagUpdateReq updates only artifact tags.
type ArtifactTagUpdateReq struct {
	ID               uint64 `json:"id"`
	Tags             string `json:"tags"`
	TagSchemaVersion string `json:"tagSchemaVersion,omitempty"`
}

// ParseJSON parses a raw JSON string to a generic object.
func ParseJSON(raw string) (map[string]any, error) {
	if raw == "" {
		return map[string]any{}, nil
	}
	var parsed map[string]any
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return nil, err
	}
	return parsed, nil
}
