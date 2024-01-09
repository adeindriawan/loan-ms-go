package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"loan-ms-go/internal/entity"
	"loan-ms-go/internal/usecase"
)

func AddUserHandler(uc *usecase.UserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		_, err = uc.AddUser(newUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to add user. Mesage: %s", err.Error())
			return

		}

		http.Redirect(w, r, "/users", http.StatusSeeOther)
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
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

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

		http.Redirect(w, r, "/users", http.StatusSeeOther)
	}
}

func GetUserByIDHandler(uc *usecase.UserUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userIDStr := vars["id"]
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		user, err := uc.GetUserByID(userID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "User with ID %s not found", userID)
			return
		}

		RenderHTMLTemplate(w, "user_details.html", user)
	}
}
