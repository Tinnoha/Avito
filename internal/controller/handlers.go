package controller

import (
	"Avito/internal/usecase"
	"net/http"
)

type HTTPHandler struct {
	pullUse usecase.PullRequestUseCase
	teamUse usecase.TeamUseCase
	userUse usecase.UserUseCase
}

func (h *HTTPHandler) TeamAdd(w http.ResponseWriter, r *http.Request) {

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
