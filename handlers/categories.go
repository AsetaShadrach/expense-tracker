package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AsetaShadrach/expense-tracker/helpers"
	"github.com/AsetaShadrach/expense-tracker/schemas"
	"github.com/go-playground/validator/v10"
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

func filterCategorysHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	w.Write([]byte("TODO"))
}

func GetOrDeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	w.Write([]byte("TODO"))
}
