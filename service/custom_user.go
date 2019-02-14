package service

import (
	"github.com/jinzhu/gorm"
	"github.com/teamlint/gen/model"
)

type customUserService struct {
	userService
}

func NewCustomUserService(db *gorm.DB) CustomUserService {
	svc := customUserService{userService: userService{db}}
	return &svc
}
func (s *customUserService) GetByName(name string) (*model.User, error) {
	var item model.User

	if err := s.DB.Where("username=?", name).Take(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.ErrRecordNotFound
		}
		return nil, err
	}

	return &item, nil
}
