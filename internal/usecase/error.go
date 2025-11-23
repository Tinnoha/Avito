package usecase

import (
	"Avito/internal/entity"
	"errors"
)

var (
	ErrNotFound    = errors.New(entity.NOT_FOUND)
	ErrNoCandidate = errors.New(entity.NO_CANDIDATE)
	ErrTeamExists  = errors.New(entity.TEAM_EXISITS)
	ErrUserExists  = errors.New(entity.USER_EXISITS)
	ErrUnexpected  = errors.New(entity.NO_PREDICTED)
	ErrPRExists    = errors.New(entity.PR_EXISTS)
	ErrPRMerged    = errors.New(entity.PR_MERGED)
	ErrNotAssigned = errors.New(entity.NOT_ASSIGNED)
)
