package usecase

import "Avito/internal/entity"

type TeamRepository interface {
	GetReviewes(AuthorId string) ([]string, error)
	NewReviewer(AuthorId string, OldReviewer string) (string, error)
	Create(entity.Team) (entity.Team, error)
	GetByName(TeamName string) (entity.Team, error)
	IsExists(TeamName string) bool
}
