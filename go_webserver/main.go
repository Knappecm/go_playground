package main

import (
	"fmt"
	loanApi "go_playground/go_webserver/api/LoanApi"
	userApi "go_playground/go_webserver/api/UserApi"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	userApi.InitializeUserApi(mux)
	loanApi.InitializeLoanApi(mux)

	fmt.Println("Server Listening to 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println(err.Error())
	}
}
