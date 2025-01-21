package userApi

import (
	"encoding/json"
	Userdata "go_playground/go_webserver/data/UserData"
	"net/http"
	"strconv"
)

func InitializeUserApi(mux *http.ServeMux) {
	mux.HandleFunc("POST /user", CreateUser)
	mux.HandleFunc("GET /user/{id}", GetUser)
	mux.HandleFunc("DELETE /user/{id}", deleteUser)
}

func CreateUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	user, err := Userdata.CreateUser(r.Body)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	jsonID, err := json.Marshal(user)
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

	user, err := Userdata.GetUser(id)

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

	err = Userdata.DeleteUser(id)
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
