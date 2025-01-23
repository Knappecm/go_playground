package main

import (
	"fmt"
	"go_playground/go_webserver/api/LoanApi"
	"go_playground/go_webserver/api/UserApi"
	"go_playground/go_webserver/bisLogic/LoanLogic"
	"go_playground/go_webserver/bisLogic/UserLogic"
	"go_playground/go_webserver/data/LoanData"
	"go_playground/go_webserver/data/UserData"
	"net/http"
)

func main() {
	userDataService := &UserData.UserDataImpl{}
	userLogicService := &UserLogic.UserLogicImpl{
		UserDataService: userDataService,
	}
	loanDataService := &LoanData.LoanDataImpl{}
	loanLogicService := &LoanLogic.LoanLogicImpl{}

	userHandler := &UserApi.UserHandler{
		UserDataService: userDataService, // Inject the business logic
	}

	loanHandler := &LoanApi.LoanHandler{
		LoanDataService:  loanDataService,
		LoanLogicService: loanLogicService,
		UserDataService:  userDataService,
		UserLogicService: userLogicService,
	}

	mux := http.NewServeMux()
	userHandler.InitializeUserApi(mux)
	loanHandler.InitializeLoanApi(mux)

	fmt.Println("Server Listening to 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err.Error())
	}
}
