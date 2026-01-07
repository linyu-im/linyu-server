package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/linyu-im/linyu-server/linyu-common/pkg/response"
	"reflect"
	"regexp"
)

var (
	Mobile = "mobile"
	Email  = "email"
)

var ValidationErrorMessages = map[string]string{
	Mobile: "param.phone-format-error",
	Email:  "param.email-format-error",
}

func init() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}
	mobileValidator(v)
	emailValidator(v)
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

// 手机号验证
func mobileValidator(v *validator.Validate) {
	_ = v.RegisterValidation(Mobile, func(fl validator.FieldLevel) bool {
		reg := regexp.MustCompile(`^1[3-9]\d{9}$`)
		return reg.MatchString(fl.Field().String())
	})
}

// 邮箱验证
func emailValidator(v *validator.Validate) {
	_ = v.RegisterValidation(Email, func(fl validator.FieldLevel) bool {
		reg := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		return reg.MatchString(fl.Field().String())
	})
}
