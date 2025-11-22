package controller

import (
	"Avito/internal/controller/dto"
	"Avito/internal/entity"
	"Avito/internal/usecase"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HTTPHandler struct {
	pullUse usecase.PullRequestUseCase
	teamUse usecase.TeamUseCase
	userUse usecase.UserUseCase
}

func NewHTTPHandler(
	pullUse usecase.PullRequestUseCase,
	teamUse usecase.TeamUseCase,
	userUse usecase.UserUseCase,
) *HTTPHandler {
	return &HTTPHandler{
		pullUse: pullUse,
		teamUse: teamUse,
		userUse: userUse,
	}
}

func HTTPError(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
) {

	errResponse := dto.ErrorResponse{
		Error: dto.Error{
			Code:    code,
			Message: message,
		},
	}

	b, err := json.MarshalIndent(errResponse, "", "    ")

	if err != nil {
		panic(err)
	}

	http.Error(w, string(b), status)
}

func (h *HTTPHandler) TeamAdd(w http.ResponseWriter, r *http.Request) {
	team := entity.Team{}

	err := json.NewDecoder(r.Body).Decode(&team)

	if err != nil {
		HTTPError(w, http.StatusBadRequest, entity.NO_PREDICTED, err.Error())
		return
	}

	team, err = h.teamUse.Create(team)

	if err != nil {
		if errors.Is(err, errors.New(entity.USER_EXISITS)) {
			HTTPError(w, http.StatusBadRequest, entity.USER_EXISITS,
				"Some of the user_id is exists yet")
			return
		} else if errors.Is(err, errors.New(entity.TEAM_EXISITS)) {
			HTTPError(w, http.StatusBadRequest, entity.TEAM_EXISITS,
				"This name of tea, is exists")
			return
		} else {
			HTTPError(w, http.StatusInternalServerError, entity.NO_PREDICTED, err.Error())
			return
		}
	}

	b, err := json.MarshalIndent(team, "", "    ")

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Error to give answer: \n", err.Error())
	}
}

func (h *HTTPHandler) TeamGet(w http.ResponseWriter, r *http.Request) {
	TeamName := r.URL.Query().Get("team_name")

	team, err := h.teamUse.GetByName(TeamName)

	if err != nil {
		if errors.Is(err, errors.New(entity.NOT_FOUND)) {
			HTTPError(w, http.StatusBadRequest, entity.NOT_FOUND,
				"This name of Team is not exsits")
			return
		} else {
			HTTPError(w, http.StatusInternalServerError, entity.NO_PREDICTED, err.Error())
			return
		}
	}

	b, err := json.MarshalIndent(team, "", "    ")

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Error to give answer: \n", err.Error())
	}
}

func (h *HTTPHandler) UserIsActive(w http.ResponseWriter, r *http.Request) {
	userdto := dto.UserActiveDTO{}

	err := json.NewDecoder(r.Body).Decode(&userdto)

	if err != nil {
		HTTPError(w, http.StatusInternalServerError, entity.NO_PREDICTED, err.Error())
		return
	}

	vasya, err := h.userUse.SetIsActive(userdto.UserId, userdto.IsActive)

	if err != nil {
		if errors.Is(err, errors.New(entity.NOT_FOUND)) {
			HTTPError(w, http.StatusBadRequest, entity.NOT_FOUND,
				"This user id is not exsits")
			return
		} else {
			HTTPError(w, http.StatusInternalServerError, entity.NO_PREDICTED, err.Error())
			return
		}
	}

	b, err := json.MarshalIndent(vasya, "", "    ")

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Error to give answer: \n", err.Error())
	}
}

