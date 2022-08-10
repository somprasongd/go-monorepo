package repository

import (
	"errors"

	"github.com/somprasongd/go-monorepo/services/auth/pkg/common"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/auth/core/ports"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/user/core/model"
	"gorm.io/gorm"
)

type authRepositoryDB struct {
	db *gorm.DB
}

func NewAuthRepositoryDB(db *gorm.DB) ports.AuthRepository {
	return &authRepositoryDB{db}
}

func (r authRepositoryDB) FindUserByEmail(email string) (*model.User, error) {
	user := model.User{}
	db := r.db.Where("email = ?", email).First(&user)
	if err := db.Error; err != nil {
		// handle error not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r authRepositoryDB) CreateUser(user *model.User) error {
	return r.db.Create(&user).Error
}

func (r authRepositoryDB) SaveProfile(user *model.User) error {
	return r.db.Save(&user).Error
}
