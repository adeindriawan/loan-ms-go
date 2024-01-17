// user_repository_test.go

package repository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"loan-ms-go/internal/entity"
)

func TestGetUserByID(t *testing.T) {
	// Create a new mock DB and repository
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer mockDB.Close()

	repo := &UserRepository{
		Repository: NewRepository(mockDB),
	}

	// Define the test cases
	tests := []struct {
		name           string
		userID         int
		expectedUser   entity.User
		expectedError  error
		mockRows       *sqlmock.Rows
	}{
		{
			name:         "User found",
			userID:       1,
			expectedUser: entity.User{ID: 1, Name: "John Doe", Email: "john@example.com"},
			expectedError: nil,
			mockRows: sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(1, "John Doe", "john@example.com"),
		},
		{
			name:         "User not found",
			userID:       2,
			expectedUser: entity.User{},
			expectedError: errors.New("sql: no rows in result set"),
			mockRows:     sqlmock.NewRows([]string{"id", "name", "email"}),
		},
		// Add more test cases as needed
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set up mock expectations
			mock.ExpectQuery("SELECT \\* FROM users WHERE id =\\?").
				WithArgs(test.userID).
				WillReturnRows(test.mockRows)

			// Call the method being tested
			user, err := repo.GetUserByID(test.userID)

			// Check the results
			if err != nil && err.Error() != test.expectedError.Error() {
				t.Errorf("Unexpected error. Got: %v, Expected: %v", err, test.expectedError)
			}

			if user != test.expectedUser {
				t.Errorf("Unexpected user. Got: %v, Expected: %v", user, test.expectedUser)
			}

			// Ensure all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
