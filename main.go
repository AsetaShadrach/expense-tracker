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

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
		Handler:      enableCORS(r),
		Addr:         os.Getenv("SERVER_ADDRESS"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}
