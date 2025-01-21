package LoanLogic

import (
	"go_playground/go_webserver/types"
	"math"
)

// func getTotalInterest(loan types.Loan) (float64, error) {

// }

func CalcMonthlyPayment(principal float64, MonthlyInterest float64, loanTermInMonths float64) float64 {
	// formula: EMI = [P * r * ( 1 + r )^n] / [(1 + r)^n - 1]
	return principal * MonthlyInterest * math.Pow(1+MonthlyInterest, loanTermInMonths) / math.Pow(1+MonthlyInterest, loanTermInMonths-1)
}

func GenerateAmortizationSchedule(principal float64, annualInterestRate float64, loanTermInMonths int) []types.LoanAmortization {
	monthlyPayment := CalcMonthlyPayment(principal, annualInterestRate, float64(loanTermInMonths))
	monthlyRate := annualInterestRate / 12 / 100
	remainingBalance := principal
	totalInterestPaid := 0.0

	monthlyBreakdown := make([]types.LoanAmortization, loanTermInMonths)
	for month := 1; month <= loanTermInMonths; month++ {
		interest := remainingBalance * monthlyRate
		principalPayment := monthlyPayment - interest
		remainingBalance -= principalPayment
		totalInterestPaid += interest

		if remainingBalance < 0 {
			remainingBalance = 0
		}

		monthBreakDown := types.LoanAmortization{Month: month, Principal: remainingBalance, InterestPaid: totalInterestPaid}
		monthlyBreakdown[month-1] = monthBreakDown

	}

	return monthlyBreakdown
}
