package repl

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

var errEmptyInput = errors.New("empty input")

func parseInput(line string) (any, error) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return nil, errEmptyInput
	}

	if strings.HasPrefix(trimmed, "[") || strings.HasPrefix(trimmed, "{") {
		value, ok, err := parseJSONValue(trimmed)
		if err != nil {
			return nil, err
		}
		if ok {
			return value, nil
		}
	}

	op, a, b, ok, err := parseBinaryExpression(trimmed)
	if err != nil {
		return nil, err
	}
	if ok {
		return []any{op, a, b}, nil
	}

	scalar, ok, err := parseScalar(trimmed)
	if err != nil {
		return nil, err
	}
	if ok {
		return scalar, nil
	}

	return nil, errors.New("unrecognized input")
}

func parseJSONValue(input string) (any, bool, error) {
	decoder := json.NewDecoder(strings.NewReader(input))
	decoder.UseNumber()

	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, false, err
	}
	if decoder.More() {
		return nil, false, errors.New("unexpected trailing JSON tokens")
	}

	normalized := normalizeJSONValue(value)
	if op, left, right, ok := normalizeBinaryOperands(normalized); ok {
		return []any{op, left, right}, true, nil
	}

	return normalized, true, nil
}

func parseBinaryExpression(input string) (string, any, any, bool, error) {
	trimmed := strings.TrimSpace(input)
	trimmed = trimEnclosing(trimmed, '(', ')')
	trimmed = trimEnclosing(trimmed, '[', ']')

	parts := strings.Fields(trimmed)
	if len(parts) != 3 {
		return "", nil, nil, false, nil
	}

	op := parts[0]
	if !isBinaryOperator(op) {
		return "", nil, nil, false, nil
	}

	left, ok, err := parseScalar(parts[1])
	if err != nil {
		return "", nil, nil, false, err
	}
	if !ok {
		return "", nil, nil, false, errors.New("invalid left operand")
	}

	right, ok, err := parseScalar(parts[2])
	if err != nil {
		return "", nil, nil, false, err
	}
	if !ok {
		return "", nil, nil, false, errors.New("invalid right operand")
	}

	return op, left, right, true, nil
}

func parseScalar(input string) (any, bool, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return nil, false, nil
	}

	if isQuotedString(trimmed) {
		return trimmed, true, nil
	}

	switch trimmed {
	case "true":
		return true, true, nil
	case "false":
		return false, true, nil
	case "nil", "null":
		return nil, true, nil
	}

	if i, err := strconv.Atoi(trimmed); err == nil {
		return i, true, nil
	}

	if f, err := strconv.ParseFloat(trimmed, 64); err == nil {
		return f, true, nil
	}

	return nil, false, nil
}

func normalizeJSONValue(value any) any {
	switch v := value.(type) {
	case json.Number:
		if strings.Contains(v.String(), ".") {
			if f, err := v.Float64(); err == nil {
				return f
			}
		} else if i, err := v.Int64(); err == nil {
			return int(i)
		}
		if f, err := v.Float64(); err == nil {
			return f
		}
		return v.String()
	case []any:
		out := make([]any, len(v))
		for i, item := range v {
			out[i] = normalizeJSONValue(item)
		}
		return out
	case map[string]any:
		out := make(map[string]any, len(v))
		for key, item := range v {
			out[key] = normalizeJSONValue(item)
		}
		return out
	default:
		return v
	}
}

func normalizeBinaryOperands(value any) (string, any, any, bool) {
	expr, ok := value.([]any)
	if !ok || len(expr) != 3 {
		return "", nil, nil, false
	}
	op, ok := expr[0].(string)
	if !ok || !isBinaryOperator(op) {
		return "", nil, nil, false
	}

	left := expr[1]
	right := expr[2]
	if s, ok := left.(string); ok {
		left = `"` + s + `"`
	}
	if s, ok := right.(string); ok {
		right = `"` + s + `"`
	}

	return op, left, right, true
}

func isBinaryOperator(op string) bool {
	switch op {
	case "+", "-", "*", "/", "%":
		return true
	default:
		return false
	}
}

func isQuotedString(s string) bool {
	if len(s) < 2 {
		return false
	}
	return (s[0] == '\'' && s[len(s)-1] == '\'') || (s[0] == '"' && s[len(s)-1] == '"')
}

func trimEnclosing(s string, open, close byte) string {
	if len(s) >= 2 && s[0] == open && s[len(s)-1] == close {
		return strings.TrimSpace(s[1 : len(s)-1])
	}
	return s
}
