package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"loan-ms-go/internal/entity"
	"loan-ms-go/internal/usecase"
)

func AddUserHandler(uc *usecase.UserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		newName := r.Form.Get("name")
		newEmail := r.Form.Get("email")
		newUser := entity.User{
			Name: newName,
			Email: newEmail,
		}

		user, err := uc.AddUser(newUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to add user. Mesage: %s", err.Error())
			return

		}

		saveUserToCache(uc.Cache, &user)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func GetUsersHandler(uc *usecase.UserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := uc.GetUsers()
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Failed to get users data. Mesage: %s", err.Error())
			return
		}

		RenderHTMLTemplate(w, "users.html", users)
	}
}

func UpdateUserHandler(uc *usecase.UserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Get the updated name and email from the form
		idStr := r.Form.Get("id")
		newName := r.Form.Get("name")
		newEmail := r.Form.Get("email")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		updatedUser := entity.User{
			ID: id,
			Name: newName,
			Email: newEmail,
		}

		_, err = uc.UpdateUser(updatedUser)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Failed to update user. Mesage: %s", err.Error())
			return
		}

		saveUserToCache(uc.Cache, &updatedUser)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func GetUserHandler(uc *usecase.UserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userIDStr := vars["id"]
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Check if user data is in cache
		cachedUser, err := getUserFromCache(uc.Cache, userIDStr)
		if err == nil {
			RenderHTMLTemplate(w, "user_details.html", cachedUser)
			return
		}

		// Data not in cache, fetch from db
		user, err := uc.GetUser(userID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "User with ID %s not found", userID)
			return
		}

		saveUserToCache(uc.Cache, &user)

		RenderHTMLTemplate(w, "user_details.html", user)
	}
}

func getUserFromCache(redisClient *redis.Client, userIDStr string) (*entity.User, error) {
	cacheKey := fmt.Sprintf("user:%s", userIDStr)
	cacheData, err := redisClient.Get(cacheKey).Result()

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

func saveUserToCache(redisClient *redis.Client, user *entity.User) {
	userIDStr := strconv.Itoa(user.ID)
	cacheKey := fmt.Sprintf("user:%s", userIDStr)
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error converting user data to JSON:", err)
		return
	}

	err = redisClient.Set(cacheKey, userJSON, 0).Err()
	if err != nil {
		fmt.Println("Error saving to Redis cache:", err)
		return
	}
}
