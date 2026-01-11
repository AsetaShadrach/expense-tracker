package helpers

import (
	"context"

	"github.com/AsetaShadrach/expense-tracker/schemas"
	"gorm.io/gorm"
)

func CreateCashFlow(ctx context.Context, data schemas.CashFlowCreateDto) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helpers.createCashFlow")
	defer span.End()

	var resp map[string]interface{}

	cashflow := schemas.CashFlow{
		Amount:          data.Amount,
		Month:           data.Month,
		Day:             data.Day,
		IncomeOrExpense: data.IncomeOrExpense,
		CategoryId:      data.CategoryId,
	}

	result := gorm.WithResult()
	err = gorm.G[schemas.CashFlow](schemas.DB, result).Create(ctx, &cashflow)

	if err != nil {
		return nil, err
	}

	resp, err = schemas.ConvertStructToMap(cashflow)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func FilterCashFlow(ctx context.Context, queryParams *map[string]interface{}) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helpers.FilterCashFlow")
	defer span.End()

	queries := *queryParams
	offset := queries["page"].(int) - 1*queries["items"].(int)

	resp, err := gorm.G[schemas.CashFlow](schemas.DB).Offset(offset).Limit(queries["items"].(int)).Where("").Find(ctx)

	if err != nil {
		return nil, err
	}

	response = map[string]interface{}{
		"items": queries["items"].(int),
		"page":  queries["page"].(int),
		"data":  resp,
	}

	return response, err
}

func UpdateCashFlow(ctx context.Context, categoryId int, updateCashflowSchema schemas.CashFlowUpdateDto) (response map[string]interface{}, err error) {
	return nil, err
}

func GetOrDeleteCashFlow(ctx context.Context, categoryId int, method string, updateCashFlowSchema schemas.CashFlowUpdateDto) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helpers.GUDCashFlow")
	defer span.End()

	var category schemas.Category

	if method == "GET" {
		category, err = gorm.G[schemas.Category](schemas.DB).Where("id = ? ", categoryId).First(ctx)
		response, _ = schemas.ConvertStructToMap(category)
	} else if method == "DELETE" {
		_, err = gorm.G[schemas.Category](schemas.DB).Where("id = ? ", categoryId).Delete(ctx)
		response = map[string]interface{}{
			"message": "successful",
		}
	} else {
		response, err = UpdateCashFlow(ctx, categoryId, updateCashFlowSchema)
	}

	return response, err
}
