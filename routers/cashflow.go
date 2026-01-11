package routers

import (
	"github.com/AsetaShadrach/expense-tracker/handlers"
	"github.com/gorilla/mux"
)

func RegisterCashFlowRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/cashflow", handlers.CreateCashFlowHandler).Methods("POST")
	r.HandleFunc("/api/v1/cashflows", handlers.FilterCashFlowHandler).Methods("GET")
	r.HandleFunc("/api/v1/cashflow/{id}", handlers.GUDCashFlowHandler).Methods("GET", "DELETE", "PUT")
}
