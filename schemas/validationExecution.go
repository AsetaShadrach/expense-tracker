package schemas

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ValidationDto interface {
	GetValidatorName() string
}

// Will always return a response.
// If error is nil, then the response was positive otherwise false and an response will also be returned
// containing a formarted error
func PerformValidation(validationDto ValidationDto, errorResponseCode string) (response []byte, err error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err = validate.Struct(validationDto)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		validationErrorsList := TranslateValidationErrors(validationErrors, validate)

		valErr := ErrorList{
			ResponseCode: errorResponseCode,
			Message:      "Invalid input",
			Errors:       validationErrorsList,
		}

		response, _ = json.MarshalIndent(valErr, "", "	")

	}

	return response, err
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

	validatorPtr.RegisterTranslation("gte", trans, func(ut ut.Translator) error {
		return ut.Add("gt", "Invalid number of value length/size/count on <{0}> . Minimum allowed is {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("gt", fe.Field(), fe.Param())

		return t
	})

	validatorPtr.RegisterTranslation("lte", trans, func(ut ut.Translator) error {
		return ut.Add("lt", "Invalid number of values/characters on <{0}> . Maximum allowed is {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("lt", fe.Field(), fe.Param())

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
