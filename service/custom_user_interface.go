package service

import "github.com/teamlint/gen/model"

type CustomUserService interface {
	UserService
	// Log LogService
	// Shop ShopService
	GetByName(name string) (*model.User, error)
}
