package usecase

import (
	"database/sql"
	"github.com/go-redis/redis"
	"loan-ms-go/internal/entity"
	"loan-ms-go/internal/repository"
)

type UserUseCase struct {
	*UseCase
	UserRepository *repository.UserRepository
}

func NewUserUseCase(cache *redis.Client, userRepository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		UseCase: NewUseCase(cache),
		UserRepository: userRepository,
	}
}

func (uc *UserUseCase) AddUser(user entity.User) (entity.User, error) {
	addedUser, err := uc.UserRepository.AddUser(user)
	if err != nil {
		return entity.User{}, err
	}

	return addedUser, nil
}

func (uc *UserUseCase) GetUsers() ([]entity.User, error) {
	users, err := uc.UserRepository.GetUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (uc *UserUseCase) GetUserByID(userID int) (entity.User, error) {
	user, err := uc.UserRepository.GetUserByID(userID)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (uc *UserUseCase) UpdateUser(user entity.User) (sql.Result, error) {
	result, err := uc.UserRepository.UpdateUser(user)
	if err != nil {
		return nil, err
	}
	return result, nil
}