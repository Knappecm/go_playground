package UserApi

import (
	"encoding/json"
	"go_playground/go_webserver/data/UserData"
	"log/slog"
	"net/http"
	"strconv"
)

type UserHandler struct {
	UserDataService UserData.UserDataService
}

func (d *UserHandler) InitializeUserApi(mux *http.ServeMux) {
	mux.HandleFunc("POST /user", d.CreateUser)
	mux.HandleFunc("GET /user/{id}", d.GetUser)
	mux.HandleFunc("DELETE /user/{id}", d.deleteUser)
	slog.Info("User Api initialized")
}

func (d *UserHandler) CreateUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	user, err := d.UserDataService.CreateUser(r.Body)
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
		slog.Error(err.Error())
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Write(jsonID)
}

func (d *UserHandler) GetUser(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error(err.Error())

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	user, err := d.UserDataService.GetUser(id)

	if err != nil {
		slog.Error(err.Error())
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
		slog.Error(err.Error())
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Write(userAtId)

}

func (d *UserHandler) deleteUser(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error(err.Error())
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	err = d.UserDataService.DeleteUser(id)
	if err != nil {
		slog.Error(err.Error())
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
