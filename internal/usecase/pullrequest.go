package usecase

import (
	"Avito/internal/entity"
	"fmt"
	"math/rand"
)

type PullRequestRepository interface {
	Create(request entity.ShortPullRequest, reviewers []string) (entity.PullRequest, error)
	Merge(PullRequestId string) (entity.PullRequest, error)
	Reassign(PullRequestId string, oldUserId string, newUserId string) error
	RequestsById(UserId string, all bool) ([]entity.ShortPullRequest, error)
	IsExists(PullRequestId string) bool
	RequestByID(PullRequestId string) (entity.PullRequest, error)
	IsMerged(PullRequestId string) (bool, error)
}

type PullRequestUseCase struct {
	pullrequestRepo PullRequestRepository
	userRepo        UserRepository
	teamRepo        TeamRepository
}

func NewPullRequestUseCase(
	pullrequestRepo PullRequestRepository,
	userRepo UserRepository,
	teamRepo TeamRepository) *PullRequestUseCase {
	return &PullRequestUseCase{
		pullrequestRepo: pullrequestRepo,
		userRepo:        userRepo,
		teamRepo:        teamRepo,
	}
}

func (uc *PullRequestUseCase) Create(feature entity.ShortPullRequest) (entity.PullRequest, error) {
	exist := uc.pullrequestRepo.IsExists(feature.PullRequestId)
	if exist {
		return entity.PullRequest{}, ErrPRExists
	}

	existAuthor := uc.userRepo.IsExists(feature.AuthorId)
	if !existAuthor {
		return entity.PullRequest{}, ErrNotFound
	}

	vasya, err := uc.userRepo.UserById(feature.AuthorId)
	if err != nil {
		return entity.PullRequest{}, ErrUnexpected
	}

	Reviewers, err := uc.teamRepo.GetReviewes(vasya.UserId, vasya.TeamName)
	if err != nil {
		return entity.PullRequest{}, ErrNoCandidate
	}

	if len(Reviewers) == 0 {
		return entity.PullRequest{}, ErrNoCandidate
	}

	var RandomNum int
	if len(Reviewers) < 2 {
		RandomNum = rand.Intn(2)
	} else {
		RandomNum = rand.Intn(3)
	}

	if RandomNum > len(Reviewers) {
		RandomNum = len(Reviewers)
	}

	rand.Shuffle(len(Reviewers), func(i, j int) {
		Reviewers[i], Reviewers[j] = Reviewers[j], Reviewers[i]
	})

	Reviewers = Reviewers[:RandomNum]

	request := entity.ShortPullRequest{
		PullRequestId:   feature.PullRequestId,
		PullRequestName: feature.PullRequestName,
		AuthorId:        feature.AuthorId,
		Status:          "open",
	}

	Request, err := uc.pullrequestRepo.Create(request, Reviewers)

	if err != nil {
		return entity.PullRequest{}, ErrUnexpected
	}

	return Request, nil
}

func (uc *PullRequestUseCase) Merge(PullRequestId string) (entity.PullRequest, error) {
	exists := uc.pullrequestRepo.IsExists(PullRequestId)
	if !exists {
		return entity.PullRequest{}, ErrNotFound
	}

	if merged, _ := uc.pullrequestRepo.IsMerged(PullRequestId); merged {
		request, err := uc.pullrequestRepo.RequestByID(PullRequestId)
		if err != nil {
			return entity.PullRequest{}, ErrUnexpected
		}
		return request, nil
	}

	request, err := uc.pullrequestRepo.Merge(PullRequestId)
	if err != nil {
		return entity.PullRequest{}, ErrUnexpected
	}

	return request, nil
}

func (uc *PullRequestUseCase) Reassign(PullRequestId string, oldUserId string) (entity.PullRequest, error) {
	exists := uc.pullrequestRepo.IsExists(PullRequestId)
	if !exists {
		return entity.PullRequest{}, ErrNotFound
	}

	merged, err := uc.pullrequestRepo.IsMerged(PullRequestId)
	if err != nil {
		return entity.PullRequest{}, ErrUnexpected
	}

	if merged {
		return entity.PullRequest{}, ErrPRMerged
	}

	request, err := uc.pullrequestRepo.RequestByID(PullRequestId)
	if err != nil {
		return entity.PullRequest{}, ErrUnexpected
	}

	var IsReviewer bool
	for _, v := range request.AssignedReviewers {
		if v == oldUserId {
			IsReviewer = true
			break
		}
	}

	if !IsReviewer {
		return entity.PullRequest{}, ErrNotAssigned
	}

	exists = uc.userRepo.IsExists(oldUserId)
	if !exists {
		return entity.PullRequest{}, ErrNotFound
	}

	vasya, err := uc.userRepo.UserById(request.AuthorId)
	if err != nil {
		return entity.PullRequest{}, ErrUnexpected
	}

	newReviewers, err := uc.teamRepo.NewReviewer(vasya.UserId, oldUserId, vasya.TeamName, PullRequestId)
	if err != nil {
		return entity.PullRequest{}, ErrUnexpected
	}

	if len(newReviewers) == 0 {
		return entity.PullRequest{}, ErrNoCandidate
	}

	RandomMember := rand.Intn(len(newReviewers))

	for k, v := range request.AssignedReviewers {
		if v == oldUserId {
			request.AssignedReviewers[k] = newReviewers[RandomMember]
		}
	}

	err = uc.pullrequestRepo.Reassign(PullRequestId, oldUserId, newReviewers[RandomMember])

	if err != nil {
		fmt.Println("10")
		fmt.Println(err)
		return entity.PullRequest{}, ErrUnexpected
	}

	_, err = uc.userRepo.SetIsActive(oldUserId, false)

	if err != nil {
		fmt.Println("11")
		fmt.Println(err)
		return entity.PullRequest{}, ErrUnexpected
	}

	return request, nil
}

func (uc *PullRequestUseCase) RequestsById(UserId string, all bool) ([]entity.ShortPullRequest, error) {
	exists := uc.userRepo.IsExists(UserId)
	if !exists {
		return nil, ErrNotFound
	}

	Requests, err := uc.pullrequestRepo.RequestsById(UserId, all)
	if err != nil {
		return nil, ErrUnexpected
	}

	return Requests, nil
}

func (uc *PullRequestUseCase) IsExists(PullRequestId string) bool {
	return uc.pullrequestRepo.IsExists(PullRequestId)
}

func (uc *PullRequestUseCase) RequestByID(PullRequestId string) (entity.PullRequest, error) {
	exists := uc.pullrequestRepo.IsExists(PullRequestId)
	if !exists {
		return entity.PullRequest{}, ErrNotFound
	}

	request, err := uc.pullrequestRepo.RequestByID(PullRequestId)
	if err != nil {
		return entity.PullRequest{}, ErrUnexpected
	}

	return request, nil
}

func (uc *PullRequestUseCase) IsMerged(PullRequestId string) (bool, error) {
	if exists := uc.pullrequestRepo.IsExists(PullRequestId); !exists {
		return false, ErrNotFound
	}

	return uc.pullrequestRepo.IsMerged(PullRequestId)
}
