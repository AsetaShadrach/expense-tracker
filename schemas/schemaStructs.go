package schemas

import (
	"context"
	"encoding/json"
	"errors"
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
	Name           string `json:"name"`
	Description    string `json:"description"`
	CategoryType   string `json:"category_type"` // Topic or Category
	ParentCategory int    `json:"parent_category"`
	CreatedBy      int    `json:"created_by"`
	UpdatedBy      int    `json:"updated_by"`
}

func (cat *Category) BeforeSave(tx *gorm.DB) (err error) {
	_, err = gorm.G[User](DB).Where("id = ?", cat.CreatedBy).First(context.Background())
	if err != nil {
		err = errors.New(fmt.Sprintf("created_by <%d> missing or not found", cat.CreatedBy))
	}
	return
}

type CashFlow struct {
	CustomGormModel
	AssociationId   int     `json:"association_id"` // Topic ID
	Amount          float64 `json:"amount"`
	Month           int     `json:"month"`
	Day             int     `json:"day"`
	IncomeOrExpense string  `json:"income_or_expense"` // Income/Expense
	Description     string  `json:"description"`
	CategoryId      int     `json:"category_id"`
}

func (cf *CashFlow) BeforeSave(tx *gorm.DB) (err error) {
	foundTopic, err := gorm.G[Category](DB).Where("id = ?", cf.AssociationId).First(context.Background())
	if err != nil {
		err = errors.New(fmt.Sprintf("association_id <%d> missing or not found", cf.AssociationId))
	}
	foundCat, err := gorm.G[Category](DB).Where("id = ?", cf.CategoryId).First(context.Background())
	if err != nil {
		err = errors.New(fmt.Sprintf("category_id <%d> missing or not found", cf.CategoryId))
	}

	if foundTopic.CategoryType != "topic" {
		err = errors.New(fmt.Sprintf("association_id <%d> is invalid. association_id must be category type 'topic'", cf.CategoryId))
	}
	if foundCat.CategoryType == "topic" {
		err = errors.New(fmt.Sprintf("category_id <%d> is a topic. Cashflow association failed", cf.CategoryId))
	}
	return
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
	UpdatedBy  string `json:"updated_by"`
	Admins     string `json:"admins"`
	Members    string `json:"members"`
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
