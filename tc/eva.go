package tc

import (
	"errors"
	"fmt"
)

var (
	ErrCheckingType           = func(msg error) error { return fmt.Errorf("typechecking error: %w", msg) }
	ErrInvalidBinaryOperation = errors.New("invalid binary operation")
	ErrInvalidType            = func(t any, e any) error { return fmt.Errorf("invalid type: '%s', expected '%s'", t, e) }
	ErrInvalidOperation       = func(op Operation) error { return fmt.Errorf("invalid operation: '%s'", op.String()) }
	ErrInvalidStringLiteral   = func(s string) error { return fmt.Errorf("invalid string literal: '%s'", s) }
)

type Type string

func (t Type) String() string {
	return string(t)
}

func ToType(s string) Type {
	return Type(s)
}

var (
	Number            Type = "number"
	SingleQuoteString Type = "single-quote-string"
	DoubleQuoteString Type = "double-quote-string"
	Boolean           Type = "boolean"
	Array             Type = "array"
	Object            Type = "object"
	Null              Type = "nil"
	Function          Type = "function"
)

type Operation string

func (t Operation) String() string {
	return string(t)
}

func ToOperation(s string) Operation {
	return Operation(s)
}

var (
	Addition       Operation = "+"
	Subtraction    Operation = "-"
	Multiplication Operation = "*"
	Division       Operation = "/"
	Modulo         Operation = "%"
)

var operationPossibilities = map[Operation][]Type{
	Addition:       []Type{Number, DoubleQuoteString, SingleQuoteString, Array},
	Subtraction:    []Type{Number, Array},
	Multiplication: []Type{Number, Array},
	Division:       []Type{Number},
	Modulo:         []Type{Number},
}

// Eva is the typechecker
type Eva struct{}

func New() *Eva {
	return &Eva{}
}

func (e *Eva) Check(expr any) (Type, error) {
	var result Type
	switch v := expr.(type) {
	case int:
		result = Number
	case string:
		result = e.detectStringType(v)
		if result == "" {
			return "", ErrInvalidStringLiteral(v)
		}
	case float64:
		result = Number
	case bool:
		result = Boolean
	case []any:
		if e.isBinaryOperation(v) {
			// + 1 2
			// validate possible operations for each operator
			operator := ToOperation(v[0].(string))
			possibleTypes, ok := operationPossibilities[operator]
			if !ok {
				return "", ErrInvalidOperation(operator)
			}

			// get the two parameters
			param1 := v[1]
			param2 := v[2]

			// perform type checking for each parameter recursively

			// check parameter 1
			result1, err := e.Check(param1)
			if err != nil {
				return "", ErrCheckingType(err)
			}

			// check parameter 2
			result2, err := e.Check(param2)
			if err != nil {
				return "", ErrCheckingType(err)
			}

			ok = e.ensureParametersTypeAreOkForOperation(possibleTypes, result1, result2)
			if !ok {
				return "", ErrInvalidBinaryOperation
			}

			// we cannot perform binary operations on different types, the second parameter must be the same type as the first parameter
			if result1 != result2 {
				return "", ErrInvalidType(result2, result1)
			}

			// the return type of any binary operation is detected from type of the first parameter
			result = result1
		} else {
			result = Array
		}
	case map[string]any:
		result = Object
	case func(...any) any:
		result = Function
	case nil:
		result = Null
	default:
		return "", ErrCheckingType(fmt.Errorf("invalid type provided: %T", v))
	}
	return result, nil
}

func (e *Eva) ensureParametersTypeAreOkForOperation(possibleTypes []Type, t1, t2 Type) bool {
	for _, t := range possibleTypes {
		if t == t1 || t == t2 {
			return true
		}
	}

	return false
}

func (e *Eva) isBinaryOperation(expr []any) bool {
	// check the length of the expression
	if len(expr) != 3 {
		return false
	}
	// get the operator
	op, ok := expr[0].(string)
	if !ok {
		return false
	}

	// check if the operator is a binary operator
	switch op {
	case "+", "-", "*", "/", "%":
		return true
	default:
		return false
	}
}

func (e *Eva) detectStringType(s string) Type {
	if len(s) == 0 {
		return DoubleQuoteString
	} else if s[0] == '\'' && s[len(s)-1] == '\'' {
		return SingleQuoteString
	} else if s[0] == '"' && s[len(s)-1] == '"' {
		return DoubleQuoteString
	}
	return ""
}
