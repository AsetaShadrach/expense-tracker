package schemas

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type CategoryInputDto struct {
	Name        string `validate:"required,alphanumunicode"`
	Description string `validate:"required,alphanumunicode"`
	Subcategory int
}

type UserInputDto struct {
	Username     string `json:"username" validate:"required,alphanum,min=2"`
	ProfilePhoto string `json:"profile_photo" validate:"omitempty,base64"`
	Email        string `json:"email" validate:"required,email"`
	Groups       []int  `json:"groups" validate:"omitempty,gt=1,dive,number"`
}

type GroupInputDto struct {
	Name       string   `validate:"required,alphanumunicode"`
	GroupPhoto string   `validate:"alpha"`
	CreatedBy  string   `validate:"alpha"`
	Admins     []string `validate:"gt=0,dive,alpha"`
}

func ValidationErrorTranslation(trans ut.Translator, validatorPtr *validator.Validate) {
	// Use the tag names and not the struct names
	validatorPtr.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		fmt.Println("Tag Name  ", name)
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}

		return name
	})

	validatorPtr.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "<{0}> is empty or null or missing a value", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())

		return t
	})

	validatorPtr.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "Invalid formart on <{0}> ", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())

		return t
	})

	validatorPtr.RegisterTranslation("alpha", trans, func(ut ut.Translator) error {
		return ut.Add("alpha", "Invalid formart on <{0}> .Expected letters ; [A-Z,a-z]", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alpha", fe.Field())

		return t
	})

	validatorPtr.RegisterTranslation("alphanumunicode", trans, func(ut ut.Translator) error {
		return ut.Add("alphanumunicode", "Invalid formart on <{0}> .Only letters and numbers allowed", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alphanumunicode", fe.Field())

		return t
	})
}

func TranslateValidationErrors(e validator.ValidationErrors, validatorPtr *validator.Validate) []string {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	ValidationErrorTranslation(trans, validatorPtr)

	var validationErrorStrs []string

	for _, er := range e {
		// can translate each error one at a time.
		validationErrorStrs = append(validationErrorStrs, er.Translate(trans))
	}
	return validationErrorStrs
}
