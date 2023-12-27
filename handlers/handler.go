package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"html/template"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	user := User{
		ID: 1,
		Name: "John Doe",
		Email: "johndoe@example.com",
	}

    // Parse the HTML template
    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // Render the HTML template
    err = tmpl.Execute(w, user)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}

func GetUserHandler(db *sql.DB, redisClient *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["id"]

		// Check if user data is in Redis cache
		cacheKey := fmt.Sprintf("user:%s", userID)
		userData, err := redisClient.Get(cacheKey).Result()
		if err == nil {
			// Data found in cache, return from cache
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(userData))
			return
		}

		// Data not in cache, fetch from MySQL
		query := "SELECT * FROM users WHERE id = ?"
		row := db.QueryRow(query, userID)

		var user User
		err = row.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "User with ID %s not found", userID)
			return
		}

		// Convert user data to JSON
		userJSON, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error converting user data to JSON")
			return
		}

		// Save user data to Redis cache
		err = redisClient.Set(cacheKey, userJSON, 0).Err()
		if err != nil {
			fmt.Println("Error saving to Redis cache:", err)
		}

		// Return user data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(userJSON)
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
		id := r.Form.Get("id")
		newName := r.Form.Get("name")
		newEmail := r.Form.Get("email")

		stmt := "UPDATE users SET name = ?, email = ? WHERE id = ?"
		_, err = db.Exec(stmt, newName, newEmail, id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Failed to update user. Mesage: %s", err.Error())
			return
		}

		// For demonstration purposes, let's print the updated data
		fmt.Printf("Updated Name: %s, Updated Email: %s, for ID: %s\n", newName, newEmail, id)

		// You can perform database updates or any other business logic here

		// Redirect to the home page after updating
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
