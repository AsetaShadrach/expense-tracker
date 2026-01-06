package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/AsetaShadrach/expense-tracker/helpers"
	"github.com/AsetaShadrach/expense-tracker/schemas"
	"github.com/AsetaShadrach/expense-tracker/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var responseVal []byte

	var validatedData schemas.CategoryInputDto

	er := json.NewDecoder(r.Body).Decode(&validatedData)

	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("An error occured. %v", er.Error())))

		return
	}

	//  using "validate" declared in users handler
	validate = validator.New(validator.WithRequiredStructEnabled())
	catErr := validate.Struct(validatedData)

	if catErr != nil {
		validationErrors := catErr.(validator.ValidationErrors)
		validationErrorsList := schemas.TranslateValidationErrors(validationErrors, validate)

		fmt.Println("Error occured during input validation --> ", validationErrorsList)
		errorList := schemas.ErrorList{
			ResponseCode: "CAT001",
			Message:      "Error occured validating category",
			Errors:       validationErrorsList,
		}

		responseVal, _ = json.Marshal(errorList)
		w.WriteHeader(http.StatusExpectationFailed)
	} else {
		response, categoryCreationErr := helpers.CreateCategory(r.Context(), validatedData)
		if categoryCreationErr != nil {
			fmt.Println(categoryCreationErr.Error())

			w.WriteHeader(http.StatusInternalServerError)
			valErr := schemas.ErrorList{
				ResponseCode: "CAT001",
				Message:      "Category creation error",
				Errors:       []string{categoryCreationErr.Error()},
			}

			responseVal, _ = json.Marshal(valErr)
		} else {
			responseVal, _ = json.Marshal(response)
			w.WriteHeader(http.StatusCreated)
		}
	}

	w.Write(responseVal)

	return
}

func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	w.Write([]byte("TODO"))
}

func FilterCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "filterCategoriesHandler")
	defer span.End()

	queryParams := make(map[string]interface{})

	singleIntParams := []string{"items", "page"}

	for key, val := range r.URL.Query() {
		if slices.Contains(singleIntParams, key) {
			queryParams[key], _ = strconv.Atoi(val[0])
		} else {
			queryParams[key] = val
		}
	}

	response, err := helpers.FilterCategories(r.Context(), &queryParams)

	if err != nil {
		errResponse := schemas.ErrorList{
			Message:      "An error occured",
			ResponseCode: "FTOO9",
			Errors:       []string{err.Error()},
		}

		errorBytes, _ := json.Marshal(errResponse)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorBytes)
	}

	fmt.Println(r.URL.Path)
	responseBytes, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func GUDCategoryHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "gudCategoryHandler")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")

	var (
		responseBytes    []byte
		err              error
		validationSchema schemas.GroupUpdateDto
	)

	if strings.Contains("PUT,PATCH", r.Method) {
		_ = json.NewDecoder(r.Body).Decode(&validationSchema)
		responseBytes, err = schemas.PerformValidation(validationSchema, "GR0017")
		if err != nil {
			utils.GeneralLogger.Error("An error occured validating group update data ", string(responseBytes))
			w.WriteHeader(http.StatusExpectationFailed)
			w.Write(responseBytes)
			return
		}
	}

	vars := mux.Vars(r)
	groupId, _ := strconv.Atoi(vars["id"])

	response, err := helpers.GUDGroup(r.Context(), groupId, r.Method, validationSchema)

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
