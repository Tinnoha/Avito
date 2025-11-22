package dto

import "Avito/internal/entity"

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Reviews struct {
	Reviews []entity.ShortPullRequest `json:"pull_requests"`
}
