package validation

import (
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"regexp"
)

func RegisterMobile(translator ut.Translator) {
	// Gin 框架默认使用 go-playground/validator 这个库来做参数验证。  Engine : 返回底层实际的验证器引擎
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if ok { // Engine() 返回的是一个 interface{} 类型，我们需要把它转换成具体的 *validator.Validate 类型，才能调用它的方法

		//规定签名必须是这个  func(fl validator.FieldLevel) bool
		_ = validate.RegisterValidation("mobile", ValidateMobile)
	}
}

func ValidateMobile(f1 validator.FieldLevel) bool {
	mobile := f1.Field().String()
	// 省略错误
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	if ok {
		return true
	}
	return false

}
