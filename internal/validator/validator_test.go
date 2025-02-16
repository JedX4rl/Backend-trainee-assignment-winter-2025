package structValidator

import (
	"testing"
)

type User struct {
	Username string `json:"username" validate:"required,min=3,max=10"`
	Age      int    `json:"age" validate:"required,min=18"`
	Email    string `json:"email" validate:"required,email"`
}

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		expectedError string
	}{
		{
			name: "Valid User",
			input: &User{
				Username: "testuser",
				Age:      25,
				Email:    "test@example.com",
			},
			expectedError: "",
		},
		{
			name: "Invalid User - Too short username",
			input: &User{
				Username: "t",
				Age:      25,
				Email:    "test@example.com",
			},
			expectedError: "field 'Username' validation failed on the 'min' tag",
		},
		{
			name: "Invalid User - Age less than 18",
			input: &User{
				Username: "testuser",
				Age:      16,
				Email:    "test@example.com",
			},
			expectedError: "field 'Age' validation failed on the 'min' tag",
		},
		{
			name: "Invalid User - Invalid email",
			input: &User{
				Username: "testuser",
				Age:      25,
				Email:    "invalid-email",
			},
			expectedError: "field 'Email' validation failed on the 'email' tag",
		},
		{
			name: "Invalid User - Missing fields",
			input: &User{
				Username: "",
				Age:      0,
				Email:    "",
			},
			expectedError: "field 'Username' validation failed on the 'required' tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.input)
			if err != nil {
				if err.Error() != tt.expectedError {
					t.Errorf("expected error '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else if tt.expectedError != "" {
				t.Errorf("expected error '%s', but got nil", tt.expectedError)
			}
		})
	}
}
