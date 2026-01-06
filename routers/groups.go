package routers

import (
	"github.com/AsetaShadrach/expense-tracker/handlers"
	"github.com/gorilla/mux"
)

func RegisterGroupRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/group", handlers.CreateGroupHandler).Methods("POST")
	r.HandleFunc("/api/v1/groups", handlers.FilterGroupsHandler).Methods("GET")
	r.HandleFunc("/api/v1/group/{id}", handlers.GUDGroupHandler).Methods("GET", "DELETE", "PUT", "PATCH")
}
