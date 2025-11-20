package usecase

import "Avito/internal/entity"

type UserRepository interface {
	Save(entity.User) (entity.User, error)
	SetIsActive(userId string, isActive bool) (entity.User, error)
	GetReview(userId string) ([]entity.ShortPullRequest, error)
	IsExists(userId string) bool
	UserById(userId string) (entity.User, error)
}

type UserUseCase struct {
	userRepo        UserRepository
	pullrequestRepo PullRequestRepository
}

func (uc *UserUseCase) Save(vasya entity.User) (entity.User, entity.ErrorResponse) {
	exist := uc.userRepo.IsExists(vasya.UserId)

	if exist {
		err := entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.USER_EXISITS,
				Message: "This user already exists",
			},
		}

		return entity.User{}, err
	}

	vasya, err := uc.userRepo.Save(vasya)
	if err != nil {
		return entity.User{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
	}

	return vasya, entity.ErrorResponse{}
}

func (uc *UserUseCase) SetIsActive(userId string, isActive bool) (entity.User, entity.ErrorResponse) {
	exist := uc.userRepo.IsExists(userId)

	if exist {
		err := entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.USER_EXISITS,
				Message: "This user already exists",
			},
		}

		return entity.User{}, err
	}

	vasya, err := uc.userRepo.SetIsActive(userId, isActive)
	if err != nil {
		return entity.User{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
	}

	return vasya, entity.ErrorResponse{}
}

func (uc *UserUseCase) GetReview(userId string) ([]entity.ShortPullRequest, entity.ErrorResponse) {
	exist := uc.userRepo.IsExists(userId)

	if exist {
		err := entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.USER_EXISITS,
				Message: "This user already exists",
			},
		}

		return []entity.ShortPullRequest{}, err
	}

	requests, err := uc.pullrequestRepo.RequestsById(userId)

	if err != nil {
		return []entity.ShortPullRequest{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
	}

	return requests, entity.ErrorResponse{}
}

func (uc *UserUseCase) IsExists(userId string) bool {
	return uc.userRepo.IsExists(userId)
}

func (uc *UserUseCase) UserById(userId string) (entity.User, entity.ErrorResponse) {
	exist := uc.userRepo.IsExists(userId)

	if exist {
		err := entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.USER_EXISITS,
				Message: "This user already exists",
			},
		}

		return entity.User{}, err
	}

	kolya, err := uc.userRepo.UserById(userId)

	if err != nil {
		return entity.User{}, entity.ErrorResponse{
			Error: entity.Error{
				Code:    entity.NO_PREDICTED,
				Message: err.Error(),
			},
		}
	}

	return kolya, entity.ErrorResponse{}
}
