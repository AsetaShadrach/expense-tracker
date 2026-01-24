package helpers

import (
	"context"
	"fmt"

	"github.com/AsetaShadrach/expense-tracker/schemas"
	"github.com/AsetaShadrach/expense-tracker/utils"
	"gorm.io/gorm"
)

func CreateCategory(ctx context.Context, categoryValidator schemas.CategoryInputDto) (response map[string]interface{}, err error) {
	category := schemas.Category{
		Name:           categoryValidator.Name,
		Description:    categoryValidator.Description,
		CategoryType:   categoryValidator.CategoryType,
		ParentCategory: categoryValidator.ParentCategory,
		CreatedBy:      categoryValidator.CreatedBy,
	}

	result := gorm.WithResult()
	groupCreateErr := gorm.G[schemas.Category](schemas.DB, result).Create(ctx, &category)

	if groupCreateErr != nil {
		return nil, groupCreateErr
	} else {
		categoryMap, _ := schemas.ConvertStructToMap(category)
		return categoryMap, nil
	}
}

func UpdateCategory(
	ctx context.Context,
	categoryId int,
	updateCategorySchema schemas.CategoryUpdateDto) (response map[string]interface{}, err error) {

	_, span := tracer.Start(ctx, "helpers.UpdateCategory")
	defer span.End()

	category := schemas.Category{
		Name:           updateCategorySchema.Name,
		Description:    updateCategorySchema.Description,
		CategoryType:   updateCategorySchema.CategoryType,
		ParentCategory: updateCategorySchema.ParentCategory,
		UpdatedBy:      updateCategorySchema.UpdatedBy,
	}

	rowsAffected, err := gorm.G[schemas.Category](schemas.DB).Where("id = ? ", categoryId).Updates(ctx, category)
	if err != nil {
		span.SetAttributes(utils.MapToAttributes(
			map[string]interface{}{"errors": err.Error()})...,
		)
		return nil, err
	}
	if rowsAffected < 1 {
		errData := schemas.ErrorList{
			Message:      "No changes made",
			ResponseCode: "CATOO1",
			Errors:       []string{fmt.Sprintf("Category with ID %d not found", categoryId)},
		}

		return schemas.ConvertStructToMap(errData)
	}

	// If update happend it means the group exists
	updatedGroup, _ := gorm.G[schemas.Group](schemas.DB).Where("id = ? ", categoryId).First(ctx)
	return schemas.ConvertStructToMap(updatedGroup)

}

func FilterCategories(ctx context.Context, queryParams *map[string]interface{}) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helpers.filterCategories")
	defer span.End()

	queries := *queryParams
	offset := queries["page"].(int) - 1*queries["items"].(int)

	resp, err := gorm.G[schemas.Category](schemas.DB).Offset(offset).Limit(queries["items"].(int)).Where("").Find(ctx)

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

func GUDCategory(
	ctx context.Context,
	categoryId int,
	method string,
	updateCategorySchema schemas.CategoryUpdateDto) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helpers.getOrDeleteCategory")
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
		response, err = UpdateCategory(ctx, categoryId, updateCategorySchema)
	}

	return response, err
}

func FetchCategoryTree(ctx context.Context, associationId int) (response map[string]interface{}, err error) {
	categories, err := gorm.G[schemas.Category](schemas.DB).Where("associationId = ? ", associationId).Group("name").Find(ctx)

	response = map[string]interface{}{
		"data": categories,
	}

	return response, err
}
