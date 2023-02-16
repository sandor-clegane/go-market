package userHandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sandor-clegane/go-market/internal/entities"
	"github.com/sandor-clegane/go-market/internal/entities/customErrors"
	"github.com/sandor-clegane/go-market/internal/service/cookieService"
	"github.com/sandor-clegane/go-market/internal/service/userService"

	"github.com/google/uuid"
)

type userHandlerImpl struct {
	userService   userService.UserService
	cookieService cookieService.CookieService
}

func New(userService userService.UserService, cookieService cookieService.CookieService) UserHandler {
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
		var ve *customErrors.UserLoginUniqueViolation
		if errors.As(err, &ve) {
			writer.WriteHeader(http.StatusConflict)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	cookieErr := u.cookieService.WriteSigned(writer, userID)
	if cookieErr != nil {
		http.Error(writer, cookieErr.Error(), http.StatusInternalServerError)
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
		var ip *customErrors.InvalidPasswordError
		if errors.As(err, &ip) {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	cookieErr := u.cookieService.WriteSigned(writer, userID)
	if cookieErr != nil {
		http.Error(writer, cookieErr.Error(), http.StatusInternalServerError)
		return
	}
}
