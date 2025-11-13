package routers

import (
	"github.com/AsetaShadrach/expense-tracker/handlers"
	"github.com/gorilla/mux"
)

func RegisterCategoryRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/user", handlers.CreateCategoryHandler).Methods("POST")
	r.HandleFunc("/api/v1/user/{id}", handlers.GetOrDeleteCategoryHandler).Methods("GET", "DELETE")
	r.HandleFunc("/api/v1/user/update/{id}", handlers.UpdateCategoryHandler).Methods("PUT", "PATCH")
}
