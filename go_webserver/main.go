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
	// went with Caches just to mock DBs 
	// setting up a DB for the project seemed overkill 
	userCache := types.UserCache{
		SafeMap: &sync.Map{},
		Count:   0,
	}

	// Initialize the LoanCache
	loanCache := types.LoanCache{
		SafeMap: &sync.Map{},
		Count:   0,
	}

	// Setting up the different Services 
	// Services can be broken down into 2 main with 3 sub services each
	// User: UserData, UserLogic, and UserApi
	userDataService := &UserData.UserDataImpl{UserCache: userCache}
	userLogicService := &UserLogic.UserLogicImpl{
		UserDataService: userDataService,
	}

	//Loan: LoanData, LoanLogic, and LoanApi
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

	// initialize a serve Mux and add the func handlers
	mux := http.NewServeMux()
	userHandler.InitializeUserApi(mux)
	loanHandler.InitializeLoanApi(mux)


	// begin the server
	fmt.Println("Server Listening to 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err.Error()) // Any unexpected errors get caught and printed for debugging
	}
}
