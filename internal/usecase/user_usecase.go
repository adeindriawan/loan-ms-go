package usecase

import (
	"database/sql"
	"loan-ms-go/internal/entity"
)

func AddUser(db *sql.DB, user entity.User) (entity.User, error) {
	query := "INSERT INTO users (name, email) VALUES (?, ?)"
	result, err := db.Exec(query, user.Name, user.Email)

	if err != nil {
		return entity.User{}, err
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return entity.User{}, err
	}

	user.ID = int(lastInsertedID)
	return user, nil
}

func GetUser(db *sql.DB, userID int) (entity.User, error) {
	query := "SELECT * FROM users WHERE id = ?"
	row := db.QueryRow(query, userID)

	var user entity.User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func UpdateUser(db *sql.DB, user entity.User) (sql.Result, error) {
	stmt := "UPDATE users SET name = ?, email = ? WHERE id = ?"
	return db.Exec(stmt, user.Name, user.Email, user.ID)
}