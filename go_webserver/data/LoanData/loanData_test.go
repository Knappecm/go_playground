package LoanData_test

import (
	"go_playground/go_webserver/data/LoanData"
	"go_playground/go_webserver/types"
	"sync"
	"testing"
)

func setup() *LoanData.LoanDataImpl {
	// Initialize a fresh LoanCache for each test
	loanCache := &types.LoanCache{
		SafeMap: &sync.Map{},
		Count:   0,
	}
	return &LoanData.LoanDataImpl{
		LoanCache: *loanCache,
	}
}

func TestGetLoan(t *testing.T) {
	loanData := setup()

	// Add a loan to the cache
	loan := types.Loan{
		Id:     1,
		UserId: 1001,
		Amount: 5000,
	}
	loanData.LoanCache.SafeMap.Store(loan.Id, loan)

	t.Run("Get existing loan", func(t *testing.T) {
		got, err := loanData.GetLoan(1, 1001)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if got.Id != loan.Id || got.UserId != loan.UserId || got.Amount != loan.Amount {
			t.Errorf("expected %v, got %v", loan, got)
		}
	})

	t.Run("Loan not found", func(t *testing.T) {
		_, err := loanData.GetLoan(2, 1001)
		if err == nil || err.Error() != "loan not found" {
			t.Errorf("expected 'loan not found' error, got %v", err)
		}
	})

	t.Run("Unauthorized loan access", func(t *testing.T) {
		_, err := loanData.GetLoan(1, 2002)
		if err == nil || err.Error() != "you don't have access to this loan" {
			t.Errorf("expected 'you don't have access to this loan' error, got %v", err)
		}
	})
}

func TestCreateLoan(t *testing.T) {
	loanData := setup()

	t.Run("Create loan", func(t *testing.T) {
		loan := types.Loan{
			UserId: 1001,
			Amount: 10000,
		}

		id, err := loanData.CreateLoan(loan)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if id != 1 {
			t.Errorf("expected loan ID to be 1, got %d", id)
		}

		// Verify loan is stored
		storedLoan, _ := loanData.LoanCache.SafeMap.Load(id)
		if storedLoan.(types.Loan).Amount != loan.Amount {
			t.Errorf("expected stored loan amount to be %v, got %v", loan.Amount, storedLoan.(types.Loan).Amount)
		}
	})
}

func TestDeleteLoan(t *testing.T) {
	loanData := setup()

	// Add a loan to the cache
	loan := types.Loan{
		Id:     1,
		UserId: 1001,
		Amount: 5000,
	}
	loanData.LoanCache.SafeMap.Store(loan.Id, loan)

	t.Run("Delete existing loan", func(t *testing.T) {
		err := loanData.DeleteLoan(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify loan is removed
		_, ok := loanData.LoanCache.SafeMap.Load(1)
		if ok {
			t.Errorf("expected loan to be deleted, but it still exists")
		}
	})

	t.Run("Delete non-existent loan", func(t *testing.T) {
		err := loanData.DeleteLoan(2)
		if err == nil || err.Error() != "loan not found" {
			t.Errorf("expected 'loan not found' error, got %v", err)
		}
	})
}
