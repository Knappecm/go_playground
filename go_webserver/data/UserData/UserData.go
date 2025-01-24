package UserData

import (
	"encoding/json"
	"errors"
	"go_playground/go_webserver/types"
	"io"
)

type UserDataService interface {
	GetUser(id int) (types.User, error)
	DoesUserExist(id int) bool
	UpdateUser(user types.User) error
	CreateUser(body io.ReadCloser) (types.User, error)
	DeleteUser(id int) error
}

type UserDataImpl struct{ UserCache types.UserCache }

// Gets the user off the cache
func (d *UserDataImpl) GetUser(id int) (types.User, error) {
	value, ok := d.UserCache.SafeMap.Load(id)
	if !ok {
		return types.User{}, errors.New("user not found")
	}
	user := value.(types.User)

	return user, nil
}

// Gets the user and returns T/F
func (d *UserDataImpl) DoesUserExist(id int) bool {
	_, ok := d.UserCache.SafeMap.Load(id)
	return ok
}

// Updates the User at that Id mem location
func (d *UserDataImpl) UpdateUser(user types.User) error {
	_, err := d.GetUser(user.Id)
	if err != nil {
		return err
	}

	d.UserCache.SafeMap.Store(user.Id, user)

	return nil
}

// Creates a new User
func (d *UserDataImpl) CreateUser(body io.ReadCloser) (types.User, error) {
	var user types.User
	var errorString string
	err := json.NewDecoder(body).Decode(&user)
	if err != nil {
		return types.User{}, err
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
		return types.User{}, errors.New(errorString)
	}

	d.UserCache.Count++
	user.Id = d.UserCache.Count
	d.UserCache.SafeMap.Store(user.Id, user)

	return user, nil
}

// Removes the user at the mem location id
func (d *UserDataImpl) DeleteUser(id int) error {
	_, ok := d.UserCache.SafeMap.Load(id)
	if !ok {
		return errors.New("user not found")
	}

	d.UserCache.SafeMap.Delete(id)

	return nil
}
