package gin

import (
	"Advanced_Shop/app/pkg/code"
	"Advanced_Shop/pkg/errors"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"strings"
)

func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// formatValidationErrors 将验证错误格式化为自然语言
func formatValidationErrors(errMap map[string]string) string {
	if len(errMap) == 0 {
		return "validation failed"
	}

	var messages []string
	for field, msg := range errMap {
		messages = append(messages, field+": "+msg)
	}

	return strings.Join(messages, "; ")
}

func HandleValidatorError(c *gin.Context, err error, trans ut.Translator) error {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors.WithCode(code.ErrValidationTranslate, err.Error())
	}

	translatedErrors := removeTopStruct(errs.Translate(trans))
	return errors.WithCode(code.ErrValidation, formatValidationErrors(translatedErrors))
}
