package helpers

import (
	"context"

	"github.com/AsetaShadrach/expense-tracker/schemas"
	"gorm.io/gorm"
)

func CreateCategory(ctx context.Context, categoryValidator schemas.CategoryInputDto) (response map[string]interface{}, err error) {
	category := schemas.Category{
		Name:        categoryValidator.Name,
		Description: categoryValidator.Description,
		Subcategory: categoryValidator.Subcategory,
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

func UpdateCategory(categoryId int) (response map[string]interface{}, err error) {
	return nil, err
}

func filterCategorys() (response map[string]interface{}, err error) {
	return nil, err
}

func GetCategory(categoryId int) (response map[string]interface{}, err error) {
	return nil, err
}

func DeleteCategory(categoryId int) (response map[string]interface{}, err error) {
	return nil, err
}
