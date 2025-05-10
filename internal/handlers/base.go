package handler

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	datastructure "tracking-service/internal/datastructures"
	errdefs "tracking-service/internal/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type BaseHandler struct{}

func (b *BaseHandler) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, datastructure.BaseResponse{
		Success: true,
		Data:    data,
	})
}

func (b *BaseHandler) SuccessWithoutData(c *gin.Context) {
	c.JSON(http.StatusNoContent, datastructure.BaseResponse{Success: true})
}

func (b *BaseHandler) SuccessWithoutContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

func (b *BaseHandler) ErrorResponse(c *gin.Context, cause error) {
	b.respondWithStatus(c, cause)
}

func (b *BaseHandler) respondWithStatus(c *gin.Context, cause error) {
	ctx := c.Request.Context()

	switch cause {
	case errdefs.ErrorNotFound:
		c.JSON(404, datastructure.ErrorResponseWithCode{
			ErrorResponse: datastructure.ErrorResponse{
				Success: false,
				Message: cause.Error(),
			},
			Details: nil,
		})
		return
	case errdefs.ErrorDuplicateKey:
		c.JSON(409, datastructure.ErrorResponseWithCode{
			ErrorResponse: datastructure.ErrorResponse{
				Success: false,
				Message: cause.Error(),
			},
			Details: nil,
		})
		return
	case errdefs.ErrorInvalidRequest:
		c.JSON(400, datastructure.ErrorResponseWithCode{
			ErrorResponse: datastructure.ErrorResponse{
				Success: false,
				Message: cause.Error(),
			},
			Details: nil,
		})
		return
	case errdefs.ErrorUnauthorized:
		c.JSON(401, datastructure.ErrorResponseWithCode{
			ErrorResponse: datastructure.ErrorResponse{
				Success: false,
				Message: cause.Error(),
			},
			Details: nil,
		})
		return
	case errdefs.ErrorForbidden:
		c.JSON(403, datastructure.ErrorResponseWithCode{
			ErrorResponse: datastructure.ErrorResponse{
				Success: false,
				Message: cause.Error(),
			},
			Details: nil,
		})
		return
	}

	log.WithContext(ctx).Errorf("Internal server error: %v", cause)
	c.JSON(500, datastructure.ErrorResponseWithCode{
		ErrorResponse: datastructure.ErrorResponse{
			Success: false,
			Message: cause.Error(),
		},
		Details: nil,
	})
}

func (b *BaseHandler) InvalidInputErrorResponse(c *gin.Context, cause error) {
	err := errdefs.ErrorInvalidRequest
	var ve validator.ValidationErrors
	if !errors.As(cause, &ve) {
		log.WithContext(c).Errorf("Invalid input binding error: %v", cause)
		b.BadRequest(c)
		return
	}

	errorsMap := make(map[string]string)
	for _, err := range ve {
		fieldName := getJSONFieldName(err)
		msg := mapValidationErrorToMessage(err)
		errorsMap[fieldName] = msg

		log.WithContext(c).Errorf("Validation error: field=%s tag=%s param=%s",
			err.Field(), err.Tag(), err.Param(),
		)
	}

	c.JSON(http.StatusBadRequest, datastructure.ErrorResponseWithCode{
		ErrorResponse: datastructure.ErrorResponse{
			Success: false,
			Message: err.Error(),
		},
		Details: errorsMap,
	})
}

func getJSONFieldName(err validator.FieldError) string {
	if t, ok := reflect.TypeOf(datastructure.Tenant{}).FieldByName(err.Field()); ok {
		tag := t.Tag.Get("json")
		if tag != "" && tag != "-" {
			return strings.Split(tag, ",")[0]
		}
	}
	return strings.ToLower(err.Field())
}

func (b *BaseHandler) BadRequest(c *gin.Context) {
	err := errdefs.ErrorInvalidRequest
	c.JSON(http.StatusBadRequest, datastructure.ErrorResponse{
		Success: false,
		Message: err.Error(),
	})
}

func (b *BaseHandler) InternalServerError(c *gin.Context) {
	err := errdefs.ErrorInternalError
	c.JSON(http.StatusInternalServerError, datastructure.ErrorResponse{
		Success: false,
		Message: err.Error(),
	})
}

func mapValidationErrorToMessage(err validator.FieldError) string {
	fieldName := err.Field()
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fieldName)
	case "email":
		return fmt.Sprintf("%s is not a valid email", fieldName)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fieldName, err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fieldName, err.Param())
	case "gtfield":
		return fmt.Sprintf("%s must be greater than %s", fieldName, err.Param())
	case "ltfield":
		return fmt.Sprintf("%s must be less than %s", fieldName, err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of %s", fieldName, err.Param())
	case "eqfield":
		return fmt.Sprintf("%s must be equal to %s", fieldName, err.Param())
	case "required_without", "required_without_all":
		return fmt.Sprintf("%s is required when %s is not present", fieldName, err.Param())
	case "required_with":
		return fmt.Sprintf("%s is required when %s is present", fieldName, err.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fieldName, err.Param())
	case "dateFormat":
		return fmt.Sprintf("%s is invalid date format", fieldName)
	case "datetimeFormat":
		return fmt.Sprintf("%s is invalid datetime format", fieldName)
	case "futureDateOnly":
		return fmt.Sprintf("%s must be in the future", fieldName)
	case "mobile":
		return fmt.Sprintf("%s is invalid mobile phone", fieldName)
	default:
		return fmt.Sprintf("%s is invalid", fieldName)
	}
}
