package repository

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"gorestAPI"
	"gorestAPI/pkg/request"
)

type ListPostgres struct {
	db *sqlx.DB
}

func NewListPostgres(db *sqlx.DB) *ListPostgres {
	return &ListPostgres{db: db}
}

func (r *ListPostgres) GetUserLists(userId int) ([]gorestAPI.Todo, error) {
	var lists []gorestAPI.Todo
	query := `
        SELECT tl.id, tl.title, tl.description
        FROM todo_lists tl
        INNER JOIN users_lists ul ON ul.todo_list_id = tl.id
        WHERE ul.user_id = $1
    `
	err := r.db.Select(&lists, query, userId)
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func (r *ListPostgres) GetListById(listId, userId int) (gorestAPI.Todo, error) {
	var list gorestAPI.Todo
	query := `
		SELECT tl.id, tl.title, tl.description
		FROM todo_lists tl
		INNER JOIN users_lists ul ON ul.todo_list_id = tl.id
		WHERE ul.todo_list_id = $1
		AND ul.user_id = $2
	`
	err := r.db.Get(&list, query, listId, userId)
	if err != nil {
		return list, err
	}
	return list, nil
}

func (r *ListPostgres) CreateList(title, description string, userId int) (int, error) {
	var listId int

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	defer tx.Rollback()

	query := `INSERT INTO todo_lists (title, description) VALUES ($1, $2) RETURNING id`
	err = tx.QueryRow(query, title, description).Scan(&listId)
	if err != nil {
		return 0, fmt.Errorf("failed to insert into todo_lists: %v", err)
	}

	query = `INSERT INTO users_lists (user_id, todo_list_id) VALUES ($1, $2)`
	_, err = tx.Exec(query, userId, listId)
	if err != nil {
		return 0, fmt.Errorf("failed to insert into users_lists: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return listId, nil
}

func (r *ListPostgres) UpdateList(req *reqs.UpdateListRequest) error {
	query := `
		UPDATE todo_lists
		SET title = $1, description = $2
		WHERE id = $3
		AND EXISTS (
			SELECT 1 
			FROM users_lists
			WHERE users_lists.todo_list_id = $3
			AND users_lists.user_id = $4
		)
		RETURNING id
	`

	var updatedId int
	err := r.db.QueryRow(query, req.Title, req.Description, req.ListId, req.UserId).Scan(&updatedId)
	if err != nil {
		// Если ошибка, возвращаем ошибку
		return err
	}
	if updatedId == 0 {
		// Если id не найдено, значит ничего не обновилось
		return fmt.Errorf("no rows updated, either list or user association does not exist")
	}
	return nil
}

func (r *ListPostgres) DeleteList(listId, userId int) error {
	query := `
		DELETE FROM todo_lists 
		WHERE id = $1
		AND EXISTS (
		    SELECT 1
		    FROM users_lists
		    WHERE users_lists.todo_list_id = todo_lists.id
		    AND users_lists.user_id = $2
		)
	`
	result, err := r.db.Exec(query, listId, userId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("list not found or access denied")
	}

	return nil
}
