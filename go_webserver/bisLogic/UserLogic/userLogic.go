package UserLogic

import (
	"go_playground/go_webserver/data/UserData"
)

type UserLogicService interface {
	AddLoanToUser(userID int, loanId int) error
}

type UserLogicImpl struct {
	UserDataService UserData.UserDataService
}

func (u *UserLogicImpl) AddLoanToUser(userID int, loanId int) error {
	user, err := u.UserDataService.GetUser(userID)
	if err != nil {
		return err
	}

	user.Loans = append(user.Loans, loanId)

	err = u.UserDataService.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}