func (h *HTTPHandler) RequestCreate(w http.ResponseWriter, r *http.Request) {
	request := dto.PullRequestDTO{}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		HTTPError(w, http.StatusInternalServerError, entity.NO_PREDICTED, err.Error())
		return
	}

	req := entity.ShortPullRequest{
		PullRequestId:   request.PullRequestId,
		PullRequestName: request.PullRequestName,
		AuthorId:        request.AuthorId,
		Status:          "open",
	}

	fullrequest, err := h.pullUse.Create(req)

	if err != nil {
		if errors.Is(err, errors.New(entity.PR_EXISTS)) {
			HTTPError(w, http.StatusBadRequest, entity.PR_EXISTS,
				"This request id is exsits yet")
			return
		} else if errors.Is(err, errors.New(entity.NOT_FOUND)) {
			HTTPError(w, http.StatusBadRequest, entity.NOT_FOUND,
				"This author is not exsists")
			return
		} else {
			HTTPError(w, http.StatusInternalServerError, entity.NO_PREDICTED, err.Error())
			return
		}
	}

	b, err := json.MarshalIndent(fullrequest, "", "    ")

	if err != nil {
		HTTPError(w, http.StatusInternalServerError, entity.NO_PREDICTED, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Error to give answer: \n", err.Error())
	}
}

func (h *HTTPHandler) RequestMerge(w http.ResponseWriter, r *http.Request) {
	requestDTO := dto.MergeRequest{}

	err := json.NewDecoder(r.Body).Decode(&requestDTO)
	if err != nil {
		HTTPError(w, http.StatusInternalServerError, entity.NO_PREDICTED, err.Error())
		return
	}

	request, err := h.pullUse.Merge(requestDTO.PullRequestId)

	if err != nil {
		if errors.Is(err, errors.New(entity.NOT_FOUND)) {
			HTTPError(w, http.StatusBadRequest, entity.NOT_FOUND,
				"This request is not exsist")
		} else {
			HTTPError(w, http.StatusInternalServerError, entity.NO_PREDICTED, err.Error())
			return
		}
	}

	b, err := json.MarshalIndent(request, "", "    ")

	if err != nil {
		HTTPError(w, http.StatusInternalServerError, entity.NO_PREDICTED, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Error to give answer: \n", err.Error())
	}

}

func (h *HTTPHandler) RequestReassign(w http.ResponseWriter, r *http.Request) {
	Resign := dto.ReASsign{}

	err := json.NewDecoder(r.Body).Decode(&Resign)

	if err != nil {
		HTTPError(w, http.StatusInternalServerError, entity.NO_PREDICTED, err.Error())
		return
	}

	request, err := h.pullUse.Reassign(Resign.PullRequestId, Resign.OldUserId)

	if err != nil {
		if errors.Is(err, errors.New(entity.NOT_FOUND)) {
			HTTPError(w, http.StatusNotFound, entity.NOT_FOUND,
				"This request is not exsist")
			return
		} else if errors.Is(err, errors.New(entity.NO_CANDIDATE)) {
			HTTPError(w, http.StatusConflict, entity.NOT_FOUND,
				"This team don`t have enough members to review")
			return
		} else if errors.Is(err, errors.New(entity.PR_MERGED)) {
			HTTPError(w, http.StatusConflict, entity.NOT_FOUND,
				"Request is merged yet")
			return
		} else if errors.Is(err, errors.New(entity.NOT_ASSIGNED)) {
			HTTPError(w, http.StatusConflict, entity.NOT_ASSIGNED,
				"This user is not a reviewer")
			return
		} else {
			HTTPError(w, http.StatusInternalServerError, entity.NO_PREDICTED, err.Error())
			return
		}
	}

	b, err := json.MarshalIndent(request, "", "    ")

	if err != nil {
		HTTPError(w, http.StatusInternalServerError, entity.NO_PREDICTED, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Error to give answer: \n", err.Error())
	}

}

func (h *HTTPHandler) UserGetRewiew(w http.ResponseWriter, r *http.Request) {
	vasya := r.URL.Query().Get("user_id")

	reviws, err := h.userUse.GetReview(vasya)

	if err != nil {
		if errors.Is(err, errors.New(entity.NOT_FOUND)) {
			HTTPError(w, http.StatusNotFound, entity.NOT_FOUND,
				"This user is not exsist")
			return
		} else if errors.Is(err, errors.New(entity.NO_PREDICTED)) {
			HTTPError(w, http.StatusInternalServerError, entity.NO_PREDICTED, err.Error())
			return
		}
	}

	answer := dto.Reviews{
		Reviews: reviws,
	}

	b, err := json.MarshalIndent(answer, "", "    ")

	if err != nil {
		HTTPError(w, http.StatusInternalServerError, entity.NO_PREDICTED, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Error to give answer: \n", err.Error())
	}
}
