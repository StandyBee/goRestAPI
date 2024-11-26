package service

import (
	"gorestAPI"
	"gorestAPI/pkg/repository"
	"gorestAPI/pkg/request/item"
)

type TodoItemService struct {
	repo repository.TodoItem
}

func NewTodoItemService(r repository.TodoItem) *TodoItemService {
	return &TodoItemService{repo: r}
}

func (s *TodoItemService) GetListItems(listId int) ([]gorestAPI.TodoItem, error) {
	return s.repo.GetListItems(listId)
}

func (s *TodoItemService) GetItem(itemId int) (gorestAPI.TodoItem, error) {
	return s.repo.GetItem(itemId)
}

func (s *TodoItemService) CreateItem(req item.CreateItemRequest, listId int) (int, error) {
	return s.repo.CreateItem(req, listId)
}

func (s *TodoItemService) UpdateItem(req *item.UpdateItemRequest, listId, itemId int) error {
	return s.repo.UpdateItem(req, listId, itemId)
}

func (s *TodoItemService) DeleteItem(itemId int) error {
	return s.repo.DeleteItem(itemId)
}
