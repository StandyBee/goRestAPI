package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gorestAPI"
	"gorestAPI/pkg/request/item"
)

type ItemPostgres struct {
	db *sqlx.DB
}

func NewItemPostgres(db *sqlx.DB) *ItemPostgres {
	return &ItemPostgres{db: db}
}

func (r *ItemPostgres) GetListItems(listId int) ([]gorestAPI.TodoItem, error) {
	var items []gorestAPI.TodoItem

	query := `
        SELECT ti.id, ti.title, ti.description, ti.done
        FROM todo_items ti 
        INNER JOIN lists_items li ON li.todo_item_id = ti.id
        WHERE li.todo_list_id = $1
    `

	err := r.db.Select(&items, query, listId)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ItemPostgres) GetItem(itemId int) (gorestAPI.TodoItem, error) {
	var res gorestAPI.TodoItem
	query := `
		SELECT * FROM todo_items WHERE id=$1
	`
	err := r.db.Get(&res, query, itemId)
	return res, err
}

func (r *ItemPostgres) CreateItem(req item.CreateItemRequest, listId int) (int, error) {
	var itemId int

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO todo_items (title, description, done) values ($1, $2, $3) RETURNING id`
	err = tx.QueryRow(query, req.Title, req.Description, req.Done).Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	query = `INSERT INTO lists_items (todo_item_id, todo_list_id) VALUES ($1, $2)`
	_, err = tx.Exec(query, itemId, listId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return itemId, nil
}

func (r *ItemPostgres) UpdateItem(req *item.UpdateItemRequest, listId, itemId int) error {
	checkQuery := `SELECT COUNT(*) FROM lists_items WHERE todo_item_id = $1 AND todo_list_id = $2`
	var count int
	err := r.db.QueryRow(checkQuery, itemId, listId).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check item ownership: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("item does not belong to the list")
	}

	query := `UPDATE todo_items SET`
	args := []interface{}{}
	argCounter := 1

	if req.Title != nil {
		query += fmt.Sprintf(" title = $%d,", argCounter)
		args = append(args, *req.Title)
		argCounter++
	}

	if req.Description != nil {
		query += fmt.Sprintf(" description = $%d,", argCounter)
		args = append(args, *req.Description)
		argCounter++
	}

	if req.Done != nil {
		query += fmt.Sprintf(" done = $%d,", argCounter)
		args = append(args, *req.Done)
		argCounter++
	}

	// Удаляем последнюю запятую и добавляем WHERE
	query = query[:len(query)-1] + fmt.Sprintf(" WHERE id = $%d", argCounter)
	args = append(args, itemId)

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

func (r *ItemPostgres) DeleteItem(itemId int) error {
	query := `DELETE FROM todo_items WHERE id=$1`
	_, err := r.db.Exec(query, itemId)
	if err != nil {
		return err
	}
	return nil
}
