package loanApi

import (
	"encoding/json"
	"go_playground/go_webserver/data/LoanData"
	Userdata "go_playground/go_webserver/data/UserData"
	"go_playground/go_webserver/types"
	"net/http"
	"strconv"
)

func InitializeLoanApi(mux *http.ServeMux) {
	mux.HandleFunc("POST /loan", CreateLoan)
	mux.HandleFunc("GET /loans/user/{id}", GetAllLoansForUser)
	mux.HandleFunc("GET /loan/{id}", GetLoan)
	mux.HandleFunc("DELETE /loan/{id}", DeleteLoan)
}

func CreateLoan(
	w http.ResponseWriter,
	r *http.Request,
) {
	Id, err := LoanData.CreateLoan(r.Body)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	jsonID, err := json.Marshal(types.Loan{Id: Id})
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
	var result map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	userId, ok := result["userId"].(int)
	if !ok {
		http.Error(
			w,
			"User ID field is invalid",
			http.StatusBadRequest,
		)
		return
	}

	if !Userdata.DoesUserExist(userId) {
		http.Error(
			w,
			"User does not exist",
			http.StatusBadRequest,
		)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	loan, err := LoanData.GetLoan(id, userId)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	if loan.UserId != userId {
		http.Error(
			w,
			"you do not have access to this loan",
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

func GetAllLoansForUser(
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

	if len(user.Loans) == 0 {
		emptyLoan, err := json.Marshal([]types.Loan{})
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}

		w.Write(emptyLoan)
	}

	allLoans := []types.Loan{}

	for _, v := range user.Loans {
		loan, err := LoanData.GetLoan(v, user.Id)
		if err != nil {
			http.Error(
				w,
				err.Error(),
				http.StatusInternalServerError,
			)
			return
		}
		allLoans = append(allLoans, loan)
	}

	w.Header().Set("Content-Type", "Application/json")
	AllLoansJson, err := json.Marshal(allLoans)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Write(AllLoansJson)

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

	err = LoanData.DeleteLoan(id)
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
