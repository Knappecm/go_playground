package LoanApi

import (
	"encoding/json"
	"fmt"
	"go_playground/go_webserver/bisLogic/LoanLogic"
	"go_playground/go_webserver/bisLogic/UserLogic"
	"go_playground/go_webserver/data/LoanData"
	"go_playground/go_webserver/data/UserData"
	"go_playground/go_webserver/types"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type LoanHandler struct {
	LoanDataService  LoanData.LoanDataService
	LoanLogicService LoanLogic.LoanLogicService
	UserDataService  UserData.UserDataService
	UserLogicService UserLogic.UserLogicService
}

// Initialize all of the loan APIs in one place
func (l *LoanHandler) InitializeLoanApi(mux *http.ServeMux) {
	mux.HandleFunc("POST /loan", l.CreateLoan)
	mux.HandleFunc("GET /loans/user/{id}", l.GetAllLoansForUser)
	mux.HandleFunc("GET /loan/{id}", l.GetLoan)
	mux.HandleFunc("DELETE /loan/{id}", l.DeleteLoan)
	mux.HandleFunc("GET /loan/{id}/breakdown", l.GetLoanBreakDown)
	slog.Info("Loan Api initialized")
}

func (l *LoanHandler) CreateLoan(
	w http.ResponseWriter,
	r *http.Request,
) {
	logTraceGroup := slog.Group(
		"request",
		"traceID", uuid.New(),
		"method", "POST",
		"apiPath", "/loan",
	)

	var loan types.Loan
	var errorString string
	err := json.NewDecoder(r.Body).Decode(&loan)
	if err != nil {
		slog.Error(err.Error(), logTraceGroup)
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	if loan.UserId <= 0 {
		errorString += "User Id is required\n"
	}

	if !l.UserDataService.DoesUserExist(loan.UserId) {
		errorString += fmt.Sprintf("User %d does not exist\n", loan.UserId)
	}

	if loan.Amount <= 0 {
		errorString += "Loan amount is required and cannot be less than or equal to 0\n"
	}
	if loan.InterestRate <= 0 {
		errorString += "Interest rate is required and cannot be less than or equal to 0\n"
	}
	if loan.LoanTermMonths <= 0 {
		errorString += "Loan term is required and cannot be less than or equal to 0\n"
	}

	if errorString != "" {
		slog.Error(errorString, logTraceGroup)
		http.Error(
			w,
			errorString,
			http.StatusBadRequest,
		)
		return
	}

	Id, err := l.LoanDataService.CreateLoan(loan)
	l.UserLogicService.AddLoanToUser(loan.UserId, Id)

	if err != nil {
		slog.Error(err.Error(), logTraceGroup)
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
		slog.Error(err.Error(), logTraceGroup)
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Write(jsonID)
}

func (l *LoanHandler) GetLoan(
	w http.ResponseWriter,
	r *http.Request,
) {

	logTraceGroup := slog.Group(
		"request",
		"traceID", uuid.New(),
		"method", "GET",
		"apiPath", "/loan/{id}",
	)

	var result map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		slog.Error(err.Error(), logTraceGroup)
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	userId, ok := result["userId"].(float64)
	if !ok {
		slog.Error("User ID field is invalid", logTraceGroup)
		http.Error(
			w,
			"User ID field is invalid",
			http.StatusBadRequest,
		)
		return
	}

	if !l.UserDataService.DoesUserExist(int(userId)) {
		slog.Error("User does not exist", logTraceGroup)
		http.Error(
			w,
			"User does not exist",
			http.StatusBadRequest,
		)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error(err.Error(), logTraceGroup)
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	loan, err := l.LoanDataService.GetLoan(id, int(userId))
	if err != nil {
		slog.Error(err.Error(), logTraceGroup)
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	if loan.UserId != int(userId) {
		slog.Error("You do not have access to this loan", logTraceGroup)

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
		slog.Error(err.Error(), logTraceGroup)

		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Write(loanAtId)

}

func (l *LoanHandler) GetAllLoansForUser(
	w http.ResponseWriter,
	r *http.Request,
) {
	logTraceGroup := slog.Group(
		"request",
		"traceID", uuid.New(),
		"method", "GET",
		"apiPath", "/loan/{id}",
	)

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error(err.Error(), logTraceGroup)

		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	user, err := l.UserDataService.GetUser(id)
	if err != nil {
		slog.Error(err.Error(), logTraceGroup)

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
			slog.Error(err.Error(), logTraceGroup)

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
		loan, err := l.LoanDataService.GetLoan(v, user.Id)
		if err != nil {
			slog.Error(err.Error(), logTraceGroup)
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
		slog.Error(err.Error(), logTraceGroup)
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Write(AllLoansJson)

}

func (l *LoanHandler) DeleteLoan(
	w http.ResponseWriter,
	r *http.Request,
) {
	logTraceGroup := slog.Group(
		"request",
		"traceID", uuid.New(),
		"method", "DELETE",
		"apiPath", "/loan/{id}",
	)
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error(err.Error(), logTraceGroup)
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	err = l.LoanDataService.DeleteLoan(id)
	if err != nil {
		slog.Error(err.Error(), logTraceGroup)
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (l *LoanHandler) GetLoanBreakDown(
	w http.ResponseWriter,
	r *http.Request,
) {
	logTraceGroup := slog.Group(
		"request",
		"traceID", uuid.New(),
		"method", "DELETE",
		"apiPath", "/loan/{id}/breakdown",
	)
	var result map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		slog.Error(err.Error(), logTraceGroup)
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	userId, ok := result["userId"].(float64)
	if !ok {
		slog.Error(err.Error(), logTraceGroup)
		http.Error(
			w,
			"User ID field is invalid",
			http.StatusBadRequest,
		)
		return
	}

	if !l.UserDataService.DoesUserExist(int(userId)) {
		slog.Error(err.Error(), logTraceGroup)
		http.Error(
			w,
			"User does not exist",
			http.StatusBadRequest,
		)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error(err.Error(), logTraceGroup)
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	loan, err := l.LoanDataService.GetLoan(id, int(userId))
	if err != nil {
		slog.Error(err.Error(), logTraceGroup)
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	LoanBreakDown := l.LoanLogicService.AmortizationSchedule(loan)

	w.Header().Set("Content-Type", "Application/json")
	LoanBreakDownJson, err := json.Marshal(LoanBreakDown)
	if err != nil {
		slog.Error(err.Error(), logTraceGroup)
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Write(LoanBreakDownJson)
}
