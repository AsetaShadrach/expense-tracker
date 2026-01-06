package schemas

import (
	"reflect"
)

type CategoryInputDto struct {
	Name        string `validate:"required,alphanumunicode"`
	Description string `validate:"required,alphanumunicode"`
	Subcategory int
}

func (dto CategoryInputDto) GetValidatorName() string {
	return reflect.TypeOf(dto).Name()
}

type UserInputDto struct {
	Username     string `json:"username" validate:"required,alphanumunicode,min=2"`
	ProfilePhoto string `json:"profile_photo" validate:"omitempty,base64"`
	Email        string `json:"email" validate:"required,email"`
	Groups       []int  `json:"groups" validate:"omitempty,gte=1,dive,number"`
}

func (dto UserInputDto) GetValidatorName() string {
	return reflect.TypeOf(dto).Name()
}

type GroupInputDto struct {
	Name       string   `json:"name" validate:"required,alphanumunicode"`
	GroupPhoto string   `json:"group_photo" validate:"omitempty,alphanumunicode"`
	CreatedBy  string   `json:"created_by" validate:"alphanumunicode"`
	Admins     []string `json:"admins" validate:"omitempty,gte=1,dive,alpha"`
	Members    []string `json:"members" validate:"omitempty,gte=1,dive,required"`
}

func (dto GroupInputDto) GetValidatorName() string {
	return reflect.TypeOf(dto).Name()
}

type GroupUpdateDto struct {
	Name       string   `json:"name" validate:"omitempty,alphanumunicode"`
	GroupPhoto string   `json:"group_photo" validate:"omitempty,alphanumunicode"`
	UpdatedBy  string   `json:"updated_by" validate:"required,alphanumunicode"`
	Admins     []string `json:"admins" validate:"omitempty,gte=1,dive,required"`
	Members    []string `json:"members" validate:"omitempty,gte=1,dive,required"`
}

func (dto GroupUpdateDto) GetValidatorName() string {
	return reflect.TypeOf(dto).Name()
}
