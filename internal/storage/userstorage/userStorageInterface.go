package userstorage

import (
	"context"

	"github.com/sandor-clegane/go-market/internal/entities"
)

var us UserStorage = &userStorageImpl{}

type UserStorage interface {
	InsertUser(ctx context.Context, userID, login, password string) error
	FindByLogin(ctx context.Context, login string) (entities.User, error)
}
