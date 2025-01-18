package data

import (
	"encoding/json"
	"errors"
	"go_playground/go_webserver/types"
	"io"
	"sync"
)

var userCache = make(map[int]types.User)

var userCacheMutex sync.RWMutex

func GetUser(id int) (types.User, error) {
	userCacheMutex.RLock()
	user, ok := userCache[id]
	userCacheMutex.RUnlock()

	if !ok {
		return types.User{}, errors.New("user not found")
	}

	return user, nil
}

func CreateUser(body io.ReadCloser) (int, error) {
	var user types.User
	var errorString string
	err := json.NewDecoder(body).Decode(&user)
	if err != nil {
		return 0, err
	}

	if user.FirstName == "" {
		errorString += "first name is required\n"
	}
	if user.LastName == "" {
		errorString += "last name is required\n"
	}
	if user.Email == "" {
		errorString += "Email is required\n"
	}

	if errorString != "" {
		return 0, errors.New(errorString)
	}

	user.Id = len(userCache) + 1

	userCacheMutex.Lock()
	userCache[user.Id] = user
	userCacheMutex.Unlock()

	return user.Id, nil
}

func DeleteUser(id int) error {
	userCacheMutex.RLock()
	_, ok := userCache[id]
	userCacheMutex.RUnlock()

	if !ok {
		return errors.New("user not found")
	}

	userCacheMutex.Lock()
	delete(userCache, id)
	userCacheMutex.Unlock()

	return nil
}
