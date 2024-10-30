package validations

import (
	"encoding/json"
	"errors"
	"ticket-booking/configs/errs"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslation "github.com/go-playground/validator/v10/translations/en"
)

var (
	_      = validator.New()
	transl ut.Translator
)

func init() {
	if val, ok := binding.Validator.Engine().(*validator.Validate); ok {
		translator := en.New()
		unt := ut.New(translator, translator)
		transl, _ = unt.GetTranslator("en")
		_ = entranslation.RegisterDefaultTranslations(val, transl)
	}
}

func ValidateRequest(validationErr error) *errs.Error {
	var jsonErr *json.UnmarshalTypeError
	var validationErrors validator.ValidationErrors

	if errors.As(validationErr, &jsonErr) {
		return errs.NewBadRequest("Invalid field type")
	} else if errors.As(validationErr, &validationErrors) {
		errorCauses := make([]*errs.Cause, 0)

		for _, e := range validationErrors {
			cause := &errs.Cause{
				Message: e.Translate(transl),
				Field:   e.Field(),
			}
			errorCauses = append(errorCauses, cause)
		}

		return errs.NewValidationError("Some fields are invalid", errorCauses)
	} else {
		return errs.NewBadRequest("Error trying to convert fields")
	}
}
