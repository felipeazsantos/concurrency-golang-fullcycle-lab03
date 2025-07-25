package validation

import (
	"encoding/json"
	"errors"

	resterr "github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/configuration/rest_err"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	validator_en "github.com/go-playground/validator/v10/translations/en"
)

var (
	Validate = validator.New()
	transl   ut.Translator
)

func init() {
	if value, ok := binding.Validator.Engine().(*validator.Validate); ok {
		en := en.New()
		enTransl := ut.New(en, en)
		transl, _ = enTransl.GetTranslator("en")
		validator_en.RegisterDefaultTranslations(value, transl)
	}
}

func ValidateErr(validationErr error) *resterr.RestErr {
	var jsonErr *json.UnmarshalTypeError
	var jsonValidation validator.ValidationErrors

	if errors.As(validationErr, &jsonErr) {
		return resterr.NewBadRequestError("invalid type error")
	} else if errors.As(validationErr, &jsonValidation) {
		errorCauses := []resterr.Causes{}

		for _, e := range jsonValidation {
			errorCauses = append(errorCauses, resterr.Causes{
				Field:   e.Field(),
				Message: e.Translate(transl),
			})
		}

		return resterr.NewBadRequestError("invalid field values", errorCauses...)
	} else {
		return resterr.NewBadRequestError("error trying to convert fields")
	}
}
