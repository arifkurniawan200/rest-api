package app

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"
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
		if !strings.EqualFold(value, word) {
			return false
		}
	}
	return true
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
		default:
			messages = append(messages, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return messages
}
