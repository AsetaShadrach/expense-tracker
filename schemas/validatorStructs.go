package schemas

import (
	"reflect"
)

type CategoryInputDto struct {
	Name           string `validate:"required,alphanumspace"`
	Description    string `validate:"omitempty,alphanumspace"`
	CategoryType   string `json:"category_type" validate:"oneof=topic category"`
	ParentCategory int    `json:"parent_category"`
	CreatedBy      int    `json:"created_by"`
}

func (dto CategoryInputDto) GetValidatorName() string {
	return reflect.TypeOf(dto).Name()
}

type CategoryUpdateDto struct {
	Name           string `validate:"omitempty,alphanumspace"`
	Description    string `validate:"omitempty,alphanumspace"`
	CategoryType   string `json:"category_type" validate:"omitempty,oneof=topic category"`
	ParentCategory int    `json:"parent_category" validate:"omitempty"`
	UpdatedBy      int    `json:"updated_by"  validate:"omitempty"`
}

func (dto CategoryUpdateDto) GetValidatorName() string {
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

type CashFlowCreateDto struct {
	Amount          float64 `json:"amount" validate:"required"`
	Month           int     `json:"month"  validate:"lte=12,gte=1"`
	Day             int     `json:"day" validate:"lte=31,gte=1"`
	IncomeOrExpense string  `json:"income_or_expense" validate:"oneof=income expense"` // Income/Expense
	Description     string  `json:"description"`
	CategoryId      int     `json:"category_id" validate:"required"`
	AssociationId   int     `json:"association_id" validate:"required"`
}

func (dto CashFlowCreateDto) GetValidatorName() string {
	return reflect.TypeOf(dto).Name()
}

type CashFlowUpdateDto struct {
	Amount          float64 `json:"amount" validate:"omitempty"`
	Month           int     `json:"month"  validate:"omitempty,lte=12,gte=1"`
	Day             int     `json:"day" validate:"omitempty,lte=31,gte=1"`
	IncomeOrExpense string  `json:"income_or_expense" validate:"omitempty,oneof=income expense"` // Income/Expense
	Description     string  `json:"description" validate:"omitempty"`
	CategoryId      int     `json:"category_id" validate:"omitempty"`
	AssociationId   int     `json:"association_id" validate:"omitempty"`
}

func (dto CashFlowUpdateDto) GetValidatorName() string {
	return reflect.TypeOf(dto).Name()
}
