package component

import (
	"regexp"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		// 如果取得失敗，建立新實例（備援）
		v = validator.New()
	}

	// 註冊 regexp 驗證器
	_ = v.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
		re, err := regexp.Compile(fl.Param())
		if err != nil {
			return false
		}
		return re.MatchString(fl.Field().String())
	})

	// 註冊 datetime_format 驗證器
	_ = v.RegisterValidation("datetime_format", datetimeFormat)

	return v
}

func datetimeFormat(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	if dateStr == "" {
		return true // 空字串交給其他 tag 處理 required
	}
	_, err := time.Parse("2006-01-02 15:04:05", dateStr)
	return err == nil
}
