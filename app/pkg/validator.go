package pkg

import (
	"net/http"
	"reflect"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func ValidateRequest[T any](request T) (map[string]string, error) {
	errs := validator.New().Struct(request)
	msgs := make(map[string]string)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			field, _ := reflect.TypeOf(request).Elem().FieldByName(err.Field())
			queryTag := getStructTag(field, "query")
			message := getStructTag(field, "message")
			msgs[queryTag] = message
		}
		return msgs, errs
	}
	return nil, nil
}

// getStructTag returns the value of the tag with the given name
func getStructTag(f reflect.StructField, tagName string) string {
	return string(f.Tag.Get(tagName))
}

func BindRequest[T any](c echo.Context, request T) error {
	err := c.Bind(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}
