// user_repository_test.go

package repository

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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
			name:         	"User found",
			userID:       	1,
			expectedUser: 	entity.User{ID: 1, Name: "John Doe", Email: "john@example.com"},
			expectedError: 	nil,
			mockRows: 		sqlmock.NewRows([]string{"id", "name", "email"}).
							AddRow(1, "John Doe", "john@example.com"),
		},
		{
			name:         	"User not found",
			userID:       	2,
			expectedUser: 	entity.User{},
			expectedError: 	errors.New("sql: no rows in result set"),
			mockRows:     	sqlmock.NewRows([]string{"id", "name", "email"}),
		},
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

func TestGetUsers(t *testing.T) {
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
		expectedUsers []entity.User
		expectedError  error
		mockRows       *sqlmock.Rows
	}{
		{
			name:         	"User found",
			expectedUsers: 	[]entity.User{
				{ID: 1, Name: "John Doe", Email: "john@example.com"},
			},
			expectedError:	nil,
			mockRows: 		sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(1, "John Doe", "john@example.com"),
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set up mock expectations
			mock.ExpectQuery("SELECT id, name, email FROM users").
				WillReturnRows(test.mockRows)

			// Call the method being tested
			users, err := repo.GetUsers()

			// Check the results using testify's assert
			assert.Equal(t, test.expectedError, err)
			assert.ElementsMatch(t, test.expectedUsers, users)

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAddUser(t *testing.T) {
	// Create a new mock DB
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer mockDB.Close()

	// Create a new UserRepository with the mock DB
	repo := &UserRepository{
		Repository: NewRepository(mockDB),
	}

	// Define the test case
	testUser := entity.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mock.ExpectExec("INSERT INTO users").
		WithArgs(testUser.Name, testUser.Email).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the method being tested
	resultUser, err := repo.AddUser(testUser)

	// Check the results
	assert.NoError(t, err)
	assert.Equal(t, 1, resultUser.ID) // Assuming the ID is set to 1 in the mock result

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser(t *testing.T) {
	// Setup Mock Database
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock DB: %v", err)
	}
	defer mockDB.Close()

	// Create UserRepository with Mock Database
	repo := &UserRepository{
		Repository: NewRepository(mockDB),
	}

	// Define Test User and Expectations
	testUser := entity.User{
		ID:    1,
		Name:  "Updated Name",
		Email: "updated@example.com",
	}
	mock.ExpectExec("UPDATE users").
		WithArgs(testUser.Name, testUser.Email, testUser.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Call the Method Being Tested
	result, err := repo.UpdateUser(testUser)

	// Check the Results
	assert.NoError(t, err)
	affectedRows, _ := result.RowsAffected()
	assert.Equal(t, int64(1), affectedRows)

	// Ensure Expectations Were Met
	assert.NoError(t, mock.ExpectationsWereMet())
}
