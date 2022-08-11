package service

import (
	"errors"

	"github.com/gofrs/uuid"
	"github.com/somprasongd/go-monorepo/common"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/module/todo/core/dto"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/module/todo/core/mapper"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/module/todo/core/model"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/module/todo/core/ports"
	"gorm.io/gorm/logger"
)

var (
	// ErrTodoNotFoundById todo not found error when find with id
	ErrTodoNotFoundById = common.NewNotFoundError("todo with given id not found")
)

type todoService struct {
	repo ports.TodoRepository
}

func NewTodoService(repo ports.TodoRepository) ports.TodoService {
	return &todoService{repo}
}

func (s todoService) Create(userId string, form dto.NewTodoForm, log logger.Interface) (*dto.TodoResponse, error) {
	// validate
	if err := common.ValidateDto(form); err != nil {
		return nil, common.NewInvalidError(err.Error())
	}

	todo := model.Todo{
		Text:   form.Text,
		UserId: uuid.FromStringOrNil(userId),
	}
	// เรียกใช้ repo เพื่อบันทึกข้อมูลใหม่
	err := s.repo.Create(&todo)
	if err != nil {
		return nil, common.ErrDbInsert
	}

	// สร้าง struct ที่ต้องการให้ handler ส่งกลับไปหา client
	serialized := mapper.TodoToDto(&todo)

	return serialized, nil
}

func (s todoService) List(userId string, page common.PagingRequest, filters dto.ListTodoFilter, log logger.Interface) (dto.TodoResponses, *common.PagingResult, error) {
	// validate
	if err := common.ValidateDto(filters); err != nil {
		return nil, nil, common.NewInvalidError(err.Error())
	}

	todos, paging, err := s.repo.Find(userId, page, filters)
	if err != nil {
		log.Error(err.Error())
		return nil, nil, common.ErrDbQuery
	}

	serialized := mapper.TodosToDto(todos)
	return serialized, paging, nil
}

func (s todoService) Get(userId string, id string, log logger.Interface) (*dto.TodoResponse, error) {
	// validate id format
	if err := s.validateId(id); err != nil {
		return nil, common.ErrIdFormat
	}

	task, err := s.repo.FindById(id, userId)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, ErrTodoNotFoundById
		}
		log.Error(err.Error())
		return nil, common.ErrDbQuery
	}

	serialized := mapper.TodoToDto(task)
	return serialized, nil
}

func (s todoService) UpdateStatus(userId string, id string, form dto.UpdateTodoForm, log logger.Interface) (*dto.TodoResponse, error) {
	// validate id format
	if err := s.validateId(id); err != nil {
		return nil, common.ErrIdFormat
	}

	err := common.ValidateDto(form)
	if err != nil {
		return nil, common.NewInvalidError(err.Error())
	}

	todo, err := s.repo.UpdateStatusById(id, userId, form.Completed)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return nil, ErrTodoNotFoundById
		}
		log.Error(err.Error())
		return nil, common.ErrDbUpdate
	}

	serialized := mapper.TodoToDto(todo)
	return serialized, nil
}

func (s todoService) Delete(userId string, id string, log logger.Interface) error {
	// validate id format
	if err := s.validateId(id); err != nil {
		return common.ErrIdFormat
	}

	err := s.repo.DeleteById(id, userId)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return ErrTodoNotFoundById
		}
		log.Error(err.Error())
		return common.ErrDbDelete
	}

	return nil
}

func (s todoService) validateId(id string) error {
	return common.ValidateDto(&dto.TodoId{ID: id})
}
