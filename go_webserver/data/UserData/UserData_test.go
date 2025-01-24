package UserData_test

import (
	"bytes"
	"go_playground/go_webserver/data/UserData"
	"go_playground/go_webserver/types"
	"io"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() *UserData.UserDataImpl {
	// Reset userCache for each test
	UserCache := &types.UserCache{
		SafeMap: &sync.Map{},
		Count:   0,
	}
	return &UserData.UserDataImpl{UserCache: *UserCache}
}

func TestCreateUser(t *testing.T) {
	userDataService := setup()

	t.Run("Successfully create a user", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"FirstName": "John", "LastName": "Doe", "Email": "john.doe@example.com"}`))
		user, err := userDataService.CreateUser(io.NopCloser(body))

		assert.NoError(t, err)
		assert.Equal(t, 1, user.Id)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, "john.doe@example.com", user.Email)

		value, _ := userDataService.UserCache.SafeMap.Load(1)
		assert.Equal(t, user, value)
	})

	t.Run("Missing required fields", func(t *testing.T) {
		body := bytes.NewReader([]byte(`{"FirstName": "", "LastName": "", "Email": ""}`))
		_, err := userDataService.CreateUser(io.NopCloser(body))

		assert.EqualError(t, err, "first name is required\nlast name is required\nEmail is required\n")
	})
}

func TestGetUser(t *testing.T) {
	userDataService := setup()

	// Add a user to the cache
	testUser := types.User{Id: 1, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}
	userDataService.UserCache.SafeMap.Store(1, testUser)

	t.Run("Get existing user", func(t *testing.T) {
		user, err := userDataService.GetUser(1)

		assert.NoError(t, err)
		assert.Equal(t, testUser, user)
	})

	t.Run("Get non-existing user", func(t *testing.T) {
		_, err := userDataService.GetUser(2)

		assert.EqualError(t, err, "user not found")
	})
}

func TestDoesUserExist(t *testing.T) {
	userDataService := setup()

	userDataService.UserCache.SafeMap.Store(1, types.User{Id: 1})

	t.Run("User exists", func(t *testing.T) {
		exists := userDataService.DoesUserExist(1)

		assert.True(t, exists)
	})

	t.Run("User does not exist", func(t *testing.T) {
		exists := userDataService.DoesUserExist(2)

		assert.False(t, exists)
	})
}

func TestUpdateUser(t *testing.T) {
	userDataService := setup()

	testUser := types.User{Id: 1, FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}
	userDataService.UserCache.SafeMap.Store(1, testUser)

	t.Run("Update existing user", func(t *testing.T) {
		updatedUser := types.User{Id: 1, FirstName: "John", LastName: "Smith", Email: "john.smith@example.com"}
		err := userDataService.UpdateUser(updatedUser)

		assert.NoError(t, err)

		value, _ := userDataService.UserCache.SafeMap.Load(1)
		assert.Equal(t, updatedUser, value)
	})

	t.Run("Update non-existing user", func(t *testing.T) {
		nonExistentUser := types.User{Id: 2, FirstName: "Jane"}
		err := userDataService.UpdateUser(nonExistentUser)

		assert.EqualError(t, err, "user not found")
	})
}

func TestDeleteUser(t *testing.T) {
	userDataService := setup()

	userDataService.UserCache.SafeMap.Store(1, types.User{Id: 1})

	t.Run("Delete existing user", func(t *testing.T) {
		err := userDataService.DeleteUser(1)

		assert.NoError(t, err)

		_, ok := userDataService.UserCache.SafeMap.Load(1)
		assert.False(t, ok)
	})

	t.Run("Delete non-existing user", func(t *testing.T) {
		err := userDataService.DeleteUser(2)

		assert.EqualError(t, err, "user not found")
	})
}
