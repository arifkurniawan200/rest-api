package app

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
	"strconv"
	"strings"
)

type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator membuat instance validator kustom baru
func NewCustomValidator() *CustomValidator {
	var (
		v   = validator.New()
		err error
	)

	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() error {
		err = v.RegisterValidation("blacklistWords", notContainsValidation)
		return err
	})

	g.Go(func() error {
		err = v.RegisterValidation("checkCategory", allowedCategory)
		return err
	})

	g.Go(func() error {
		err = v.RegisterValidation("inRange", inRange)
		return err
	})

	err = g.Wait()
	if err != nil {
		return nil
	}

	return &CustomValidator{v}
}

// Validate validate using custom validation
func (cv *CustomValidator) Validate(c echo.Context, i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

// notContainsValidation custom function to block several words
func notContainsValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	invalidWords := []string{"sex", "gay", "lesbian"}

	for _, word := range invalidWords {
		if strings.Contains(strings.ToLower(value), word) {
			return false
		}
	}
	return true
}

// notContainsValidation custom function to block several words
func allowedCategory(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	invalidWords := []string{"photo", "sketch", "cartoon", "animation"}

	for _, word := range invalidWords {
		if strings.EqualFold(value, word) {
			return true
		}
	}
	return false
}

func inRange(fl validator.FieldLevel) bool {
	param := fl.Param()
	params := strings.Split(param, "-")

	if len(params) != 2 {
		return false
	}

	v1, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		return false
	}

	v2, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return false
	}

	value := fl.Field().Int()

	return value >= v1 && value <= v2
}

func _FormatValidationError(err error) []string {
	var messages []string
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("%s is required", err.Field()))
		case "blacklistWords":
			messages = append(messages, fmt.Sprintf("%s must not contain words %s", err.Field(), "sex, gay or lesbian"))
		case "url":
			messages = append(messages, fmt.Sprintf("%s is not a valid URL", err.Field()))
		case "email":
			messages = append(messages, fmt.Sprintf("%s is not a valid email", err.Field()))
		case "checkCategory":
			messages = append(messages, fmt.Sprintf("%s allowed categories are %s", err.Field(), "photo, sketch, cartoon or animation"))
		case "min":
			messages = append(messages, fmt.Sprintf("%s must be greater than %s character (minimum %s char)", err.Field(), err.Param(), err.Param()))
		case "max":
			messages = append(messages, fmt.Sprintf("%s must be less than %s character (maximum %s char)", err.Field(), err.Param(), err.Param()))
		case "inRange":
			param := err.Param()
			params := strings.Split(param, "-")
			if len(params) != 2 {
				continue
			}
			v1, _ := strconv.ParseInt(params[0], 10, 64)
			v2, _ := strconv.ParseInt(params[1], 10, 64)
			messages = append(messages, fmt.Sprintf("%s must be less than %v and greater than %v", err.Field(), v2, v1))
		default:
			messages = append(messages, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return messages
}
