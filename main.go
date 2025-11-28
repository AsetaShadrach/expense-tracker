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
		fmt.Println("DB connection succesful. Beginning auto migration ---")
		// Auto migrate
		err := schemas.DB.AutoMigrate(schemas.User{}, schemas.Category{}, schemas.Group{})

		if err != nil {
			fmt.Printf("An error occured during migration : %v \n\n", err.Error())
			panic(err.Error())
		}

		fmt.Println("Migration succesfful. Beginning to load routers ---")
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
