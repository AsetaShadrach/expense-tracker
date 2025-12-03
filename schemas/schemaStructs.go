package schemas

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CustomGormModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type Category struct {
	CustomGormModel
	Name        string `json:"name" gorm:"unique"`
	Description string `json:"description"`
	Subcategory int    `json:"subcategory"`
}

type User struct {
	CustomGormModel
	Username     string `gorm:"unique" json:"username"`
	ProfilePhoto string `json:"profile_photo"`
	Email        string `json:"email"`
	Groups       string `json:"groups"`
}

type Group struct {
	CustomGormModel
	Name       string `json:"name" gorm:"unique"`
	GroupPhoto string `json:"photo_url"`
	CreatedBy  string `json:"created_by"`
	Admins     string `json:"admins"`
}

type ErrorList struct {
	ResponseCode string   `json:"response_code"`
	Message      string   `json:"message"`
	Errors       []string `json:"errors"`
}

// Receives an interface that is a struct and converts it to a map
func ConvertStructToMap(customStruct interface{}) (userMap map[string]interface{}, err error) {
	structBytes, _ := json.Marshal(customStruct)
	var structJson map[string]interface{}

	unMshalErr := json.Unmarshal(structBytes, &structJson)
	if unMshalErr != nil {
		fmt.Println("An error occured unmarshalling ", unMshalErr)
		return nil, unMshalErr
	}

	return structJson, nil

}
