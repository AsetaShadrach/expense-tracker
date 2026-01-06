package helpers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/AsetaShadrach/expense-tracker/schemas"
	"github.com/AsetaShadrach/expense-tracker/utils"
	"gorm.io/gorm"
)

var tracer = *utils.Tracer

func CreateUser(ctx context.Context, data schemas.UserInputDto) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helper.createUser")
	defer span.End()

	user := schemas.User{
		Username:     data.Username,
		Email:        data.Email,
		Groups:       strings.Trim(strings.Replace(fmt.Sprint(data.Groups), " ", ",", -1), "[]"),
		ProfilePhoto: data.ProfilePhoto,
	}

	result := gorm.WithResult()
	err = gorm.G[schemas.User](schemas.DB, result).Create(ctx, &user)

	if err != nil {
		utils.GeneralLogger.Error("An error writing to db occured ", slog.Any("Errors", err))
		return nil, err
	}

	userJson, conversionErr := schemas.ConvertStructToMap(user)
	if conversionErr != nil {
		return nil, err
	}

	utils.GeneralLogger.Info("User created succesfully. ID  --> ", slog.Int("id", int(user.ID)))

	return userJson, nil
}

func UpdateUser(userId int) (response map[string]interface{}, err error) {
	return nil, err
}

func FilterUsers(ctx context.Context, queryParamsPtr *map[string]string) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helper.filterUsers")
	defer span.End()

	queryParams := *queryParamsPtr

	userList := []schemas.User{}
	schemas.DB.Limit(10).Where("username LIKE ?", queryParams["name"]+"%").Find(&userList)

	userListResponse := map[string]interface{}{
		"page":  1,
		"items": 10,
		"data":  userList,
	}

	return userListResponse, nil
}

func GetorDeleteUser(ctx context.Context, userId int, method string) (response map[string]interface{}, err error) {
	_, span := tracer.Start(ctx, "helper.getOrDeleteUser")
	defer span.End()

	user, err := gorm.G[schemas.User](schemas.DB).Where("id = ?", userId).First(ctx)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New(fmt.Sprintf("User with id %d not found ", userId))
	}

	if strings.ToUpper(method) == "GET" {

		userJson, conversionErr := schemas.ConvertStructToMap(user)
		if conversionErr != nil {
			return nil, conversionErr
		}

		return userJson, nil
	} else {
		_, delErr := gorm.G[schemas.User](schemas.DB).Where("id = ?", userId).Delete(ctx)

		if delErr != nil {
			return nil, delErr
		}
		response := make(map[string]interface{})
		response["message"] = "User succesfully deleted"
		return response, nil

	}
}
