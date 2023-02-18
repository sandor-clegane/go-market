package userservice

import (
	"context"
	"encoding/base64"

	"github.com/sandor-clegane/go-market/internal/entities"
	"github.com/sandor-clegane/go-market/internal/entities/customerrors"
	"github.com/sandor-clegane/go-market/internal/storage/userstorage"
)

type userServiceImpl struct {
	userStorage userstorage.UserStorage
}

func (u userServiceImpl) Create(ctx context.Context, user entities.UserRequest, userID string) error {
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(user.Password))

	return u.userStorage.InsertUser(ctx, userID, user.Login, encodedPassword)
}

func (u userServiceImpl) Login(ctx context.Context, user entities.UserRequest) (string, error) {
	foundUser, err := u.userStorage.FindByLogin(ctx, user.Login)
	if err != nil {
		return "", err
	}
	decodedPassword, err := base64.StdEncoding.DecodeString(foundUser.Password)
	if err != nil {
		return "", err
	}
	if user.Password != string(decodedPassword) {
		return "", customerrors.NewInvalidPasswordError(user.Password)
	}

	return foundUser.ID, nil
}

func New(userRepository userstorage.UserStorage) UserService {
	return &userServiceImpl{
		userRepository,
	}
}
