package routers

import (
	"github.com/AsetaShadrach/expense-tracker/handlers"
	"github.com/gorilla/mux"
)

func RegisterUserRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/user", handlers.CreateUserHandler).Methods("POST")
	r.HandleFunc("/api/v1/user/{id}", handlers.GetOrDeleteUserHandler).Methods("GET", "DELETE")
	r.HandleFunc("/api/v1/user/update/{id}", handlers.UpdateUserHandler).Methods("PUT", "PATCH")
	r.HandleFunc("/api/v1/users", handlers.FilterUsersHandler).Methods("GET")
}
