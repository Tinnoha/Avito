package usecase

import (
	"Avito/internal/entity"
	"errors"
	"fmt"
)

type TeamRepository interface {
	GetReviewes(AuthorId string, teamName string) ([]string, error)
	NewReviewer(AuthorId string, OldReviewer string, TeamName string) (string, error)
	Create(entity.Team) (entity.Team, error)
	GetByName(TeamName string) (entity.Team, error)
	IsExists(TeamName string) bool
}

type TeamUseCase struct {
	teamRepo TeamRepository
	userRepo UserRepository
	pullRepo PullRequestRepository
}

func (uс *TeamUseCase) GetReviewes(AuthorId string) ([]string, entity.ErrorResponse) {
	exists := uс.userRepo.IsExists(AuthorId)

	if !exists {
		return []string{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NOT_FOUND,
				Message: "Author is not exists",
			},
		}
	}

	vasya, err := uс.userRepo.UserById(AuthorId)

	if err != nil {
		return []string{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
	}

	reviewers, err := uс.teamRepo.GetReviewes(vasya.UserId, vasya.TeamName)

	if err != nil {
		if errors.Is(err, errors.New(entity.NO_CANDIDATE)) {
			return []string{}, entity.ErrorResponse{
				Error: entity.Error{
					Code:    entity.NO_CANDIDATE,
					Message: "No candidate for this reviews",
				},
			}
		} else {
			return []string{}, entity.ErrorResponse{
				Error: entity.Error{
					Code:    entity.NO_PREDICTED,
					Message: err.Error(),
				},
			}
		}
	}

	return reviewers, entity.ErrorResponse{}
}

func (uc *TeamUseCase) NewReviewer(AuthorId string, OldReviewer string) (string, entity.ErrorResponse) {
	exists := uc.userRepo.IsExists(AuthorId)

	if !exists {
		return "", entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NOT_FOUND,
				Message: "Author is not exists",
			},
		}
	}

	exists = uc.userRepo.IsExists(OldReviewer)

	if !exists {
		return "", entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NOT_FOUND,
				Message: "Author is not exists",
			},
		}
	}

	vasya, err := uc.userRepo.UserById(AuthorId)

	if err != nil {
		return "", entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
	}

	reviewers, err := uc.teamRepo.NewReviewer(vasya.UserId, OldReviewer, vasya.TeamName)

	if err != nil {
		return "", entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
	}

	_, err = uc.userRepo.SetIsActive(OldReviewer, false)

	if err != nil {
		return "", entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
	}

	return reviewers, entity.ErrorResponse{}
}

func (uc *TeamUseCase) Create(team entity.Team) (entity.Team, entity.ErrorResponse) {
	for _, v := range team.Members {
		exists := uc.userRepo.IsExists(v.UserId)

		if exists {
			return entity.Team{}, entity.ErrorResponse{
				Error: entity.Error{
					Code:    entity.USER_EXISITS,
					Message: fmt.Sprintf("User Id %s is exists yet", v.UserId),
				},
			}
		}
	}

	exists := uc.teamRepo.IsExists(team.TeamName)

	if exists {
		return entity.Team{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.TEAM_EXISITS,
				Message: fmt.Sprintf("Team name %s is exists yet", team.TeamName),
			},
		}
	}

	savedTeam, err := uc.teamRepo.Create(team)

	if err != nil {
		return entity.Team{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
	}

	return savedTeam, entity.ErrorResponse{}
}

func (uc *TeamUseCase) GetByName(TeamName string) (entity.Team, entity.ErrorResponse) {
	exists := uc.teamRepo.IsExists(TeamName)

	if !exists {
		return entity.Team{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NOT_FOUND,
				Message: fmt.Sprintf("Team name %s is not exists yet", TeamName),
			},
		}
	}

	team, err := uc.teamRepo.GetByName(TeamName)

	if err != nil {
		return entity.Team{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
	}

	return team, entity.ErrorResponse{}
}

func (uc *TeamUseCase) IsExists(TeamName string) bool {
	return uc.teamRepo.IsExists(TeamName)
}
