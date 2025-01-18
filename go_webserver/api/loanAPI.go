package api

import (
	"encoding/json"
	"go_playground/go_webserver/data"
	"go_playground/go_webserver/types"
	"net/http"
	"strconv"
)

func InitializeLoanApi(mux *http.ServeMux) {
	mux.HandleFunc("POST /loan", CreateLoan)
	mux.HandleFunc("GET /loan/{id}", GetLoan)
	mux.HandleFunc("DELETE /loan/{id}", DeleteLoan)
}

func CreateLoan(
	w http.ResponseWriter,
	r *http.Request,
) {
	Id, err := data.CreateLoan(r.Body)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	jsonID, err := json.Marshal(types.CreateUserReturnType{Id: Id})
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

func GetLoan(
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

	loan, err := data.GetLoan(id, 1)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	loanAtId, err := json.Marshal(loan)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Write(loanAtId)

}

func DeleteLoan(
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
