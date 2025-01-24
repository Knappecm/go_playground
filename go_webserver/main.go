package main

import (
	"fmt"
	"go_playground/go_webserver/api/LoanApi"
	"go_playground/go_webserver/api/UserApi"
	"go_playground/go_webserver/bisLogic/LoanLogic"
	"go_playground/go_webserver/bisLogic/UserLogic"
	"go_playground/go_webserver/data/LoanData"
	"go_playground/go_webserver/data/UserData"
	"go_playground/go_webserver/types"
	"net/http"
	"sync"
)

func main() {
	// Initialize the UserCache
	userCache := types.UserCache{
		SafeMap: &sync.Map{},
		Count:   0,
	}

	// Initialize the LoanCache
	loanCache := types.LoanCache{
		SafeMap: &sync.Map{},
		Count:   0,
	}

	userDataService := &UserData.UserDataImpl{UserCache: userCache}
	userLogicService := &UserLogic.UserLogicImpl{
		UserDataService: userDataService,
	}
	
	loanDataService := &LoanData.LoanDataImpl{LoanCache: loanCache}
	loanLogicService := &LoanLogic.LoanLogicImpl{}

	userHandler := &UserApi.UserHandler{
		UserDataService: userDataService,
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
