package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/AsetaShadrach/expense-tracker/routers"
	"github.com/AsetaShadrach/expense-tracker/schemas"
	utils "github.com/AsetaShadrach/expense-tracker/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func main() {
	envLoadErr := godotenv.Load()
	if envLoadErr != nil {
		log.Fatal("Error loading .env file")
	}

	// Initiate the logger
	utils.InitiateLogger()

	tp, traceInitErr := utils.InitTracer()

	if traceInitErr != nil {
		log.Fatal(traceInitErr)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			utils.GeneralLogger.Error("Error shutting down trace provider", slog.Any("Errors", err))
		}
	}()

	err := schemas.ConnectDb()
	if err != nil {
		utils.GeneralLogger.Info("TODO : Initiate gracefull shutdown -- > Server.Close ?")
	} else {
		utils.GeneralLogger.Info("DB connection succesful. Beginning auto migration ---")
		// Auto migrate
		err := schemas.DB.AutoMigrate(
			schemas.User{},
			schemas.Category{},
			schemas.Group{},
			schemas.CashFlow{},
		)

		if err != nil {
			utils.GeneralLogger.Error(fmt.Sprintf("An error occured during migration : %v \n\n", err.Error()))
			panic(err.Error())
		}

		utils.GeneralLogger.Info("Migration succesful. Beginning to load routers ---")
	}

	r := mux.NewRouter()
	r.Use(otelmux.Middleware("expense-tracker-server"))

	// InsertMiddlewares
	r.Use(
		utils.ErrorResolver,
		utils.InstrumentRequest,
	)

	routers.RegisterCashFlowRoutes(r)
	routers.RegisterCategoryRoutes(r)
	routers.RegisterGroupRoutes(r)
	routers.RegisterUserRoutes(r)

	utils.GeneralLogger.Info("Routers loaded succesfully ---")

	srv := &http.Server{
		Handler:      r,
		Addr:         os.Getenv("SERVER_ADDRESS"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
