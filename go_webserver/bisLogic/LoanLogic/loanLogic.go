package LoanLogic

import (
	"go_playground/go_webserver/types"
	"math"
)

type LoanLogicService interface {
	MonthlyPayment(loan types.Loan, monthlyInterest float64) float64
	AmortizationSchedule(loan types.Loan) types.LoanBreakDown
}

type LoanLogicImpl struct{}

func (u *LoanLogicImpl) MonthlyPayment(loan types.Loan, monthlyInterest float64) float64 {
	// formula: EMI = [P * r * ( 1 + r )^n] / [(1 + r)^n - 1]

	monthlyPayment := monthlyInterest * math.Pow((1+monthlyInterest), float64(loan.LoanTermMonths))
	monthlyPayment = monthlyPayment / (math.Pow((1+monthlyInterest), float64(loan.LoanTermMonths)) - 1)
	monthlyPayment = monthlyPayment * loan.Amount

	return monthlyPayment
}

func (u *LoanLogicImpl) AmortizationSchedule(loan types.Loan) types.LoanBreakDown {
	monthlyRate := loan.InterestRate / 12 / 100
	loanBreakdown := types.LoanBreakDown{
		MonthlyPayment: math.Round(u.MonthlyPayment(loan, monthlyRate)*100) / 100,
	}

	remainingBalance := loan.Amount
	totalPrincipalPaid := 0.0
	monthlyBreakdown := make([]types.LoanAmortization, loan.LoanTermMonths)

	for month := 1; month <= loan.LoanTermMonths; month++ {
		interest := math.Round((remainingBalance*monthlyRate)*100) / 100
		principalPayment := math.Round((loanBreakdown.MonthlyPayment-interest)*100) / 100
		remainingBalance = math.Round((remainingBalance-principalPayment)*100) / 100
		loanBreakdown.TotalInterest = math.Round((loanBreakdown.TotalInterest+interest)*100) / 100
		totalPrincipalPaid = math.Round((totalPrincipalPaid+principalPayment)*100) / 100

		if remainingBalance < 1 {
			remainingBalance = 0
		}

		monthBreakDown :=
			types.LoanAmortization{
				Month:              month,
				PrincipalThisMonth: principalPayment,
				InterestThisMonth:  interest,
				TotalInterestPaid:  loanBreakdown.TotalInterest,
				TotalPrincipalPaid: totalPrincipalPaid,
				TotalRemaining:     remainingBalance,
			}

		monthlyBreakdown[month-1] = monthBreakDown

	}
	loanBreakdown.MonthlyBreakDown = monthlyBreakdown
	loanBreakdown.TotalPaid = math.Round((loanBreakdown.TotalInterest+loan.Amount)*100) / 100

	return loanBreakdown
}
