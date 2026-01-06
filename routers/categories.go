package routers

import (
	"github.com/AsetaShadrach/expense-tracker/handlers"
	"github.com/gorilla/mux"
)

func RegisterCategoryRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/category", handlers.CreateCategoryHandler).Methods("POST")
	r.HandleFunc("/api/v1/categories", handlers.FilterCategoriesHandler).Methods("GET")
	r.HandleFunc("/api/v1/category/{id}", handlers.GUDCategoryHandler).Methods("GET", "PUT", "DELETE")
}
