package userStorage

import (
	"context"
	"database/sql"

	"Gophermarket/go-musthave-diploma-tpl/internal/entities"
	"Gophermarket/go-musthave-diploma-tpl/internal/entities/customErrors"

	"github.com/omeid/pgerror"
)

const (
	insertUserQuery = "" +
		"INSERT INTO users (id, login, password) " +
		"VALUES ($1, $2, $3)"
	getUserByLoginQuery = "" +
		"SELECT id, login, password " +
		"FROM users " +
		"WHERE login = $1"
)

type userStorageImpl struct {
	db *sql.DB
}

func (us *userStorageImpl) InsertUser(ctx context.Context, userID, login, password string) error {
	_, err := us.db.ExecContext(ctx, insertUserQuery, userID, login, password)
	if err != nil {
		if e := pgerror.UniqueViolation(err); e != nil {
			return customErrors.NewUserLoginUniqueViolationError(login, err)
		}
		return err
	}
	return nil
}

func (us *userStorageImpl) FindByLogin(ctx context.Context, login string) (entities.User, error) {
	row := us.db.QueryRowContext(ctx, getUserByLoginQuery, login)
	var user entities.User
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func New(db *sql.DB) UserStorage {
	return &userStorageImpl{
		db: db,
	}
}
