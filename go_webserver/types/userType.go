package types

type User struct {
	Id        int
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type CreateUserReturnType struct {
	Id int `json:"id"`
}