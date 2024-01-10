package usecase

import (
	"fmt"
	"encoding/json"
	"strconv"
	"database/sql"
	"github.com/go-redis/redis"
	"loan-ms-go/services"
	"loan-ms-go/internal/entity"
	"loan-ms-go/internal/repository"
)

type UserUseCase struct {
	*UseCase
	UserRepository *repository.UserRepository
}

func NewUserUseCase(cache *redis.Client, logger *services.Logger, userRepository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		UseCase: NewUseCase(cache, logger),
		UserRepository: userRepository,
	}
}

func (uc *UserUseCase) AddUser(user entity.User) (entity.User, error) {
	addedUser, err := uc.UserRepository.AddUser(user)
	if err != nil {
		return entity.User{}, err
	}

	uc.saveUserToCache(&addedUser)

	return addedUser, nil
}

func (uc *UserUseCase) GetUsers() ([]entity.User, error) {
	users, err := uc.UserRepository.GetUsers()
	if err != nil {
		return nil, err
	}

	uc.Logger.InfoLogger.Println("Users data fetched.")

	return users, nil
}

func (uc *UserUseCase) GetUserByID(userID int) (entity.User, error) {
	cachedUser, err := uc.getUserFromCache(strconv.Itoa(userID))
	if err == nil {
		return *cachedUser, nil
	}
	
	user, err := uc.UserRepository.GetUserByID(userID)
	if err != nil {
		return entity.User{}, err
	}

	uc.saveUserToCache(&user)

	return user, nil
}

func (uc *UserUseCase) UpdateUser(user entity.User) (sql.Result, error) {
	result, err := uc.UserRepository.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	uc.saveUserToCache(&user)
	return result, nil
}

func (uc *UserUseCase) getUserFromCache(userIDStr string) (*entity.User, error) {
	cacheKey := fmt.Sprintf("user:%s", userIDStr)
	cacheData, err := uc.Cache.Get(cacheKey).Result()

	if err != nil {
		return nil, err
	}

	var cachedUser entity.User
	err = json.Unmarshal([]byte(cacheData), &cachedUser)
	if err != nil {
		return nil, err
	}

	return &cachedUser, nil
}

func (uc *UserUseCase) saveUserToCache(user *entity.User) {
	userIDStr := strconv.Itoa(user.ID)
	cacheKey := fmt.Sprintf("user:%s", userIDStr)
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error converting user data to JSON:", err)
		return
	}

	err = uc.Cache.Set(cacheKey, userJSON, 0).Err()
	if err != nil {
		fmt.Println("Error saving to Redis cache:", err)
		return
	}
}