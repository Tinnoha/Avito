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

func (uc *PullRequestUseCase) Create(feature entity.ShortPullRequest) (entity.PullRequest, entity.ErrorResponse) {
	exist := uc.pullrequestRepo.IsExists(feature.PullRequestId)

	if exist {
		err := entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.PR_EXISTS,
				Message: "This pull request already exists",
			},
		}

		return entity.PullRequest{}, err
	}

	existAuthor := uc.userRepo.IsExists(feature.AuthorId)

	if !existAuthor {
		err := entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NOT_FOUND,
				Message: "This author is not exists",
			},
		}

		return entity.PullRequest{}, err
	}

	vasya, err := uc.userRepo.UserById(feature.AuthorId)

	if err != nil {
		return entity.PullRequest{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
	}

	Reviewers, err := uc.teamRepo.GetReviewes(vasya.UserId, vasya.TeamName)

	if err != nil {
		return entity.PullRequest{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
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

	return Request, entity.ErrorResponse{}
}

func (uc *PullRequestUseCase) Merge(PullRequestId string) (entity.PullRequest, entity.ErrorResponse) {
	exists := uc.pullrequestRepo.IsExists(PullRequestId)

	if !exists {
		return entity.PullRequest{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NOT_FOUND,
				Message: "This pull request is not exists",
			},
		}
	}

	request, err := uc.pullrequestRepo.Merge(PullRequestId)

	if err != nil {
		return entity.PullRequest{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
	}

	return request, entity.ErrorResponse{}
}

func (uc *PullRequestUseCase) Reassign(PullRequestId string, oldUserId string) (entity.PullRequest, entity.ErrorResponse) {
	exists := uc.pullrequestRepo.IsExists(PullRequestId)

	if !exists {
		return entity.PullRequest{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NOT_FOUND,
				Message: "This pull request is not exists",
			},
		}
	}

	request, err := uc.pullrequestRepo.RequestByID(PullRequestId)

	if err != nil {
		return entity.PullRequest{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
	}

	vasya, err := uc.userRepo.UserById(request.AuthorId)

	if err != nil {
		return entity.PullRequest{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
	}

	newReviewer, err := uc.teamRepo.NewReviewer(oldUserId, vasya.UserId, vasya.TeamName)

	if err != nil {
		if errors.Is(err, errors.New(entity.NO_CANDIDATE)) {
			return entity.PullRequest{}, entity.ErrorResponse{
				Error: entity.Error{
					Code:    entity.NO_CANDIDATE,
					Message: err.Error(),
				},
			}
		}
		return entity.PullRequest{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
	}

	for _, v := range request.AssignedReviewers {
		if v == oldUserId {
			v = newReviewer
		}
	}

	return request, entity.ErrorResponse{}
}

func (uc *PullRequestUseCase) RequestsById(UserId string) ([]entity.ShortPullRequest, entity.ErrorResponse) {
	exists := uc.userRepo.IsExists(UserId)

	if !exists {
		return []entity.ShortPullRequest{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NOT_FOUND,
				Message: "This user is not exists",
			},
		}
	}

	Requests, err := uc.pullrequestRepo.RequestsById(UserId)

	if err != nil {
		return []entity.ShortPullRequest{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
	}

	return Requests, entity.ErrorResponse{}
}

func (us *PullRequestUseCase) IsExists(PullRequestId string) bool {
	return us.pullrequestRepo.IsExists(PullRequestId)
}

func (uc *PullRequestUseCase) RequestByID(PullRequestId string) (entity.PullRequest, entity.ErrorResponse) {
	exists := uc.pullrequestRepo.IsExists(PullRequestId)

	if !exists {
		return entity.PullRequest{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NOT_FOUND,
				Message: "This pull request is not exists",
			},
		}
	}

	request, err := uc.pullrequestRepo.RequestByID(PullRequestId)

	if err != nil {
		return entity.PullRequest{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
	}

	return request, entity.ErrorResponse{}
}
