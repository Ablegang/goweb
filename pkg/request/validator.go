package request

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/zh"
)

func Zh(obj interface{}) error {

	validate := validator.New()
	trans, _ := ut.New(zh.New()).GetTranslator("zh")
	_ = translations.RegisterDefaultTranslations(validate, trans)

	if err := validate.Struct(obj); err != nil {
		errs, _ := err.(validator.ValidationErrors)
		for _, e := range errs {
			return errors.New(e.Translate(trans))
		}
	}

	return nil
}

// 绑定参数
func Bind(c *gin.Context, obj interface{}) error {
	err := c.ShouldBind(obj)
	if err != nil {
		return errors.New("请求错误")
	}

	err = Zh(obj)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
