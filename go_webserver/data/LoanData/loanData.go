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

type LoanDataImpl struct{}

var loanCache types.LoanCache

func (l *LoanDataImpl) GetLoan(id int, userId int) (types.Loan, error) {
	value, ok := loanCache.SafeMap.Load(id)
	if !ok {
		return types.Loan{}, errors.New("loan not found")
	}
	loan := value.(types.Loan)

	if loan.UserId != userId {
		return types.Loan{}, errors.New("you don't have access to this loan")
	}

	return loan, nil
}

func (l *LoanDataImpl) CreateLoan(loan types.Loan) (int, error) {

	loanCache.Count++
	loan.Id = loanCache.Count
	loanCache.SafeMap.Store(loan.Id, loan)

	return loan.Id, nil
}

func (l *LoanDataImpl) DeleteLoan(id int) error {
	_, ok := loanCache.SafeMap.Load(id)
	if !ok {
		return errors.New("loan not found")
	}

	loanCache.SafeMap.Delete(id)

	return nil
}
