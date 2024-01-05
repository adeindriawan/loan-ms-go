package usecase

import (
	"database/sql"
	"loan-ms-go/internal/entity"
	"loan-ms-go/internal/repository"
)

func AddUser(db *sql.DB, user entity.User) (entity.User, error) {
	addedUser, err := repository.AddUser(db, user)
	if err != nil {
		return entity.User{}, err
	}

	return addedUser, nil
}

func GetUsers(db *sql.DB) ([]entity.User, error) {
	users, err := repository.GetUsers(db)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetUser(db *sql.DB, userID int) (entity.User, error) {
	user, err := repository.GetUserByID(db, userID)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func UpdateUser(db *sql.DB, user entity.User) (sql.Result, error) {
	result, err := repository.UpdateUser(db, user)
	if err != nil {
		return nil, err
	}
	return result, nil
}