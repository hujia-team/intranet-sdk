package models

type LocalSkillUploadResult struct {
	StatusCode int                       `json:"statusCode"`
	BodyText   string                    `json:"bodyText,omitempty"`
	Parsed     *LocalSkillUploadResponse `json:"parsed,omitempty"`
	ParseError string                    `json:"parseError,omitempty"`
}

type LocalSkillUploadResponse struct {
	Code int                          `json:"code"`
	Msg  string                       `json:"msg"`
	Data LocalSkillUploadResponseData `json:"data"`
}

type LocalSkillUploadResponseData struct {
	Created        bool           `json:"created"`
	Skipped        bool           `json:"skipped"`
	UploadToken    *string        `json:"uploadToken,omitempty"`
	PublisherName  *string        `json:"publisherName,omitempty"`
	PublisherEmail *string        `json:"publisherEmail,omitempty"`
	PublisherTeam  *string        `json:"publisherTeam,omitempty"`
	Skill          LocalSkillInfo `json:"skill"`
}

type LocalSkillInfo struct {
	ID          *uint64 `json:"id,omitempty"`
	Slug        *string `json:"slug,omitempty"`
	Name        *string `json:"name,omitempty"`
	Version     *string `json:"version,omitempty"`
	Source      *string `json:"source,omitempty"`
	Path        *string `json:"path,omitempty"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"isActive,omitempty"`
	Status      *uint32 `json:"status,omitempty"`
}

type LocalSkillTokenResetRequest struct {
	Slug string `json:"slug"`
}

type LocalSkillTokenResetResult struct {
	StatusCode int                           `json:"statusCode"`
	BodyText   string                        `json:"bodyText,omitempty"`
	Parsed     *LocalSkillTokenResetResponse `json:"parsed,omitempty"`
	ParseError string                        `json:"parseError,omitempty"`
}

type LocalSkillTokenResetResponse struct {
	Code int                              `json:"code"`
	Msg  string                           `json:"msg"`
	Data LocalSkillTokenResetResponseData `json:"data"`
}

type LocalSkillTokenResetResponseData struct {
	Skill       LocalSkillInfo `json:"skill"`
	UploadToken string         `json:"uploadToken"`
}
