package helper

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/wolftotem4/golava-core/instance"
	"github.com/wolftotem4/golava-core/lang"
)

func GetTranslator(i *instance.Instance, soft bool) ut.Translator {
	options := []lang.TranslatorOption{lang.Fallback(i.App.Base().Translation.GetFallback())}

	if soft {
		options = append(options, lang.Soft)
	}

	return i.GetUserPreferredTranslator(options...)
}
