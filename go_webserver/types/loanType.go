package types

type Loan struct {
	Id             int     `json:"id"`
	UserId         int     `json:"userId"`
	Amount         float32 `json:"amount"`
	InterestRate   float32 `json:"interestRate"`
	LoanTermMonths int     `json:"loanTermMonths"`
}
