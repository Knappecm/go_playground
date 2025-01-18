package LoanData

import (
	"encoding/json"
	"errors"
	"go_playground/go_webserver/bisLogic/UserLogic.go"
	Userdata "go_playground/go_webserver/data/UserData"
	"go_playground/go_webserver/types"
	"io"
	"strconv"
	"sync"
)

var loanCache = make(map[int]types.Loan)
var loanCacheMutex sync.RWMutex

func GetLoan(id int, userId int) (types.Loan, error) {
	loanCacheMutex.RLock()
	loan, ok := loanCache[id]
	loanCacheMutex.RUnlock()

	if !ok {
		return types.Loan{}, errors.New("loan not found")
	}

	if loan.UserId != userId {
		return types.Loan{}, errors.New("you don't have access to this loan")
	}

	return loan, nil
}

func CreateLoan(body io.ReadCloser) (int, error) {
	var loan types.Loan
	var errorString string
	err := json.NewDecoder(body).Decode(&loan)
	if err != nil {
		return 0, err
	}

	if loan.UserId <= 0 {
		errorString += "User Id is required\n"
	}

	_, err = Userdata.GetUser(loan.UserId)
	if err != nil {
		errorString += "User " + strconv.Itoa(loan.UserId) + " does not exist\n"
	}

	if loan.Amount <= 0 {
		errorString += "Loan amount is required and cannot be less than or equal to 0\n"
	}
	if loan.InterestRate <= 0 {
		errorString += "Interest rate is required and cannot be less than or equal to 0\n"
	}
	if loan.LoanTermMonths <= 0 {
		errorString += "Loan term is required and cannot be less than or equal to 0\n"
	}

	if errorString != "" {
		return 0, errors.New(errorString)
	}

	loan.Id = len(loanCache) + 1

	loanCacheMutex.Lock()
	loanCache[loan.Id] = loan
	loanCacheMutex.Unlock()

	UserLogic.AddLoanToUser(loan.UserId, loan.Id)

	return loan.Id, nil
}

func DeleteLoan(id int) error {
	loanCacheMutex.RLock()
	_, ok := loanCache[id]
	loanCacheMutex.RUnlock()

	if !ok {
		return errors.New("loan not found")
	}

	loanCacheMutex.Lock()
	delete(loanCache, id)
	loanCacheMutex.Unlock()

	return nil
}
