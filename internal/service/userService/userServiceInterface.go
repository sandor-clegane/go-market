package userService

import (
	"context"

	"github.com/sandor-clegane/go-market/internal/entities"
)

var _ UserService = &userServiceImpl{}

type UserService interface {
	Create(ctx context.Context, user entities.UserRequest, userID string) error
	Login(ctx context.Context, user entities.UserRequest) (string, error)
}
