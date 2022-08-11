package ports

import (
	"github.com/somprasongd/go-monorepo/common"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/module/todo/core/dto"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/module/todo/core/model"
)

// interface สำหรับ output port
type TodoRepository interface {
	Create(m *model.Todo) error
	Find(userId string, page common.PagingRequest, filters dto.ListTodoFilter) (model.Todos, *common.PagingResult, error)
	FindById(id string, userId string) (*model.Todo, error)
	UpdateStatusById(id string, userId string, status bool) (*model.Todo, error)
	DeleteById(id string, userId string) error
}

// interface สำหรับ input port
type TodoService interface {
	Create(userId string, newTodo dto.NewTodoForm, reqId string) (*dto.TodoResponse, error)
	List(userId string, page common.PagingRequest, filters dto.ListTodoFilter, reqId string) (dto.TodoResponses, *common.PagingResult, error)
	Get(userId string, id string, reqId string) (*dto.TodoResponse, error)
	UpdateStatus(userId string, id string, updateTodo dto.UpdateTodoForm, reqId string) (*dto.TodoResponse, error)
	Delete(userId string, id string, reqId string) error
}
