package repository

import (
	"database/sql"
)

type Repository[T any] struct {}

func (r *Repository[T]) FindById(db *sql.DB, entity *T) error {
	_, err := db.Exec("UPDATE")
}