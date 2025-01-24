package types

import "sync"

type User struct {
	Id        int
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Loans     []int  `json:"loans"`
}

type UserCache struct {
	Count   int
	SafeMap *sync.Map
}
