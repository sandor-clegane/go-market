package userhandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sandor-clegane/go-market/internal/entities"
	"github.com/sandor-clegane/go-market/internal/entities/customerrors"
	"github.com/sandor-clegane/go-market/internal/service/cookieservice"
	"github.com/sandor-clegane/go-market/internal/service/userservice"

	"github.com/google/uuid"
)

type userHandlerImpl struct {
	userService   userservice.UserService
	cookieService cookieservice.CookieService
}

func New(userService userservice.UserService, cookieService cookieservice.CookieService) UserHandler {
	return &userHandlerImpl{userService, cookieService}
}

//Create  Хендлер: POST /api/user/register.
//Регистрация производится по паре логин/пароль. Каждый логин должен быть уникальным.
//После успешной регистрации должна происходить автоматическая аутентификация пользователя.
func (u *userHandlerImpl) Create(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var user entities.UserRequest
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	userID := uuid.New().String()
	err = u.userService.Create(request.Context(), user, userID)

	if err != nil {
		var ve *customerrors.LoginUniqueViolationError
		if errors.As(err, &ve) {
			writer.WriteHeader(http.StatusConflict)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	err = u.cookieService.WriteSigned(writer, userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

//Login Хендлер: POST /api/user/login.
//Аутентификация производится по паре логин/пароль.
func (u *userHandlerImpl) Login(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var user entities.UserRequest
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	userID, err := u.userService.Login(request.Context(), user)
	if err != nil {
		var ip *customerrors.InvalidPasswordError
		if errors.As(err, &ip) {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = u.cookieService.WriteSigned(writer, userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
