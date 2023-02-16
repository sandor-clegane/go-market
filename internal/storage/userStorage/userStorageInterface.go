package userStorage

import (
	"context"

	"Gophermarket/go-musthave-diploma-tpl/internal/entities"
)

var us UserStorage = &userStorageImpl{}

type UserStorage interface {
	InsertUser(ctx context.Context, userID, login, password string) error
	FindByLogin(ctx context.Context, login string) (entities.User, error)
}
