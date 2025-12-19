package helpers

import (
	"context"
	"strings"

	"github.com/AsetaShadrach/expense-tracker/schemas"
	"github.com/AsetaShadrach/expense-tracker/utils"
	"gorm.io/gorm"
)

func CreateGroup(ctx context.Context, data schemas.GroupInputDto) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helpers.createGroup")
	defer span.End()

	group := schemas.Group{
		Name:       data.Name,
		GroupPhoto: data.GroupPhoto,
		CreatedBy:  data.CreatedBy,
		Admins:     strings.Join(data.Admins, ","),
	}

	result := gorm.WithResult()
	err = gorm.G[schemas.Group](schemas.DB, result).Create(ctx, &group)
	if err != nil {
		groupCreationError := schemas.ErrorList{
			ResponseCode: "GR001",
			Message:      "An error occured",
			Errors:       []string{"err"},
		}

		utils.GeneralLogger.Error("an error occured %v", result)

		response, _ = schemas.ConvertStructToMap(groupCreationError)
	}

	response, _ = schemas.ConvertStructToMap(group)

	return response, err
}

func UpdateGroup(ctx context.Context, groupId int, updateSchema schemas.GroupUpdateDto) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helpers.updateGroup")
	defer span.End()

	group := schemas.Group{
		Name:       updateSchema.Name,
		GroupPhoto: updateSchema.GroupPhoto,
		UpdatedBy:  updateSchema.UpdatedBy,
		Admins:     strings.Join(updateSchema.Admins, ","),
	}

	_, err = gorm.G[schemas.Group](schemas.DB).Where("id = ? ", groupId).Updates(ctx, group)

	if err != nil {
		groupUpdateError := schemas.ErrorList{
			ResponseCode: "GR001",
			Message:      "An error occured",
			Errors:       []string{"err"},
		}

		utils.GeneralLogger.Error("an error occured %v", err.Error())

		response, _ = schemas.ConvertStructToMap(groupUpdateError)
	}

	response, _ = schemas.ConvertStructToMap(group)

	return nil, err
}

func FilterGroups(ctx context.Context, queryParams *map[string]string) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helpers.filterGroups")
	defer span.End()

	return nil, err
}

func GetOrDeleteGroup(ctx context.Context, groupId int, method string) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helpers.getOrDeleteGroup")
	defer span.End()

	var group schemas.Group

	if method == "GET" {
		group, err = gorm.G[schemas.Group](schemas.DB).Where("id = ? ", groupId).First(ctx)
		response, _ = schemas.ConvertStructToMap(group)
	} else {
		_, err = gorm.G[schemas.Group](schemas.DB).Where("id = ? ", groupId).Delete(ctx)
		response = map[string]interface{}{
			"message": "successful",
		}
	}

	return response, err
}
