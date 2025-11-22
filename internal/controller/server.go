package controller

import (
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	Hendlers HTTPHandler
}

func NewHTTPServer(Hendlers HTTPHandler) *HTTPServer {
	return &HTTPServer{
		Hendlers: Hendlers,
	}
}

func (s *HTTPServer) Run() error {
	router := mux.NewRouter()

	router.Path("/team/add").Methods(http.MethodPost).HandlerFunc(s.Hendlers.TeamAdd)
	router.Path("/team/get").Methods(http.MethodGet).HandlerFunc(s.Hendlers.TeamGet)

	router.Path("/users/setIsActive").Methods(http.MethodPost).HandlerFunc(s.Hendlers.UserIsActive)
	router.Path("/users/getReview").Methods(http.MethodGet).HandlerFunc(s.Hendlers.UserGetRewiew)

	router.Path("/pullRequest/create").Methods(http.MethodPost).HandlerFunc(s.Hendlers.RequestCreate)
	router.Path("/pullRequest/merge").Methods(http.MethodPost).HandlerFunc(s.Hendlers.RequestMerge)
	router.Path("/pullRequest/reassign").Methods(http.MethodPost).HandlerFunc(s.Hendlers.RequestReassign)

	return http.ListenAndServe(":8080", router)
}
