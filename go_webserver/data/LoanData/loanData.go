package LoanData

import (
	"errors"
	"go_playground/go_webserver/types"
)

type LoanDataService interface {
	GetLoan(id int, userId int) (types.Loan, error)
	CreateLoan(loan types.Loan) (int, error)
	DeleteLoan(id int) error
}

type LoanDataImpl struct{ LoanCache types.LoanCache }

// GetLoan retrieves a loan from the loan cache by its ID.
// It checks if the loan exists and if the user has access to it.
// Returns the loan and an error if the loan is not found or the user does not have access.
func (l *LoanDataImpl) GetLoan(id int, userId int) (types.Loan, error) {
	value, ok := l.LoanCache.SafeMap.Load(id)
	if !ok {
		return types.Loan{}, errors.New("loan not found")
	}
	loan := value.(types.Loan)

	if loan.UserId != userId {
		return types.Loan{}, errors.New("you don't have access to this loan")
	}

	return loan, nil
}

// CreateLoan adds a new loan to the loan cache.
// It increments the loan count, assigns a new ID, and stores the loan in the cache.
// Returns the new loan's ID and any error encountered.
func (l *LoanDataImpl) CreateLoan(loan types.Loan) (int, error) {

	l.LoanCache.Count++
	loan.Id = l.LoanCache.Count
	l.LoanCache.SafeMap.Store(loan.Id, loan)

	return loan.Id, nil
}

// DeleteLoan removes a loan from the loan cache by its ID.
// Returns an error if the loan is not found in the cache.
func (l *LoanDataImpl) DeleteLoan(id int) error {
	_, ok := l.LoanCache.SafeMap.Load(id)
	if !ok {
		return errors.New("loan not found")
	}

	l.LoanCache.SafeMap.Delete(id)

	return nil
}
