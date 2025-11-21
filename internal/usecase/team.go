package usecase

import (
	"Avito/internal/entity"
	"errors"
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
}

func (uс *TeamUseCase) GetReviewes(AuthorId string) ([]string, error) {
	exists := uс.userRepo.IsExists(AuthorId)

	if !exists {
		return []string{}, errors.New(entity.NOT_FOUND)
	}

	vasya, err := uс.userRepo.UserById(AuthorId)

	if err != nil {
		return []string{}, errors.New(entity.NO_PREDICTED)
	}

	reviewers, err := uс.teamRepo.GetReviewes(vasya.UserId, vasya.TeamName)

	if err != nil {
		if errors.Is(err, errors.New(entity.NO_CANDIDATE)) {
			return []string{}, errors.New(entity.NO_CANDIDATE)
		} else {
			return []string{}, errors.New(entity.NO_PREDICTED)
		}
	}

	return reviewers, nil
}

func (uc *TeamUseCase) NewReviewer(AuthorId string, OldReviewer string) (string, error) {
	exists := uc.userRepo.IsExists(AuthorId)

	if !exists {
		return "", errors.New(entity.NOT_FOUND)
	}

	exists = uc.userRepo.IsExists(OldReviewer)

	if !exists {
		return "", errors.New(entity.NOT_FOUND)
	}

	vasya, err := uc.userRepo.UserById(AuthorId)

	if err != nil {
		return "", errors.New(entity.NO_PREDICTED)
	}

	reviewers, err := uc.teamRepo.NewReviewer(vasya.UserId, OldReviewer, vasya.TeamName)

	if err != nil {
		return "", errors.New(entity.NO_PREDICTED)
	}

	_, err = uc.userRepo.SetIsActive(OldReviewer, false)

	if err != nil {
		return "", errors.New(entity.NO_PREDICTED)
	}

	return reviewers, nil
}

func (uc *TeamUseCase) Create(team entity.Team) (entity.Team, error) {
	for _, v := range team.Members {
		exists := uc.userRepo.IsExists(v.UserId)

		if exists {
			return entity.Team{}, errors.New(entity.USER_EXISITS)
		}
	}

	exists := uc.teamRepo.IsExists(team.TeamName)

	if exists {
		return entity.Team{}, errors.New(entity.TEAM_EXISITS)
	}

	savedTeam, err := uc.teamRepo.Create(team)

	if err != nil {
		return entity.Team{}, errors.New(entity.NO_PREDICTED)
	}

	return savedTeam, nil
}

func (uc *TeamUseCase) GetByName(TeamName string) (entity.Team, error) {
	exists := uc.teamRepo.IsExists(TeamName)

	if !exists {
		return entity.Team{}, errors.New(entity.NOT_FOUND)
	}

	team, err := uc.teamRepo.GetByName(TeamName)

	if err != nil {
		return entity.Team{}, errors.New(entity.NOT_FOUND)
	}

	return team, nil
}

func (uc *TeamUseCase) IsExists(TeamName string) bool {
	return uc.teamRepo.IsExists(TeamName)
}
