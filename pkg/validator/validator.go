package validator

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	trans    ut.Translator
	validate *validator.Validate
}

func NewValidator() Validator {
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	return Validator{
		validate: validate,
		trans:    trans,
	}
}

type ValidationErrors map[string]string

func (ve ValidationErrors) Error() string {
	j, err := json.Marshal(ve)
	if err != nil {
		return fmt.Sprintf("failed to marshal validation errors: %v", err)
	}
	return string(j)
}

func (v *Validator) Validate(dto interface{}) *ValidationErrors {
	err := v.validate.Struct(dto)
	if err != nil {

		if errs, ok := err.(validator.ValidationErrors); ok {

			translations := errs.Translate(v.trans)

			convertedErrors := ValidationErrors{}
			for field, message := range translations {
				convertedErrors[field] = message
			}

			return &convertedErrors
		}

		return nil
	}

	return nil
}

func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	regex := regexp.MustCompile(pattern)

	return regex.MatchString(email)
}

// func ValidatePassword(password string) bool {
// 	allowedSymbols := "!@#$%^&*()-_+="

// 	containsLowercase := regexp.MustCompile(`[a-z]`).MatchString(password)
// 	containsUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
// 	containsDigit := regexp.MustCompile(`\d`).MatchString(password)

// 	var containsSymbol bool
// 	for _, char := range password {
// 		if strings.ContainsRune(allowedSymbols, char) {
// 			containsSymbol = true
// 			break
// 		}
// 	}

// 	return containsLowercase && containsUppercase && containsDigit && containsSymbol && len(password) >= 8
// }

// func ValidateRequestRegister(ctx context.Context, req dto.RegisterRequest) error {
// 	traceID := helper.GetTraceID(ctx)
// 	switch {
// 	case req.Email == "" || !ValidateEmail(req.Email):
// 		return response.ErrInvalidEmail.WithTraceID(traceID)
// 	case req.Password == "" || !ValidatePassword(req.Password):
// 		return response.ErrInvalidPassword.WithTraceID(traceID)
// 	}
// 	return nil
// }
