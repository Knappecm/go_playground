package api

import (
	"encoding/json"
	"go_playground/go_webserver/data"
	"go_playground/go_webserver/types"
	"net/http"
	"strconv"
)

func InitializeUserApi(mux *http.ServeMux) {
	mux.HandleFunc("POST /users", CreateUser)
	mux.HandleFunc("GET /users/{id}", GetUser)
	mux.HandleFunc("DELETE /users/{id}", deleteUser)
}

func CreateUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	userId, err := data.CreateUser(r.Body)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	jsonID, err := json.Marshal(types.CreateUserReturnType{Id: userId})
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Write(jsonID)
}

func GetUser(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	user, err := data.GetUser(id)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	userAtId, err := json.Marshal(user)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Write(userAtId)

}

func deleteUser(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	err = data.DeleteUser(id)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
