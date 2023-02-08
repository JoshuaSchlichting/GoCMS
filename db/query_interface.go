package db

import (
	"context"

)

type QueriesInterface interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteUser(ctx context.Context, id int32) error
	GetUser(ctx context.Context, id int32) (User, error)
	ListFiles(ctx context.Context) ([]File, error)
	ListUsers(ctx context.Context) ([]User, error)
	UpdatUser(ctx context.Context, arg UpdatUserParams) (User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}
