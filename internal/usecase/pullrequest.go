package usecase

import (
	"Avito/internal/entity"
	"errors"
	"time"
)

type PullRequestRepository interface {
	Create(entity.ShortPullRequest) (entity.PullRequest, error)
	Merge(PullRequestId string) (entity.PullRequest, error)
	Reassign(PullRequestId string, oldUserId string) (entity.PullRequest, error)
	RequestsById(UserId string) ([]entity.ShortPullRequest, error)
	IsExists(PullRequestId string) bool
	RequestByID(PullRequestId string) (entity.PullRequest, error)
}

type PullRequestUseCase struct {
	pullrequestRepo PullRequestRepository
	userRepo        UserRepository
	teamRepo        TeamRepository
}

func (uc *PullRequestUseCase) Create(feature entity.ShortPullRequest) (entity.PullRequest, error) {
	exist := uc.pullrequestRepo.IsExists(feature.PullRequestId)

	if exist {
		return entity.PullRequest{}, errors.New(entity.PR_EXISTS)
	}

	existAuthor := uc.userRepo.IsExists(feature.AuthorId)

	if !existAuthor {
		return entity.PullRequest{}, errors.New(entity.NOT_FOUND)
	}

	vasya, err := uc.userRepo.UserById(feature.AuthorId)

	if err != nil {
		return entity.PullRequest{}, errors.New(entity.NO_PREDICTED)
	}

	Reviewers, err := uc.teamRepo.GetReviewes(vasya.UserId, vasya.TeamName)

	if err != nil {
		return entity.PullRequest{}, errors.New(entity.NO_PREDICTED)
	}

	Request := entity.PullRequest{
		PullRequestId:     feature.PullRequestId,
		PullRequestName:   feature.PullRequestName,
		AuthorId:          feature.AuthorId,
		Status:            "open",
		AssignedReviewers: Reviewers,
		CreatedAt:         time.Now(),
		MergeAt:           nil,
	}

	return Request, nil
}

func (uc *PullRequestUseCase) Merge(PullRequestId string) (entity.PullRequest, error) {
	exists := uc.pullrequestRepo.IsExists(PullRequestId)

	if !exists {
		return entity.PullRequest{}, errors.New(entity.NOT_FOUND)
	}

	request, err := uc.pullrequestRepo.Merge(PullRequestId)

	if err != nil {
		return entity.PullRequest{}, errors.New(entity.NO_PREDICTED)
	}

	return request, nil
}

func (uc *PullRequestUseCase) Reassign(PullRequestId string, oldUserId string) (entity.PullRequest, error) {
	exists := uc.pullrequestRepo.IsExists(PullRequestId)

	if !exists {
		return entity.PullRequest{}, errors.New(entity.NOT_FOUND)
	}

	request, err := uc.pullrequestRepo.RequestByID(PullRequestId)

	if err != nil {
		return entity.PullRequest{}, errors.New(entity.NO_PREDICTED)
	}

	vasya, err := uc.userRepo.UserById(request.AuthorId)

	if err != nil {
		return entity.PullRequest{}, errors.New(entity.NO_PREDICTED)
	}

	newReviewer, err := uc.teamRepo.NewReviewer(oldUserId, vasya.UserId, vasya.TeamName)

	if err != nil {
		if errors.Is(err, errors.New(entity.NO_CANDIDATE)) {
			return entity.PullRequest{}, errors.New(entity.NO_CANDIDATE)
		}
		return entity.PullRequest{}, errors.New(entity.NO_PREDICTED)
	}

	for _, v := range request.AssignedReviewers {
		if v == oldUserId {
			v = newReviewer
		}
	}

	return request, nil
}

func (uc *PullRequestUseCase) RequestsById(UserId string) ([]entity.ShortPullRequest, error) {
	exists := uc.userRepo.IsExists(UserId)

	if !exists {
		return []entity.ShortPullRequest{}, errors.New(entity.NOT_FOUND)
	}

	Requests, err := uc.pullrequestRepo.RequestsById(UserId)

	if err != nil {
		return []entity.ShortPullRequest{}, errors.New(entity.NO_PREDICTED)
	}

	return Requests, nil
}

func (us *PullRequestUseCase) IsExists(PullRequestId string) bool {
	return us.pullrequestRepo.IsExists(PullRequestId)
}

func (uc *PullRequestUseCase) RequestByID(PullRequestId string) (entity.PullRequest, error) {
	exists := uc.pullrequestRepo.IsExists(PullRequestId)

	if !exists {
		return entity.PullRequest{}, errors.New(entity.NOT_FOUND)
	}

	request, err := uc.pullrequestRepo.RequestByID(PullRequestId)

	if err != nil {
		return entity.PullRequest{}, errors.New(entity.NO_PREDICTED)
	}

	return request, nil
}
