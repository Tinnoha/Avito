package usecase

import (
	"Avito/internal/entity"
	"fmt"
	"math/rand"
)

type TeamRepository interface {
	GetReviewes(AuthorId string, teamName string) ([]string, error)
	NewReviewer(AuthorId string, OldReviewer string, TeamName string, PullRequestId string) ([]string, error)
	Create(entity.Team) (entity.Team, int, error)
	GetByName(TeamName string) (entity.Team, error)
	IsExists(TeamName string) bool
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

func (uc *TeamUseCase) GetReviewes(AuthorId string) ([]string, error) {
	exists := uc.userRepo.IsExists(AuthorId)

	if !exists {
		return nil, ErrNotFound
	}

	vasya, err := uc.userRepo.UserById(AuthorId)
	if err != nil {
		return nil, ErrUnexpected
	}

	reviewers, err := uc.teamRepo.GetReviewes(vasya.UserId, vasya.TeamName)
	if err != nil {
		return nil, ErrNoCandidate
	}
	var RandomNum int
	if len(reviewers) < 2 {
		RandomNum = rand.Intn(2)
	} else {
		RandomNum = rand.Intn(3)
	}

	if RandomNum > len(reviewers) {
		RandomNum = len(reviewers)
	}

	rand.Shuffle(len(reviewers), func(i, j int) {
		reviewers[i], reviewers[j] = reviewers[j], reviewers[i]
	})

	return reviewers[:RandomNum], nil
}

func (uc *TeamUseCase) NewReviewer(AuthorId string, OldReviewer string, PullRequestId string) (string, error) {
	exists := uc.userRepo.IsExists(AuthorId)
	if !exists {
		return "", ErrNotFound
	}

	exists = uc.userRepo.IsExists(OldReviewer)
	if !exists {
		return "", ErrNotFound
	}

	vasya, err := uc.userRepo.UserById(AuthorId)
	if err != nil {
		return "", ErrUnexpected
	}

	reviewers, err := uc.teamRepo.NewReviewer(vasya.UserId, OldReviewer, vasya.TeamName, PullRequestId)
	if err != nil {
		return "", ErrUnexpected
	}

	if len(reviewers) == 0 {
		return "", ErrNoCandidate
	}

	RandomMember := rand.Intn(len(reviewers))
	fmt.Println("2111")

	_, err = uc.userRepo.SetIsActive(OldReviewer, false)
	if err != nil {
		fmt.Println("w")
		return "", ErrUnexpected
	}

	return reviewers[RandomMember], nil
}

func (uc *TeamUseCase) Create(team entity.Team) (entity.Team, error) {
	for _, v := range team.Members {
		exists := uc.userRepo.IsExists(v.UserId)
		if exists {
			return entity.Team{}, ErrUserExists
		}
	}

	exists := uc.teamRepo.IsExists(team.TeamName)
	if exists {
		return entity.Team{}, ErrTeamExists
	}

	savedTeam, id, err := uc.teamRepo.Create(team)
	if err != nil {
		return entity.Team{}, ErrUnexpected
	}

	for _, v := range team.Members {
		vasya := entity.User{
			TeamName: team.TeamName,
			UserId:   v.UserId,
			Username: v.Username,
			IsActive: v.IsActive,
		}

		_, err := uc.userRepo.Save(vasya, id)
		if err != nil {
			return entity.Team{}, ErrUnexpected
		}
	}

	return savedTeam, nil
}

func (uc *TeamUseCase) GetByName(TeamName string) (entity.Team, error) {
	exists := uc.teamRepo.IsExists(TeamName)
	if !exists {
		return entity.Team{}, ErrNotFound
	}

	team, err := uc.teamRepo.GetByName(TeamName)
	if err != nil {
		return entity.Team{}, ErrUnexpected
	}

	return team, nil
}

func (uc *TeamUseCase) IsExists(TeamName string) bool {
	return uc.teamRepo.IsExists(TeamName)
}
