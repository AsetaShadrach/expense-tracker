package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AsetaShadrach/expense-tracker/helpers"
	"github.com/AsetaShadrach/expense-tracker/schemas"
	"github.com/go-playground/validator/v10"
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
	fmt.Println(r.URL.Path)
	w.Write([]byte("TODO"))
}

func GUDCashFlowHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	w.Write([]byte("TODO"))
}
