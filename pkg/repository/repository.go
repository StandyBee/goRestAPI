package repository

import (
	"github.com/jmoiron/sqlx"
	"gorestAPI"
	"gorestAPI/pkg/request"
	"gorestAPI/pkg/request/item"
)

type Authorization interface {
	CreateUser(user gorestAPI.User) (int, error)
	GetUser(username, password string) (gorestAPI.User, error)
}

type TodoList interface {
	GetUserLists(userId int) ([]gorestAPI.Todo, error)
	GetListById(listId, userId int) (gorestAPI.Todo, error)
	CreateList(title, desc string, userId int) (int, error)
	UpdateList(req *reqs.UpdateListRequest) error
	DeleteList(listId, userId int) error
}

type TodoItem interface {
	GetListItems(listId int) ([]gorestAPI.TodoItem, error)
	GetItem(itemId int) (gorestAPI.TodoItem, error)
	CreateItem(req item.CreateItemRequest, listId int) (int, error)
	UpdateItem(req *item.UpdateItemRequest, listId, itemId int) error
	DeleteItem(itemId int) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {

	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewListPostgres(db),
		TodoItem:      NewItemPostgres(db),
	}
}
