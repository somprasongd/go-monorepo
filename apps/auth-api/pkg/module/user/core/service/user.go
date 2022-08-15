package service

import (
	"errors"
	"fmt"

	"github.com/somprasongd/go-monorepo/common"
	"github.com/somprasongd/go-monorepo/common/logger"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/user/core/dto"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/user/core/mapper"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/user/core/ports"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/util"
)

var (
	// ErrUserNotFoundById user not found error when find with id
	ErrUserNotFoundById     = common.NewNotFoundError("user with given id not found")
	ErrHashPassword         = common.NewUnexpectedError("hash password error")
	ErrUserEmailDuplication = common.NewBadRequestError("duplicate email")
)

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &userService{repo}
}

func (s userService) Create(form dto.NewUserForm, log logger.Interface) (*dto.UserResponse, error) {
	// validate
	if err := common.ValidateDto(form); err != nil {
		return nil, common.NewInvalidError(err.Error())
	}

	u, err := s.repo.FindByEmail(form.Email)

	if err != nil && !errors.Is(err, common.ErrRecordNotFound) {
		log.Error(err.Error())
		return nil, common.ErrDbQuery
	}

	if u != nil {
		return nil, ErrUserEmailDuplication
	}

	user := mapper.CreateUserFormToModel(form)
	hashPwd, err := util.HashPassword(form.Password)
	if err != nil {
		log.Error(err.Error())
		return nil, ErrHashPassword
	}
	fmt.Println(form.Password, hashPwd)
	user.Password = hashPwd

	err = s.repo.Create(user)
	if err != nil {
		log.Error(err.Error())
		return nil, common.ErrDbInsert
	}

	serialized := mapper.UserToDto(user)

	return serialized, nil
}

func (s userService) List(page common.PagingRequest, log logger.Interface) (dto.UserResponses, *common.PagingResult, error) {
	users, paging, err := s.repo.Find(page)
	if err != nil {
		log.Error(err.Error())
		return nil, nil, common.ErrDbQuery
	}

	serialized := mapper.UsersToDto(users)
	return serialized, paging, nil
}

func (s userService) Get(id string, log logger.Interface) (*dto.UserResponse, error) {
	// validate id format
	if err := s.validateId(id); err != nil {
		return nil, common.ErrIdFormat
	}

	user, err := s.repo.FindById(id)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, ErrUserNotFoundById
		}
		log.Error(err.Error())
		return nil, common.ErrDbQuery
	}

	serialized := mapper.UserToDto(user)
	return serialized, nil
}

func (s userService) UpdatePassword(id string, form dto.UpdateUserPasswordForm, log logger.Interface) (*dto.UserResponse, error) {
	// validate id format
	if err := s.validateId(id); err != nil {
		return nil, common.ErrIdFormat
	}

	err := common.ValidateDto(form)
	if err != nil {
		return nil, common.NewInvalidError(err.Error())
	}

	user, err := s.repo.FindById(id)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, ErrUserNotFoundById
		}
		log.Error(err.Error())
		return nil, common.ErrDbQuery
	}

	hashPwd, err := util.HashPassword(form.Password)

	if err != nil {
		log.Error(err.Error())
		return nil, ErrHashPassword
	}

	user.Password = hashPwd

	err = s.repo.UpdatePasswordById(id, user)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, ErrUserNotFoundById
		}
		log.Error(err.Error())
		return nil, common.ErrDbUpdate
	}

	serialized := mapper.UserToDto(user)
	return serialized, nil
}

func (s userService) Delete(id string, log logger.Interface) error {
	// validate id format
	if err := s.validateId(id); err != nil {
		return common.ErrIdFormat
	}

	err := s.repo.DeleteById(id)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return ErrUserNotFoundById
		}
		log.Error(err.Error())
		return common.ErrDbDelete
	}

	return nil
}

func (s userService) validateId(id string) error {
	return common.ValidateDto(&dto.UserId{ID: id})
}
