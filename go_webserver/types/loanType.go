package types

import "sync"

type Loan struct {
	Id             int     `json:"id"`
	UserId         int     `json:"userId"`
	Amount         float64 `json:"amount"`
	InterestRate   float64 `json:"interestRate"`
	LoanTermMonths int     `json:"loanTermMonths"`
	LoanAlias      string  `json:"loanAlias"`
}

type LoanAmortization struct {
	Month              int     `json:"month"`
	PrincipalThisMonth float64 `json:"principalThisMonth"`
	InterestThisMonth  float64 `json:"interestThisMonth"`
	TotalInterestPaid  float64 `json:"totalInterestPaid"`
	TotalPrincipalPaid float64 `json:"totalPrincipalPaid"`
	TotalRemaining     float64 `json:"totalRemaining"`
}

type LoanBreakDown struct {
	MonthlyPayment   float64            `json:"monthlyPayment"`
	TotalPaid        float64            `json:"totalPaid"`
	TotalInterest    float64            `json:"totalInterest"`
	MonthlyBreakDown []LoanAmortization `json:"monthlyBreakDown"`
}

type LoanCache struct {
	Count   int
	SafeMap *sync.Map
}
