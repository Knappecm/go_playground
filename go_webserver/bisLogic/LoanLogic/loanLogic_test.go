package LoanLogic

import (
	"go_playground/go_webserver/types"
	"math"
	"testing"
)

func TestMonthlyPayment(t *testing.T) {
	loanLogic := &LoanLogicImpl{}

	t.Run("Valid Monthly Payment Calculation", func(t *testing.T) {
		loan := types.Loan{
			Amount:         100_000, // Loan principal
			LoanTermMonths: 360,     // Loan term in months
		}
		monthlyInterest := 0.005 // Monthly interest rate (0.5%)

		expectedMonthlyPayment := 599.55 // Known expected value for these inputs

		got := math.Round(loanLogic.MonthlyPayment(loan, monthlyInterest)*100) / 100
		if got != expectedMonthlyPayment {
			t.Errorf("expected %v, got %v", expectedMonthlyPayment, got)
		}
	})
}

func TestAmortizationSchedule(t *testing.T) {
	loanLogic := &LoanLogicImpl{}

	t.Run("Valid Amortization Schedule", func(t *testing.T) {
		loan := types.Loan{
			Amount:         100_000, // Loan principal
			InterestRate:   6,       // Annual interest rate (6%)
			LoanTermMonths: 360,     // Loan term in months
		}

		// Expected values
		expectedMonthlyPayment := 599.55   // Pre-computed monthly payment
		expectedTotalPaid := 215838.45     // Total payment over the loan term
		expectedTotalInterest := 115838.45 // Total interest paid

		loanBreakdown := loanLogic.AmortizationSchedule(loan)

		// Check monthly payment
		if loanBreakdown.MonthlyPayment != expectedMonthlyPayment {
			t.Errorf("expected monthly payment %v, got %v", expectedMonthlyPayment, loanBreakdown.MonthlyPayment)
		}

		// Check total paid
		if loanBreakdown.TotalPaid != expectedTotalPaid {
			t.Errorf("expected total paid %v, got %v", expectedTotalPaid, loanBreakdown.TotalPaid)
		}

		// Check total interest
		if loanBreakdown.TotalInterest != expectedTotalInterest {
			t.Errorf("expected total interest %v, got %v", expectedTotalInterest, loanBreakdown.TotalInterest)
		}

		// Check amortization schedule details
		for i, monthly := range loanBreakdown.MonthlyBreakDown {
			if i == 0 {
				continue
			}
			if monthly.Month != i+1 {
				t.Errorf("expected month %v, got %v", i+1, monthly.Month)
			}

			if monthly.InterestThisMonth > loanBreakdown.MonthlyBreakDown[i-1].InterestThisMonth {
				t.Errorf("expected Interest this Month be lower than lasts %v, got %v", monthly.InterestThisMonth, loanBreakdown.MonthlyBreakDown[i-1].InterestThisMonth)
			}

			if monthly.PrincipalThisMonth < loanBreakdown.MonthlyBreakDown[i-1].PrincipalThisMonth {
				t.Errorf("expected principal this Month be larger than lasts %v, got %v", monthly.PrincipalThisMonth, loanBreakdown.MonthlyBreakDown[i-1].PrincipalThisMonth)

			}

		}
	})
}
