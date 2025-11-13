package schemas

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// This function is here incase there a anything you may want to do before or after db connection
func ConnectDb() error {

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_CONNECTION_STRING")), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed to connect database : %v", err)
		return err
	}

	// Assign the connection pointer to the global DB pointer
	DB = db

	return nil
}
