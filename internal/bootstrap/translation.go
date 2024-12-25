package bootstrap

import (
	"reflect"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/wolftotem4/golangloader"
	"github.com/wolftotem4/golangloader/defaults"
	"github.com/wolftotem4/golava-core/validation"
)

var validatorReplaced = false

func initTranslation(fallback string) (*ut.UniversalTranslator, error) {
	var validate *validator.Validate

	if validatorReplaced {
		validate = binding.Validator.Engine().(*validator.Validate)
	} else {
		validate = validator.New(validator.WithRequiredStructEnabled())
		validate.SetTagName("binding")
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			return fld.Tag.Get("json")
		})
		binding.Validator = validation.NewMoldModifyValidator(validate)

		validatorReplaced = true
	}

	uni, err := defaults.Locales.NewUniversalTranslator(fallback)
	if err != nil {
		return nil, err
	}

	err = golangloader.LoadTranslate("lang", uni)
	if err != nil {
		return nil, err
	}

	err = uni.VerifyTranslations()
	if err != nil {
		return nil, err
	}

	err = defaults.Locales.RegisterValidateTranslation(validate, uni)
	if err != nil {
		return nil, err
	}

	return uni, nil
}
