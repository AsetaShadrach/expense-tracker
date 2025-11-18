package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AsetaShadrach/expense-tracker/schemas"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)

	var validatedData schemas.UserInputDto

	decodingError := json.NewDecoder(r.Body).Decode(&validatedData)
	if decodingError != nil {
		fmt.Println("An error occured decoding")
		w.Write([]byte(decodingError.Error()))
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(validatedData)

	if err != nil {
		fmt.Println("Error occured during input validation")

		validationErrors := err.(validator.ValidationErrors)
		validationErrStrs := schemas.TranslateValidationErrors(validationErrors, validate)

		valErr := schemas.ErrorList{
			ResponseCode: "I001",
			Message:      "Invalid input",
			Errors:       validationErrStrs,
		}

		val, marshalErr := json.MarshalIndent(valErr, "", "	")

		if marshalErr != nil {
			fmt.Println("-------------- >  ", marshalErr)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write(val)
		return
	}

	fmt.Println("All went well ", validatedData)
	w.Write([]byte("TODO"))
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	w.Write([]byte("TODO"))
}

func filterUsersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	w.Write([]byte("TODO"))
}

func GetOrDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	w.Write([]byte("TODO"))
}
