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
	"github.com/gorilla/mux"
)

func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "createGroupHandler")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")

	var (
		validationData schemas.GroupInputDto
		responseBytes  []byte
		err            error
	)

	_ = json.NewDecoder(r.Body).Decode(&validationData)

	responseBytes, err = schemas.PerformValidation(validationData, "GRO17")
	if err != nil {
		utils.GeneralLogger.Error("Error occured during group create validation --> ", slog.Any("Errors", string(responseBytes)))
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write(responseBytes)

		return
	}
	response, err := helpers.CreateGroup(r.Context(), validationData)
	if err != nil {
		errs := schemas.ErrorList{
			ResponseCode: "GROO1",
			Message:      "An error occured creating the group",
			Errors:       []string{err.Error()},
		}

		utils.GeneralLogger.Error("An error occured attempting to creating group  ", err)

		responseBytes, _ = json.Marshal(errs)
		w.WriteHeader(http.StatusExpectationFailed)
		w.Write(responseBytes)

		return
	}

	responseBytes, err = json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)

	return
}

func UpdateGroupHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "updateGroupHandler")
	defer span.End()

	var (
		responseBytes    []byte
		err              error
		validationSchema schemas.GroupUpdateDto
	)

	responseBytes, err = schemas.PerformValidation(validationSchema, "GR0017")
	if err != nil {
		utils.GeneralLogger.Error("An error occured validating group update data ", string(responseBytes))

		w.WriteHeader(http.StatusExpectationFailed)
		w.Write(responseBytes)
		return
	}

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	response, err := helpers.UpdateGroup(r.Context(), id, validationSchema)
	responseBytes, _ = json.Marshal(response)

	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func filterGroupsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	w.Write([]byte("TODO"))
}

func GetOrDeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "getOrDeleteGroupHandler")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")

	var responseBytes []byte

	vars := mux.Vars(r)
	groupId, _ := strconv.Atoi(vars["id"])
	response, err := helpers.GetOrDeleteGroup(r.Context(), groupId, r.Method)

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
