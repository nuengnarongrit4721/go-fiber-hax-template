package validation

import (
	"fmt"
	"net/mail"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	sharederrs "gofiber-hax/internal/shared/errors"

	"github.com/gofiber/fiber/v2"
)

type FieldError struct {
	Field string `json:"field"`
	Rule  string `json:"rule"`
}

type Error struct {
	Fields []FieldError `json:"fields"`
}

func (e Error) Error() string {
	return sharederrs.ErrInvalidInput.Error()
}

func (e Error) Unwrap() error {
	return sharederrs.ErrInvalidInput
}

func BindAndValidate(c *fiber.Ctx, dst any) error {
	if err := c.BodyParser(dst); err != nil {
		return Error{}
	}
	return ValidateStruct(dst)
}

func ValidateStruct(input any) error {
	if input == nil {
		return Error{}
	}

	value := reflect.ValueOf(input)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return Error{}
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("%w: validator expects struct input", sharederrs.ErrInvalidInput)
	}

	var errs []FieldError
	valueType := value.Type()
	for i := range valueType.NumField() {
		fieldType := valueType.Field(i)
		rules := strings.TrimSpace(fieldType.Tag.Get("validate"))
		if rules == "" || rules == "-" {
			continue
		}

		fieldValue := value.Field(i)
		for _, rule := range strings.Split(rules, ",") {
			rule = strings.TrimSpace(rule)
			if rule == "" {
				continue
			}
			if err := validateRule(value, fieldType, fieldValue, rule); err != nil {
				errs = append(errs, *err)
			}
		}
	}

	if len(errs) > 0 {
		return Error{Fields: errs}
	}
	return nil
}

func validateRule(parent reflect.Value, fieldType reflect.StructField, fieldValue reflect.Value, rule string) *FieldError {
	name := jsonName(fieldType)
	stringValue := fieldStringValue(fieldValue)

	switch {
	case rule == "required":
		if strings.TrimSpace(stringValue) == "" {
			return &FieldError{Field: name, Rule: "required"}
		}
	case strings.HasPrefix(rule, "min="):
		min, err := strconv.Atoi(strings.TrimPrefix(rule, "min="))
		if err != nil {
			return &FieldError{Field: name, Rule: "min"}
		}
		if utf8.RuneCountInString(strings.TrimSpace(stringValue)) < min {
			return &FieldError{Field: name, Rule: fmt.Sprintf("min=%d", min)}
		}
	case rule == "email":
		if stringValue == "" {
			return nil
		}
		addr, err := mail.ParseAddress(stringValue)
		if err != nil || addr.Address != stringValue {
			return &FieldError{Field: name, Rule: "email"}
		}
	case strings.HasPrefix(rule, "eqfield="):
		otherName := strings.TrimSpace(strings.TrimPrefix(rule, "eqfield="))
		otherField := parent.FieldByName(otherName)
		if !otherField.IsValid() || stringValue != fieldStringValue(otherField) {
			return &FieldError{Field: name, Rule: "eqfield"}
		}
	}

	return nil
}

func fieldStringValue(value reflect.Value) string {
	for value.IsValid() && value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return ""
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return ""
	}
	if value.Kind() == reflect.String {
		return value.String()
	}
	return fmt.Sprint(value.Interface())
}

func jsonName(field reflect.StructField) string {
	tag := strings.Split(field.Tag.Get("json"), ",")[0]
	if tag != "" && tag != "-" {
		return tag
	}
	return strings.ToLower(field.Name)
}
