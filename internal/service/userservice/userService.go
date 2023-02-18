package userservice

import (
	"context"
	"crypto/md5"
	"io"

	"github.com/sandor-clegane/go-market/internal/entities"
	"github.com/sandor-clegane/go-market/internal/entities/customerrors"
	"github.com/sandor-clegane/go-market/internal/storage/userstorage"
)

type userServiceImpl struct {
	userStorage userstorage.UserStorage
}

func (u userServiceImpl) Create(ctx context.Context, user entities.UserRequest, userID string) error {
	h := md5.New()
	io.WriteString(h, user.Password)
	encodedPassword := string(h.Sum(nil))

	return u.userStorage.InsertUser(ctx, userID, user.Login, encodedPassword)
}

func (u userServiceImpl) Login(ctx context.Context, user entities.UserRequest) (string, error) {
	//expected hash
	foundUser, err := u.userStorage.FindByLogin(ctx, user.Login)
	if err != nil {
		return "", err
	}
	//input hash
	h := md5.New()
	io.WriteString(h, user.Password)
	encodedPassword := string(h.Sum(nil))

	if foundUser.Password != encodedPassword {
		return "", customerrors.NewInvalidPasswordError(user.Password)
	}

	return foundUser.ID, nil
}

func New(userRepository userstorage.UserStorage) UserService {
	return &userServiceImpl{
		userRepository,
	}
}
