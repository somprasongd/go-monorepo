package repository

import (
	"errors"

	"github.com/somprasongd/go-monorepo/common"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/module/todo/core/dto"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/module/todo/core/model"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/module/todo/core/ports"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type todoRepositoryDB struct {
	db *gorm.DB
}

func NewTodoRepositoryDB(db *gorm.DB) ports.TodoRepository {
	return &todoRepositoryDB{db}
}

func (r todoRepositoryDB) Create(todo *model.Todo) error {
	return r.db.Create(&todo).Error
}

func (r todoRepositoryDB) Find(userId string, page common.PagingRequest, filters dto.ListTodoFilter) (model.Todos, *common.PagingResult, error) {
	todos := []*model.Todo{}
	db := r.db

	db = db.Where(`user_id = ?`, userId)

	if len(filters.Term) > 0 {
		db = db.Where(`text ilike ?`,
			"%"+filters.Term+"%")
	}
	if filters.Completed != nil {
		status := model.OPEN.String()
		if *filters.Completed {
			status = model.DONE.String()
		}
		db = db.Where(`status = ?`, status)
	}
	pg := common.Pagination{
		PagingRequest: page,
		Query:         db,
		Records:       &todos,
	}
	paging, err := pg.Paginate()

	if err != nil {
		return nil, nil, err
	}

	return todos, paging, nil
}

func (r todoRepositoryDB) FindById(id string, userId string) (*model.Todo, error) {
	todo := model.Todo{}
	db := r.db.Where("id = ? and user_id = ?", id, userId).First(&todo)
	if err := db.Error; err != nil {
		// handle error not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, err
	}
	return &todo, nil
}

func (r todoRepositoryDB) UpdateStatusById(id string, userId string, status bool) (*model.Todo, error) {
	todo := model.Todo{}
	if status {
		todo.Done()
	} else {
		todo.Open()
	}
	db := r.db.Model(&todo).
		Clauses(clause.Returning{}).
		Where("id = ? and user_id = ?", id, userId).
		Updates(&todo)
	if err := db.Error; err != nil {
		return nil, err
	}
	// handle not found error
	if db.RowsAffected == 0 {
		return nil, common.ErrRecordNotFound
	}
	return &todo, nil
}

func (r todoRepositoryDB) DeleteById(id string, userId string) error {
	db := r.db.Where("id = ? and user_id = ?", id, userId).Delete(&model.Todo{})
	if err := db.Error; err != nil {
		return err
	}
	// handle not found error
	if db.RowsAffected == 0 {
		return common.ErrRecordNotFound
	}
	return nil
}
