package data

import (
	"encoding/json"
	"errors"
	"go_playground/go_webserver/types"
	"io"
	"sync"
)

var loanCache = make(map[int]types.User)

var loanCacheMutex sync.RWMutex

func GenerateLoanData(id int) (types.User, error) {
	loanCacheMutex.RLock()
	user, ok := loanCache[id]
	loanCacheMutex.RUnlock()

	if !ok {
		return types.User{}, errors.New("user not found")
	}

	return user, nil
}

func CreateLoan(body io.ReadCloser) (int, error) {
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

	user.Id = len(loanCache) + 1

	loanCacheMutex.Lock()
	loanCache[user.Id] = user
	loanCacheMutex.Unlock()

	return user.Id, nil
}

func DeleteLoan(id int) error {
	loanCacheMutex.RLock()
	_, ok := loanCache[id]
	loanCacheMutex.RUnlock()

	if !ok {
		return errors.New("user not found")
	}

	loanCacheMutex.Lock()
	delete(loanCache, id)
	loanCacheMutex.Unlock()

	return nil
}
