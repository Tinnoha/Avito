package entity

import "time"

type PullRequest struct {
	PullRequestId     string     `db:"pull_request_id" json:"pull_request_id"`
	PullRequestName   string     `db:"pull_request_name" json:"pull_request_name"`
	AuthorId          string     `db:"author_id" json:"author_id"`
	Status            string     `db:"status" json:"status"`
	AssignedReviewers []string   `db:"assigned_reviewers" json:"assigned_reviewers"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	MergeAt           *time.Time `db:"merged_at" json:"merged_at"`
}

type ShortPullRequest struct {
	PullRequestId   string `db:"pull_request_id" json:"pull_request_id"`
	PullRequestName string `db:"pull_request_name" json:"pull_request_name"`
	AuthorId        string `db:"author_id" json:"author_id"`
	Status          string `db:"status" json:"status"`
}
