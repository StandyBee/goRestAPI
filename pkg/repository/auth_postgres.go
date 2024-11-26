package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"gorestAPI"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user gorestAPI.User) (int, error) {
	var userId int
	err := r.db.QueryRow("INSERT INTO users (name, username, password) values ($1, $2, $3) RETURNING id", user.Name, user.Username, user.Password).
		Scan(&userId)

	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (r *AuthPostgres) GetUser(username string, password string) (gorestAPI.User, error) {
	var user gorestAPI.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username = $1 AND password = $2", usersTable)
	err := r.db.Get(&user, query, username, password)
	if err != nil {
		return gorestAPI.User{}, err
	}
	return user, nil
}
