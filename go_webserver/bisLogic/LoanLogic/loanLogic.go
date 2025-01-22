package LoanLogic

import (
	"go_playground/go_webserver/types"
	"math"
)

// func getTotalInterest(loan types.Loan) (float64, error) {

// }

func CalcMonthlyPayment(principal float64, MonthlyInterest float64, loanTermInMonths float64) float64 {
	// formula: EMI = [P * r * ( 1 + r )^n] / [(1 + r)^n - 1]

	monthlyPayment := MonthlyInterest * math.Pow((1+MonthlyInterest), loanTermInMonths)
	monthlyPayment = monthlyPayment / (math.Pow((1+MonthlyInterest), loanTermInMonths) - 1)
	monthlyPayment = monthlyPayment * principal

	return monthlyPayment
}

func GenerateAmortizationSchedule(principal float64, annualInterestRate float64, loanTermInMonths int) types.LoanBreakDown {
	monthlyRate := annualInterestRate / 12 / 100
	monthlyPayment := CalcMonthlyPayment(principal, monthlyRate, float64(loanTermInMonths))
	remainingBalance := principal
	totalInterestPaid := 0.0
	totalPrincipalPaid := 0.0
	monthlyBreakdown := make([]types.LoanAmortization, loanTermInMonths)
	for month := 1; month <= loanTermInMonths; month++ {
		interest := math.Round((remainingBalance*monthlyRate)*100) / 100
		principalPayment := math.Round((monthlyPayment-interest)*100) / 100
		remainingBalance = math.Round((remainingBalance-principalPayment)*100) / 100
		totalInterestPaid = math.Round((totalInterestPaid+interest)*100) / 100
		totalPrincipalPaid = math.Round((totalPrincipalPaid+principalPayment)*100) / 100
		if remainingBalance < 1 {
			remainingBalance = 0
		}

		monthBreakDown :=
			types.LoanAmortization{
				Month:              month,
				PrincipalThisMonth: principalPayment,
				InterestThisMonth:  interest,
				TotalInterestPaid:  totalInterestPaid,
				TotalPrincipalPaid: totalPrincipalPaid,
				TotalRemaining:     remainingBalance,
			}

		monthlyBreakdown[month-1] = monthBreakDown

	}
	totalPaid := totalInterestPaid + principal
	monthlyPayment = math.Round(monthlyPayment*100) / 100
	totalPaid = math.Round(totalPaid*100) / 100
	totalInterestPaid = math.Round(totalInterestPaid*100) / 100
	return types.LoanBreakDown{MonthlyPayment: monthlyPayment, TotalPaid: totalPaid, TotalInterest: totalInterestPaid, MonthlyBreakDown: monthlyBreakdown}
}
