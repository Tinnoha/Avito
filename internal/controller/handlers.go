package controller

import (
	"Avito/internal/entity"
	"Avito/internal/usecase"
	"encoding/json"
	"net/http"
)

type HTTPHandler struct {
	pullUse usecase.PullRequestUseCase
	teamUse usecase.TeamUseCase
	userUse usecase.UserUseCase
}

func HTTPError (
	w http.ResponseWriter,
	status int,
	code string,
	message string,
	){

	errResponse := entity.ErrorResponse{
		Error: entity.Error{
			Code: code,
			Message: message,
		},
	}

	b,err := json.MarshalIndent(errResponse,"","    ")

	if err != nil{
		panic(err)
	}

	http.Error(w,string(b),status)
}

func (h *HTTPHandler) TeamAdd(w http.ResponseWriter, r *http.Request) {
	team := entity.Team{}

	err := json.NewDecoder(r.Body).Decode(&team)

	if err != nil{
		HTTPError(w,http.StatusBadRequest,entity.NO_PREDICTED,err.Error())
	}

	team, errResponse := h.teamUse.Create(team)

	if errResponse

}

func (h *HTTPHandler) TeamGet(w http.ResponseWriter, r *http.Request) {

}

func (h *HTTPHandler) UserIsActive(w http.ResponseWriter, r *http.Request) {

}

func (h *HTTPHandler) RequsetCreate(w http.ResponseWriter, r *http.Request) {

}

func (h *HTTPHandler) RequestMerge(w http.ResponseWriter, r *http.Request) {

}

func (h *HTTPHandler) RequestReassign(w http.ResponseWriter, r *http.Request) {

}

func (h *HTTPHandler) UserGetRewiew(w http.ResponseWriter, r *http.Request) {

}
