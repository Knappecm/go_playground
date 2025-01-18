package UserLogic

import (
	Userdata "go_playground/go_webserver/data/UserData"
)

func AddLoanToUser(userID int, loanId int) error {
	user, err := Userdata.GetUser(userID)
	if err != nil {
		return err
	}

	user.Loans = append(user.Loans, loanId)

	err = Userdata.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}
