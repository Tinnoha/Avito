package usecase

import (
	"Avito/internal/entity"
	"errors"
	"math/rand"
)

type TeamRepository interface {
	GetReviewes(AuthorId string, teamName string, count int) ([]string, error)
	NewReviewer(AuthorId string, OldReviewer string, TeamName string) (string, error)
	Create(entity.Team) (entity.Team, error)
	GetByName(TeamName string) (entity.Team, error)
	IsExists(TeamName string) bool
	CountActiveMembers(TeamName string) (int, error)
}

type TeamUseCase struct {
	teamRepo TeamRepository
	userRepo UserRepository
}

func NewTeamUseCase(
	userRepo UserRepository,
	teamRepo TeamRepository) *TeamUseCase {
	return &TeamUseCase{
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
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

	team, err := uс.teamRepo.GetByName(vasya.TeamName)

	if err != nil {
		return []string{}, errors.New(entity.NO_PREDICTED)
	}

	CountOfMembers, err := uс.teamRepo.CountActiveMembers(team.TeamName)

	if err != nil {
		return []string{}, errors.New(entity.NO_PREDICTED)
	}

	var RandomNum int
	if CountOfMembers-1 < 2 {
		RandomNum = rand.Intn(2)
	} else {
		RandomNum = rand.Intn(3)
	}

	reviewers, err := uс.teamRepo.GetReviewes(vasya.UserId, vasya.TeamName, RandomNum)

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
		return entity.Team{}, errors.New(entity.NO_PREDICTED)
	}

	return team, nil
}

func (uc *TeamUseCase) IsExists(TeamName string) bool {
	return uc.teamRepo.IsExists(TeamName)
}

func (uc *TeamUseCase) CountActiveMembers(TeamName string) (int, error) {
	exists := uc.teamRepo.IsExists(TeamName)

	if !exists {
		return 0, errors.New(entity.NOT_FOUND)
	}

	count, err := uc.teamRepo.CountActiveMembers(TeamName)

	if err != nil {
		return 0, errors.New(entity.NO_PREDICTED)
	}

	return count, nil
}
