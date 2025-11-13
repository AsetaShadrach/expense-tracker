package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AsetaShadrach/expense-tracker/routers"
	"github.com/AsetaShadrach/expense-tracker/schemas"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	envLoadErr := godotenv.Load()
	if envLoadErr != nil {
		log.Fatal("Error loading .env file")
	}

	err := schemas.ConnectDb()
	if err != nil {
		fmt.Println("TODO : Initiate gracefull shutdown -- > Server.Close ?")
	} else {
		fmt.Println("DB connection succesful. Beginning to load routers ---")
	}

	r := mux.NewRouter()

	routers.RegisterCashFlowRoutes(r)
	routers.RegisterCategoryRoutes(r)
	routers.RegisterGroupRoutes(r)
	routers.RegisterUserRoutes(r)

	fmt.Println("Routers loaded succesfully ---")

	srv := &http.Server{
		Handler:      r,
		Addr:         os.Getenv("SERVER_ADDRESS"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
