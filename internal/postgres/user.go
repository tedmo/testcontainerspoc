package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/tedmo/testcontainerspoc/internal/postgres/sqlc"
)
import "github.com/tedmo/testcontainerspoc/internal/app"

type UserRepo struct {
	DB      *sql.DB
	Querier sqlc.Querier
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		DB:      db,
		Querier: sqlc.New(),
	}
}

func (repo *UserRepo) FindUserByID(ctx context.Context, id int64) (*app.User, error) {
	sqlUser, err := repo.Querier.FindUserByID(ctx, repo.DB, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, app.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return sqlUser.DomainModel(), nil
}

func (repo *UserRepo) CreateUser(ctx context.Context, user *app.CreateUserReq) (*app.User, error) {

	sqlUser, err := repo.Querier.CreateUser(ctx, repo.DB, user.Name)
	if err != nil {
		return nil, err
	}

	return sqlUser.DomainModel(), nil
}

func (repo *UserRepo) FindUsers(ctx context.Context) ([]app.User, error) {

	sqlUsers, err := repo.Querier.FindUsers(ctx, repo.DB)
	if err != nil {
		return nil, err
	}

	var users []app.User
	for _, sqlUser := range sqlUsers {
		users = append(users, *sqlUser.DomainModel())
	}

	if users == nil {
		users = []app.User{}
	}

	return users, nil
}
