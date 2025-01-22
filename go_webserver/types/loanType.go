package types

type Loan struct {
	Id             int     `json:"id"`
	UserId         int     `json:"userId"`
	Amount         float32 `json:"amount"`
	InterestRate   float32 `json:"interestRate"`
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
