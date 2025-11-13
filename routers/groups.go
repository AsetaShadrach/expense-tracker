package routers

import (
	"github.com/AsetaShadrach/expense-tracker/handlers"
	"github.com/gorilla/mux"
)

func RegisterGroupRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/user", handlers.CreateGroupHandler).Methods("POST")
	r.HandleFunc("/api/v1/user/{id}", handlers.GetOrDeleteGroupHandler).Methods("GET", "DELETE")
	r.HandleFunc("/api/v1/user/update/{id}", handlers.UpdateGroupHandler).Methods("PUT", "PATCH")
}
