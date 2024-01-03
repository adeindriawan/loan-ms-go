package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"loan-ms-go/internal/entity"
	"loan-ms-go/internal/usecase"
)

func AddUserHandler(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
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

		user, err := usecase.AddUser(db, newUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to add user. Mesage: %s", err.Error())
			return

		}

		saveUserToCache(redisClient, &user)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func GetUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := usecase.GetUsers(db)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Failed to get users data. Mesage: %s", err.Error())
			return
		}

		RenderHTMLTemplate(w, "users.html", users)
	}
}

func UpdateUserHandler(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
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

		_, err = usecase.UpdateUser(db, updatedUser)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Failed to update user. Mesage: %s", err.Error())
			return
		}

		saveUserToCache(redisClient, &updatedUser)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func GetUserHandler(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userIDStr := vars["id"]
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Check if user data is in Redis cache
		cachedUser, err := getUserFromCache(redisClient, userIDStr)
		if err == nil {
			RenderHTMLTemplate(w, "user_details.html", cachedUser)
			return
		}

		// Data not in cache, fetch from MySQL
		user, err := usecase.GetUser(db, userID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "User with ID %s not found", userID)
			return
		}

		saveUserToCache(redisClient, &user)

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
