package routers

import (
	"github.com/AsetaShadrach/expense-tracker/handlers"
	"github.com/gorilla/mux"
)

func RegisterCashFlowRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/cashflow", handlers.CreateCashFlowEntryHandler).Methods("POST")
	r.HandleFunc("/api/v1/cashflow/{id}", handlers.GetOrDeleteCashFlowEntryHandler).Methods("GET", "DELETE")
	r.HandleFunc("/api/v1/cashflow/update/{id}", handlers.UpdateCashFlowEntryHandler).Methods("PUT", "PATCH")
}
