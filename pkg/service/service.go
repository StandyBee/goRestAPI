package service

import (
	"gorestAPI"
	"gorestAPI/pkg/repository"
	"gorestAPI/pkg/request"
	"gorestAPI/pkg/request/item"
)

type Authorization interface {
	CreateUser(user gorestAPI.User) (int, error)
	Login(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	GetUserLists(userId int) ([]gorestAPI.Todo, error)
	GetListById(listId, userId int) (gorestAPI.Todo, error)
	CreateList(title, description string, userId int) (int, error)
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

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem),
	}
}
