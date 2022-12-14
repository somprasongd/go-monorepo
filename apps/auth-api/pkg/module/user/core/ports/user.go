package ports

import (
	"github.com/somprasongd/go-monorepo/common"
	"github.com/somprasongd/go-monorepo/common/logger"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/user/core/dto"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/user/core/model"
)

// interface สำหรับ output port
type UserRepository interface {
	Create(m *model.User) error
	Find(page common.PagingRequest) (model.Users, *common.PagingResult, error)
	FindById(id string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	UpdatePasswordById(id string, m *model.User) error
	DeleteById(id string) error
}

// interface สำหรับ input port
type UserService interface {
	Create(newUser dto.NewUserForm, log logger.Interface) (*dto.UserResponse, error)
	List(page common.PagingRequest, log logger.Interface) (dto.UserResponses, *common.PagingResult, error)
	Get(id string, log logger.Interface) (*dto.UserResponse, error)
	UpdatePassword(id string, updateUser dto.UpdateUserPasswordForm, log logger.Interface) (*dto.UserResponse, error)
	Delete(id string, log logger.Interface) error
}
