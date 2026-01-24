package helpers

import (
	"context"
	"fmt"
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

		return schemas.ConvertStructToMap(groupCreationError)
	}

	return schemas.ConvertStructToMap(group)

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

	var rowsAffected int

	rowsAffected, err = gorm.G[schemas.Group](schemas.DB).Where("id = ? ", groupId).Updates(ctx, group)
	if err != nil {
		groupUpdateError := schemas.ErrorList{
			ResponseCode: "GR001",
			Message:      "An error occured",
			Errors:       []string{fmt.Sprintf("An error occured updating group with id %d", groupId)},
		}

		span.SetAttributes(utils.MapToAttributes(
			map[string]interface{}{"errors": []string{err.Error()}})...,
		)

		return schemas.ConvertStructToMap(groupUpdateError)
	}

	if rowsAffected < 1 {
		groupUpdateError := schemas.ErrorList{
			ResponseCode: "GR000",
			Message:      "An error occured",
			Errors:       []string{fmt.Sprintf("Group with ID %d not found", groupId)},
		}

		span.SetAttributes(utils.MapToAttributes(
			map[string]interface{}{"errors": groupUpdateError.Errors})...,
		)

		return schemas.ConvertStructToMap(groupUpdateError)

	}
	// If update happend it means the group exists
	updatedGroup, _ := gorm.G[schemas.Group](schemas.DB).Where("id = ? ", groupId).First(ctx)
	return schemas.ConvertStructToMap(updatedGroup)

}

func FilterGroups(ctx context.Context, queryParams *map[string]interface{}) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helpers.filterGroups")
	defer span.End()

	queryData := *queryParams
	offset := (queryData["page"].(int) - 1) * queryData["items"].(int)
	resp, err := gorm.G[schemas.Group](schemas.DB).Offset(offset).Limit(queryData["items"].(int)).Where("").Find(ctx)
	if err != nil {
		return nil, err
	}

	response = map[string]interface{}{
		"page":  queryData["page"].(int),
		"items": len(resp),
		"data":  resp,
	}

	return response, nil
}

func GUDGroup(ctx context.Context, groupId int, method string, updateGroupSchema schemas.GroupUpdateDto) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helpers.gudGroup")
	defer span.End()

	var group schemas.Group

	if method == "GET" {
		group, err = gorm.G[schemas.Group](schemas.DB).Where("id = ? ", groupId).First(ctx)
		return schemas.ConvertStructToMap(group)
	} else if method == "DELETE" {
		_, err = gorm.G[schemas.Group](schemas.DB).Where("id = ? ", groupId).Delete(ctx)
		response = map[string]interface{}{
			"message": "successful",
		}
		return response, nil
	} else {
		return UpdateGroup(ctx, groupId, updateGroupSchema)
	}
}
