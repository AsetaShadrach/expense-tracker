package schemas

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name        string
	Description string
	SubCategory int
}

type User struct {
	gorm.Model
	Username     string
	ProfilePhoto string
	Email        int
	Groups       int
}

type Group struct {
	gorm.Model
	Name       string   `json:"name"`
	GroupPhoto string   `json:"photoUrl"`
	CreatedBy  string   `json:"createdBy"`
	Admins     []string `json:"admins"`
}
