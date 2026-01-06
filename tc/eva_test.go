package tc

import "testing"

func TestEva_CheckType(t *testing.T) {
	tests := []struct {
		name         string
		input        any
		wantErr      bool
		expectedType Type
	}{
		{
			name:         "1 - Number Type",
			input:        1,
			expectedType: Number,
		},
		{
			name:         "'hello' - String Type",
			input:        `'hello'`,
			expectedType: SingleQuoteString,
		},
		{
			name:         "\"hello\" - String Type",
			input:        `"hello"`,
			expectedType: DoubleQuoteString,
		},
		{
			name:         "'' - String Type",
			input:        `''`,
			expectedType: SingleQuoteString,
		},
		{
			name:         "true - Boolean Type",
			input:        true,
			expectedType: Boolean,
		},
		{
			name:         "nil - Nil Type",
			input:        nil,
			expectedType: Null,
		},
		{
			name:         "empty string - String Type",
			input:        "",
			expectedType: DoubleQuoteString,
		},
		{
			name:         "Addition on numbers - binary operations",
			input:        []any{"+", 1, 9},
			expectedType: Number,
		},
		{
			name:         "Addition on strings - binary operations",
			input:        []any{"+", `"hello"`, `"world"`},
			expectedType: DoubleQuoteString,
		},
		{
			name:         "Subtraction on numbers - binary operations",
			input:        []any{"-", 1, 9},
			expectedType: Number,
		},
		{
			name:         "Multiplication on numbers - binary operations",
			input:        []any{"*", 1, 9},
			expectedType: Number,
		},
		{
			name:         "Division on numbers - binary operations",
			input:        []any{"/", 1, 9},
			expectedType: Number,
		},
		{
			name:         "Modulo on numbers - binary operations",
			input:        []any{"%", 1, 9},
			expectedType: Number,
		},
		{
			name:    "Subtraction on strings - binary operations",
			input:   []any{"-", `"hello"`, `"world"`},
			wantErr: true,
		},
		{
			name:    "Multiplication on strings - binary operations",
			input:   []any{"*", `"hello"`, `"world"`},
			wantErr: true,
		},
		{
			name:    "Subtraction on boolean - binary operations",
			input:   []any{"-", true, false},
			wantErr: true,
		},
		{
			name:    "Modulo on strings - binary operations",
			input:   []any{"%", `"hello"`, `"world"`},
			wantErr: true,
		},
		{
			name:    "Division on strings - binary operations",
			input:   []any{"/", `"hello"`, `"world"`},
			wantErr: true,
		},
		{
			name:    "Division on boolean - binary operations",
			input:   []any{"/", true, false},
			wantErr: true,
		},
		{
			name:    "Addition on boolean - binary operations",
			input:   []any{"+", true, false},
			wantErr: true,
		},
		{
			name:         "Addition on arrays - binary operations",
			input:        []any{"+", []any{1, 2}, []any{3, 4}},
			expectedType: Array,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			eva := New()
			result, err := eva.Check(test.input)

			if test.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != test.expectedType {
				t.Errorf("expected type: %v, got=%v", test.expectedType, result)
			}
		})
	}
}
