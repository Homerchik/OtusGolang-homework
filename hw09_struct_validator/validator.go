package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrIntLessThanMin    = errors.New("less than minimum value")
	ErrIntGreaterThanMax = errors.New("greater than maximum value")
	ErrIntNotInList      = errors.New("not in range")
	ErrStrExceedsLen     = errors.New("exceeds maximum length")
	ErrStrNotMatchRegexp = errors.New("doesn't match regexp")
	ErrStrNotInList      = errors.New("value is not in permitted values list")
	ErrNotDefinedRule    = errors.New("unknown rule")
	ErrGeneral           = errors.New("error during processing")
)

type IntTag struct {
	rule string
	val  []int
}

type StringTag struct {
	rule, val string
}

func NewIntTag(tagline string) (*IntTag, error) {
	parsed := strings.SplitN(tagline, ":", 2)
	valuesStr := strings.Split(parsed[1], ",")
	valuesInt := make([]int, len(valuesStr))
	for i, v := range valuesStr {
		vInt, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("can't parse rule for tagline %v: error %w: %w", tagline, err, ErrGeneral)
		}
		valuesInt[i] = vInt
	}
	return &IntTag{rule: parsed[0], val: valuesInt}, nil
}

func (t *IntTag) Validate(v int) error {
	switch t.rule {
	case "min":
		if v < t.val[0] {
			return fmt.Errorf("value %v, minimum %v: %w", v, t.val[0], ErrIntLessThanMin)
		}
	case "max":
		if v > t.val[0] {
			return fmt.Errorf("value %v, maximum %v: %w", v, t.val[0], ErrIntGreaterThanMax)
		}
	case "in":
		for _, val := range t.val {
			if v == val {
				return nil
			}
		}
		return fmt.Errorf("value %v, range %v: %w", v, t.val, ErrIntNotInList)
	default:
		return fmt.Errorf("rule %v is not defined: %w", t.rule, ErrNotDefinedRule)
	}
	return nil
}

func NewStringTag(tagline string) *StringTag {
	parsed := strings.SplitN(tagline, ":", 2)
	return &StringTag{rule: parsed[0], val: parsed[1]}
}

func (t *StringTag) Validate(v string) error {
	switch t.rule {
	case "len":
		return lenValidate(v, t.val)
	case "regexp":
		return regexpValidate(v, t.val)
	case "in":
		return isPartOfListValidate(v, t.val)
	default:
		return fmt.Errorf("rule %v is not defined: %w", t.rule, ErrNotDefinedRule)
	}
}

func lenValidate(s string, expLen string) error {
	intLen, err := strconv.Atoi(expLen)
	if err != nil {
		return fmt.Errorf("can't parse rule value %v: %w: %w", expLen, err, ErrGeneral)
	}
	if len(s) > intLen {
		return fmt.Errorf("value %v, max length %v: %w", s, expLen, ErrStrExceedsLen)
	}
	return nil
}

func regexpValidate(s, r string) error {
	if match, err := regexp.Match(r, []byte(s)); err != nil {
		return fmt.Errorf("can't compile expression %v: %w", r, ErrGeneral)
	} else if !match {
		return fmt.Errorf("value %v, expression %v: %w", s, r, ErrStrNotMatchRegexp)
	}
	return nil
}

func isPartOfListValidate(s, l string) error {
	if !strings.Contains(l, s) {
		return fmt.Errorf("value %v, list %v: %w", s, l, ErrStrNotInList)
	}
	return nil
}

func (v ValidationErrors) Error() string {
	sb := strings.Builder{}
	for _, err := range v {
		sb.WriteString(fmt.Sprintf("Field:\n%v\nErrors:\n%v\n\n", err.Field, err.Err))
	}
	return sb.String()
}

func Validate(v interface{}) error {
	var errs error
	st := reflect.ValueOf(v)
	stTyp := reflect.TypeOf(v)
	if st.Kind() != reflect.Struct {
		return nil
	}
	for i := range st.NumField() {
		var err error
		stValue := st.FieldByIndex([]int{i})
		stField := stTyp.Field(i)
		tagline := stField.Tag.Get("validate")
		if tagline == "" {
			continue
		}
		fieldName := stField.Name
		//nolint:exhaustive
		switch stValue.Kind() {
		case reflect.String, reflect.Int:
			err = ValidateValue(stValue, fieldName, tagline)
		case reflect.Slice:
			err = ValidateSlice(stValue, fieldName, tagline)
		default:
			continue
		}
		if err != nil {
			errs = errors.Join(err, errs)
		}
	}
	return errs
}

func ValidateValue(v reflect.Value, name, tagline string) error {
	errs := []error{}
	tagsStr := strings.Split(tagline, "|")
	for _, tagStr := range tagsStr {
		if v.Kind() == reflect.Int {
			if tag, err := NewIntTag(tagStr); err != nil {
				errs = append(errs, err)
			} else if err := tag.Validate(int(v.Int())); err != nil {
				errs = append(errs, err)
			}
		} else if err := NewStringTag(tagStr).Validate(v.String()); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return ValidationErrors{ValidationError{name, errors.Join(errs...)}}
	}
	return nil
}

func ValidateSlice(v reflect.Value, name, tagline string) error {
	var errs error
	for i := range v.Len() {
		item := v.Index(i)
		if err := ValidateValue(item, name, tagline); err != nil {
			errs = errors.Join(err, errs)
		}
	}
	return errs
}
