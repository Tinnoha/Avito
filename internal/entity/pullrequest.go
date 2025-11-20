package entity

import "time"

type PullRequest struct {
	PullRequestId     string
	PullRequestName   string
	AuthorId          string
	Status            string
	AssignedReviewers []string
	CreatedAt         time.Time
	MergeAt           time.Time
}

type ShortPullRequest struct {
	PullRequestId   string
	PullRequestName string
	AuthorId        string
	Status          string
}
