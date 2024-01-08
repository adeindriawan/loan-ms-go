package repository

import (
	"database/sql"
	"loan-ms-go/internal/entity"
)

type UserRepository struct {
	*Repository
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		Repository: NewRepository(db),
	}
}

func (r *UserRepository) GetUsers() ([]entity.User, error) {
	rows, err := r.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

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

func GetUsers(db *sql.DB) ([]entity.User, error) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByID(db *sql.DB, userID int) (entity.User, error) {
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
