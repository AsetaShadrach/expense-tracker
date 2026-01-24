package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AsetaShadrach/expense-tracker/helpers"
	"github.com/AsetaShadrach/expense-tracker/schemas"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func CreateCashFlowHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "createCashFlowHandler")
	defer span.End()

	var (
		validatedData schemas.CashFlowCreateDto
		responseBytes []byte
	)

	er := json.NewDecoder(r.Body).Decode(&validatedData)
	if er != nil {
		errorList := schemas.ErrorList{
			ResponseCode: "CF001",
			Message:      "An error occured",
			Errors:       []string{er.Error()},
		}

		w.WriteHeader(http.StatusBadRequest)
		responseBytes, _ = json.Marshal(errorList)
		w.Write(responseBytes)

		return
	}

	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(validatedData)

	if err != nil {
		validationErrors := validator.ValidationErrors(err.(validator.ValidationErrors))
		validationErrorsList := schemas.TranslateValidationErrors(validationErrors, validate)

		errorList := schemas.ErrorList{
			ResponseCode: "CF001",
			Message:      "An error occured",
			Errors:       validationErrorsList,
		}

		w.WriteHeader(http.StatusExpectationFailed)
		responseBytes, _ = json.Marshal(errorList)
		w.Write(responseBytes)
		return
	}

	resp, err := helpers.CreateCashFlow(r.Context(), validatedData)
	if err != nil {
		errorList := schemas.ErrorList{
			ResponseCode: "CF001",
			Message:      "An error occured adding cashflow",
			Errors:       []string{err.Error()},
		}

		w.WriteHeader(http.StatusBadRequest)
		responseBytes, _ = json.Marshal(errorList)
		w.Write(responseBytes)
		return
	}
	responseBytes, _ = json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func FilterCashFlowHandler(w http.ResponseWriter, r *http.Request) {

	var responseBytes []byte

	vars := mux.Vars(r)
	topicId, _ := strconv.Atoi(vars["id"])

	response, err := helpers.FetchCashFlowTree(r.Context(), topicId)

	if err != nil {
		response := schemas.ErrorList{
			ResponseCode: "GR002",
			Message:      "An error occured",
			Errors:       []string{err.Error()},
		}

		responseBytes, _ = json.Marshal(response)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		responseBytes, _ = json.Marshal(response)
		w.WriteHeader(http.StatusOK)
	}

	w.Write(responseBytes)
	return

}

func GUDCashFlowHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	w.Write([]byte("TODO"))
}

func SummaryCashFlowHandler(w http.ResponseWriter, r *http.Request) {

	var responseBytes []byte

	vars := mux.Vars(r)
	topicId, _ := strconv.Atoi(vars["id"])

	response, err := helpers.FetchCashFlowTree(r.Context(), topicId)

	if err != nil {
		response := schemas.ErrorList{
			ResponseCode: "GR002",
			Message:      "An error occured",
			Errors:       []string{err.Error()},
		}

		responseBytes, _ = json.Marshal(response)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		responseBytes, err = json.Marshal(response)
		if err != nil {
			responseBytes = []byte(fmt.Sprintf("%s", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}

	w.Write(responseBytes)
	return

}
