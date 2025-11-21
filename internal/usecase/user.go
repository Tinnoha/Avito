package usecase

import (
	"Avito/internal/entity"
	"errors"
)

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

func (uc *UserUseCase) Save(vasya entity.User) (entity.User, error) {
	exist := uc.userRepo.IsExists(vasya.UserId)

	if exist {
		return entity.User{}, errors.New(entity.USER_EXISITS)
	}

	vasya, err := uc.userRepo.Save(vasya)
	if err != nil {
		return entity.User{}, errors.New(entity.NO_PREDICTED)
	}

	return vasya, nil
}

func (uc *UserUseCase) SetIsActive(userId string, isActive bool) (entity.User, error) {
	exist := uc.userRepo.IsExists(userId)

	if exist {
		return entity.User{}, errors.New(entity.USER_EXISITS)
	}

	vasya, err := uc.userRepo.SetIsActive(userId, isActive)
	if err != nil {
		return entity.User{}, errors.New(entity.NO_PREDICTED)
	}

	return vasya, nil
}

func (uc *UserUseCase) GetReview(userId string) ([]entity.ShortPullRequest, error) {
	exist := uc.userRepo.IsExists(userId)

	if exist {
		return []entity.ShortPullRequest{}, errors.New(entity.USER_EXISITS)
	}

	requests, err := uc.pullrequestRepo.RequestsById(userId)

	if err != nil {
		return []entity.ShortPullRequest{}, errors.New(entity.NO_PREDICTED)
	}

	return requests, nil
}

func (uc *UserUseCase) IsExists(userId string) bool {
	return uc.userRepo.IsExists(userId)
}

func (uc *UserUseCase) UserById(userId string) (entity.User, error) {
	exist := uc.userRepo.IsExists(userId)

	if exist {
		return entity.User{}, errors.New(entity.USER_EXISITS)
	}

	kolya, err := uc.userRepo.UserById(userId)

	if err != nil {
		return entity.User{}, errors.New(entity.NO_PREDICTED)
	}

	return kolya, nil
}
