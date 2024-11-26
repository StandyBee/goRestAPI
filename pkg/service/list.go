package service

import (
	"gorestAPI"
	"gorestAPI/pkg/repository"
	"gorestAPI/pkg/request"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) GetUserLists(userId int) ([]gorestAPI.Todo, error) {
	return s.repo.GetUserLists(userId)
}

func (s *TodoListService) GetListById(listId, userId int) (gorestAPI.Todo, error) {
	return s.repo.GetListById(listId, userId)
}

func (s *TodoListService) CreateList(title, desc string, userId int) (int, error) {
	return s.repo.CreateList(title, desc, userId)
}

func (s *TodoListService) UpdateList(req *reqs.UpdateListRequest) error {
	return s.repo.UpdateList(req)
}

func (s *TodoListService) DeleteList(listId, userId int) error {
	return s.repo.DeleteList(listId, userId)
}
