package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/AsetaShadrach/expense-tracker/helpers"
	"github.com/AsetaShadrach/expense-tracker/schemas"
	"github.com/AsetaShadrach/expense-tracker/utils"
	logging "github.com/AsetaShadrach/expense-tracker/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate *validator.Validate
var tracer = *utils.Tracer

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var (
		responseBytes    []byte
		err              error
		validationSchema schemas.UserInputDto
	)

	decodingError := json.NewDecoder(r.Body).Decode(&validationSchema)
	if decodingError != nil {
		fmt.Println("An error occured decoding")
		w.Write([]byte(decodingError.Error()))
		return
	}

	responseBytes, err = schemas.PerformValidation(validationSchema, "USR001")
	if err != nil {
		logging.GeneralLogger.Error("Error occured during input validation --> ", slog.Any("Errors", responseBytes))
		w.WriteHeader(http.StatusExpectationFailed)
	} else {

		userCreationResponse, userCreationError := helpers.CreateUser(r.Context(), validationSchema)

		if userCreationError != nil {
			fmt.Println(userCreationError.Error())

			valErr := schemas.ErrorList{
				ResponseCode: "USR001",
				Message:      "User creation error",
				Errors:       []string{userCreationError.Error()},
			}

			responseBytes, _ = json.MarshalIndent(valErr, "", "	")

			w.WriteHeader(http.StatusInternalServerError)
		} else {
			responseBytes, _ = json.MarshalIndent(userCreationResponse, "", "	")
			w.WriteHeader(http.StatusOK)
			logging.GeneralLogger.Info("Finalizing user creation")
		}
	}

	w.Write(responseBytes)
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

		w.WriteHeader(http.StatusOK)
		w.Write(reponseBytes)
	}
}
