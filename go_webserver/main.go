package main

import (
	"fmt"
	"go_playground/go_webserver/api"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	api.InitializeUserApi(mux)
	api.InitializeLoanApi(mux)

	fmt.Println("Server Listening to 8080")
	http.ListenAndServe(":8080", mux)
}
