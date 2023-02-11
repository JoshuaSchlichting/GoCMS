package db

import (
	"context"
)

type QueriesInterface interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteUser(ctx context.Context, id int64) error
	GetUser(ctx context.Context, id int64) (User, error)
	ListFiles(ctx context.Context) ([]File, error)
	ListUsers(ctx context.Context) ([]User, error)
	UpdatUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
	UploadFile(ctx context.Context, arg UploadFileParams) (File, error)
	GetUserByName(ctx context.Context, name string) (User, error)
	UpdateOrganization(ctx context.Context, arg UpdateOrganizationParams) (Organization, error)
}
