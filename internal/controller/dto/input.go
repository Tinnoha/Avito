package dto

type UserActiveDTO struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type PullRequestDTO struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
}

type MergeRequest struct {
	PullRequestId string `json:"pull_request_id"`
}

type ReASsign struct {
	PullRequestId string `json:"pull_request_id"`
	OldUserId     string `json:"old_user_id"`
}
