package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	TestCase struct {
		name           string
		stToValidate   interface{}
		expectedErrors []error
	}
)

func TestValidate(t *testing.T) {
	tests := []TestCase{
		{
			name:           "No errors",
			stToValidate:   User{"1", "Mike", 22, "a@a.com", "admin", []string{"38211192321", "38211192322"}, []byte("")},
			expectedErrors: nil,
		},
		{
			name:           "Error validating int too small",
			stToValidate:   User{"1", "Mike", 16, "a@a.com", "admin", []string{"38211192321", "38211192322"}, []byte("")},
			expectedErrors: []error{ErrIntLessThanMin},
		},
		{
			name:           "Error validating int too big",
			stToValidate:   User{"1", "Mike", 99, "a@a.com", "admin", []string{"38211192321", "38211192322"}, []byte("")},
			expectedErrors: []error{ErrIntGreaterThanMax},
		},
		{
			name:           "Error validating int is in list",
			stToValidate:   Response{204, ""},
			expectedErrors: []error{ErrIntNotInList},
		},
		{
			name:           "Error validating string length",
			stToValidate:   App{"123456"},
			expectedErrors: []error{ErrStrExceedsLen},
		},
		{
			name:           "Error validating string pattern",
			stToValidate:   User{"1", "Mike", 22, "a.com", "admin", []string{"38211192321", "38211192322"}, []byte("")},
			expectedErrors: []error{ErrStrNotMatchRegexp},
		},
		{
			name:           "Error validating string in list",
			stToValidate:   User{"1", "Mike", 22, "a@a.com", "superuser", []string{"38211192321", "38211192322"}, []byte("")},
			expectedErrors: []error{ErrStrNotInList},
		},
		{
			name:           "Error in slice",
			stToValidate:   User{"1", "Mike", 99, "a.com", "superuser", []string{"+38211192321", "38211192322"}, []byte("")},
			expectedErrors: []error{ErrStrExceedsLen},
		},
		{
			name:           "Multiple errors during validation",
			stToValidate:   User{"1", "Mike", 99, "a.com", "superuser", []string{"38211192321", "38211192322"}, []byte("")},
			expectedErrors: []error{ErrIntGreaterThanMax, ErrStrNotMatchRegexp, ErrStrNotInList},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("case %v", tt.name), func(t *testing.T) {
			tt := tt
			t.Parallel()
			if tt.expectedErrors == nil {
				require.NoError(t, Validate(tt.stToValidate))
			} else {
				actErr := Validate(tt.stToValidate)
				for _, err := range tt.expectedErrors {
					require.ErrorAs(t, actErr, &err)
				}
			}
			// checkError(t, errs, tt.expectedErrors)
			// require.Equal(t, tt.expectedErrors, errors, "Expected errors not equal to the validation ones")
			// 		// Place your code here.
			_ = tt
		})
	}
}
