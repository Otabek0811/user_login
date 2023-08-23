package storage

import (
	"app/api/models"
	"context"
)

type StorageI interface {
	CloseDB()
	User() UserRepoI
	Phone() PhoneRepoI
}
type UserRepoI interface {
	Create(ctx context.Context, req *models.CreateUser) (string, error)
	GetByID(ctx context.Context, req *models.UserPrimaryKey) (*models.User, error)
	GetList(ctx context.Context, req *models.GetListUserRequest) (resp *models.GetListUserResponse, err error)
	Update(ctx context.Context, req *models.UpdateUser) (int64, error)
	Delete(ctx context.Context, req *models.UserPrimaryKey) (int64, error)
}

type PhoneRepoI interface {
	Create(ctx context.Context, req *models.CreatePhone) (string, error)
	GetByID(ctx context.Context, req *models.PhonePrimaryKey) (*models.Phone, error)
	GetList(ctx context.Context, req *models.GetListPhoneRequest) (resp *models.GetListPhoneResponse, err error)
	Update(ctx context.Context, req *models.UpdatePhone) (int64, error)
	Delete(ctx context.Context, req *models.PhonePrimaryKey) (int64, error)
}
