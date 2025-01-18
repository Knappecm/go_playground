package types

type User struct {
	Id        int
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Loans     []int  `json:"loans"`
}
