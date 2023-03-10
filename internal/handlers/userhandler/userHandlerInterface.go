package userhandler

import "net/http"

var _ UserHandler = (*userHandlerImpl)(nil)

type UserHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}
