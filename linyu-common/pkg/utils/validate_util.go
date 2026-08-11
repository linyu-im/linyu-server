package utils

import (
	"errors"
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
)

var (
	Mobile   = "mobile"
	Email    = "email"
	Account  = "account"
	Password = "password"
	Eqfield  = "eqfield"
)

var ValidationErrorMessages = map[string]string{
	Eqfield: "param.password-not-match",
}

func init() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}
	_ = RegisterRegexValidation(v, Mobile, `^1[3-9]\d{9}$`, "param.phone-format-error")
	_ = RegisterRegexValidation(v, Email, `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, "param.email-format-error")
	_ = RegisterRegexValidation(v, Account, `^[a-zA-Z][a-zA-Z0-9_-]{3,19}$`, "param.account-format-error")
	_ = RegisterAllRegexValidation(v, Password, "param.password-format-error",
		`^[\x21-\x7E]{8,20}$`,
		`[A-Za-z]`,
		`[0-9]`,
		`[^A-Za-z0-9]`,
	)
}

func ShouldBindBodyWithJSONAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			var errsMsg []string
			t := reflect.TypeOf(obj).Elem()
			for _, e := range validationErrors {
				field, _ := t.FieldByName(e.StructField())
				if msg, ok := ValidationErrorMessages[e.Tag()]; ok {
					response.Fail(c, msg)
					return false
				} else {
					errsMsg = append(errsMsg, field.Tag.Get("json"))
				}
			}
			response.FailWithErrData(c, "param.validate-failed", map[string]interface{}{
				"errors": errsMsg,
			})
			return false
		}
		response.Fail(c, err.Error())
		return false
	}
	return true
}

// RegisterRegexValidation 注册基于正则的通用校验器，并绑定对应报错文案
func RegisterRegexValidation(v *validator.Validate, tag, pattern, errMsg string) error {
	return RegisterAllRegexValidation(v, tag, errMsg, pattern)
}

// RegisterAllRegexValidation 注册多正则校验器（全部命中才通过），并绑定报错文案
func RegisterAllRegexValidation(v *validator.Validate, tag, errMsg string, patterns ...string) error {
	regs := make([]*regexp.Regexp, 0, len(patterns))
	for _, pattern := range patterns {
		reg, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}
		regs = append(regs, reg)
	}
	if errMsg != "" {
		ValidationErrorMessages[tag] = errMsg
	}
	return v.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		for _, reg := range regs {
			if !reg.MatchString(val) {
				return false
			}
		}
		return true
	})
}
