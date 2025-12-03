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

var validate *validator.Validate

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	var responseVal []byte

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
		validationErrors := err.(validator.ValidationErrors)
		validationErrorsList := schemas.TranslateValidationErrors(validationErrors, validate)

		fmt.Println("Error occured during input validation --> ", validationErrorsList)

		valErr := schemas.ErrorList{
			ResponseCode: "USR001",
			Message:      "Invalid input",
			Errors:       validationErrorsList,
		}

		var marshalErr error
		responseVal, marshalErr = json.MarshalIndent(valErr, "", "	")

		if marshalErr != nil {
			fmt.Println("-------------- >  ", marshalErr)
		}

		w.WriteHeader(http.StatusExpectationFailed)
	} else {

		userCreationResponse, userCreationError := helpers.CreateUser(r.Context(), validatedData)

		if userCreationError != nil {
			fmt.Println(userCreationError.Error())

			valErr := schemas.ErrorList{
				ResponseCode: "USR001",
				Message:      "User creation error",
				Errors:       []string{userCreationError.Error()},
			}

			responseVal, _ = json.MarshalIndent(valErr, "", "	")

			w.WriteHeader(http.StatusInternalServerError)
		} else {
			responseVal, _ = json.MarshalIndent(userCreationResponse, "", "	")
			w.WriteHeader(http.StatusOK)
		}
	}

	fmt.Println("Finalizing user creation")
	w.Write(responseVal)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	w.Write([]byte("TODO"))
}

func FilterUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()
	params := make(map[string]string)

	for j, k := range queryParams {
		params[j] = k[0]
	}

	resp, err := helpers.FilterUsers(r.Context(), &params)

	if err != nil {
		valErr := schemas.ErrorList{
			ResponseCode: "P002",
			Message:      "Error occured",
			Errors:       []string{err.Error()},
		}

		responseVal, _ := json.Marshal(valErr)

		w.WriteHeader(http.StatusBadRequest)

		w.Write(responseVal)
	} else {
		reponseBytes, _ := json.Marshal(resp)

		w.WriteHeader(http.StatusOK)
		w.Write(reponseBytes)
	}
}

func GetOrDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userId, _ := strconv.Atoi(vars["id"])

	resp, err := helpers.GetorDeleteUser(r.Context(), userId, r.Method)

	if err != nil {
		valErr := schemas.ErrorList{
			ResponseCode: "P002",
			Message:      "Error occured",
			Errors:       []string{err.Error()},
		}

		responseVal, _ := json.Marshal(valErr)

		w.WriteHeader(http.StatusBadRequest)

		w.Write(responseVal)
	} else {

		reponseBytes, _ := json.Marshal(resp)
		fmt.Println(r.URL.Path)

		w.WriteHeader(http.StatusOK)
		w.Write(reponseBytes)
	}
}
