package userHandler

import "net/http"

var _ UserHandler = &userHandlerImpl{}

type UserHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}
