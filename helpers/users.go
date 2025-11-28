package helpers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/AsetaShadrach/expense-tracker/schemas"
	"gorm.io/gorm"
)

func CreateUser(ctx context.Context, data schemas.UserInputDto) (response map[string]interface{}, err error) {
	user := schemas.User{
		Username:     data.Username,
		Email:        data.Email,
		Groups:       strings.Trim(strings.Replace(fmt.Sprint(data.Groups), " ", ",", -1), "[]"),
		ProfilePhoto: data.ProfilePhoto,
	}

	result := gorm.WithResult()
	err = gorm.G[schemas.User](schemas.DB, result).Create(ctx, &user)

	if err != nil {
		fmt.Println("An error writing to db occured ", err)
		return nil, err
	}

	userJson, conversionErr := schemas.ConvertStructToMap(user)
	if conversionErr != nil {
		return nil, err
	}

	fmt.Println("User created succesfully. ID  --> ", user.ID)

	return userJson, nil
}

func UpdateUser(userId int) (response map[string]interface{}, err error) {
	return nil, err
}

func FilterUsers(ctx context.Context, queryParamsPtr *map[string]string) (response map[string]interface{}, err error) {
	queryParams := *queryParamsPtr

	userList := []schemas.User{}
	schemas.DB.Limit(10).Where("username LIKE ?", queryParams["name"]+"%").Find(&userList)

	userListResponse := make(map[string]interface{})

	userListResponse["page"] = 1
	userListResponse["items"] = 10
	userListResponse["data"] = userList

	return userListResponse, nil
}

func GetorDeleteUser(ctx context.Context, userId int, method string) (response map[string]interface{}, err error) {
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
