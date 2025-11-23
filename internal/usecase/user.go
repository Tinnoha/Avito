package usecase

import (
	"Avito/internal/entity"
)

type UserRepository interface {
	Save(entity.User, int) (entity.User, error)
	SetIsActive(userId string, isActive bool) (entity.User, error)
	IsExists(userId string) bool
	UserById(userId string) (entity.User, error)
}

type UserUseCase struct {
	userRepo        UserRepository
	pullrequestRepo PullRequestRepository
}

func NewUserUseCase(
	userRepo UserRepository,
	pullrequestRepo PullRequestRepository) *UserUseCase {
	return &UserUseCase{
		userRepo:        userRepo,
		pullrequestRepo: pullrequestRepo,
	}
}

func (uc *UserUseCase) Save(vasya entity.User, id int) (entity.User, error) {
	exist := uc.userRepo.IsExists(vasya.UserId)
	if exist {
		return entity.User{}, ErrUserExists
	}

	vasya, err := uc.userRepo.Save(vasya, id)
	if err != nil {
		return entity.User{}, ErrUnexpected
	}

	return vasya, nil
}

func (uc *UserUseCase) SetIsActive(userId string, isActive bool) (entity.User, error) {
	exist := uc.userRepo.IsExists(userId)
	if !exist {
		return entity.User{}, ErrNotFound
	}

	vasya, err := uc.userRepo.SetIsActive(userId, isActive)
	if err != nil {
		return entity.User{}, ErrUnexpected
	}

	return vasya, nil
}

func (uc *UserUseCase) GetReview(userId string, all bool) ([]entity.ShortPullRequest, error) {
	exist := uc.userRepo.IsExists(userId)
	if !exist {
		return nil, ErrNotFound
	}

	requests, err := uc.pullrequestRepo.RequestsById(userId, all)
	if err != nil {
		return nil, ErrUnexpected
	}

	return requests, nil
}

func (uc *UserUseCase) IsExists(userId string) bool {
	return uc.userRepo.IsExists(userId)
}

func (uc *UserUseCase) UserById(userId string) (entity.User, error) {
	exist := uc.userRepo.IsExists(userId)
	if !exist {
		return entity.User{}, ErrNotFound
	}

	user, err := uc.userRepo.UserById(userId)
	if err != nil {
		return entity.User{}, ErrUnexpected
	}

	return user, nil
}
