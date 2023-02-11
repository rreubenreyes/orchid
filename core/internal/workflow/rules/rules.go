package rules

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Operator string
type Format string

const (
	Equals            Operator = "$eq"
	LessThan          Operator = "$lt"
	LessThanEquals    Operator = "$lte"
	GreaterThan       Operator = "$gt"
	GreaterThanEquals Operator = "$gte"
	Contains          Operator = "$in"
	And               Operator = "$and"
	Or                Operator = "$or"
	Xor               Operator = "$xor"
	Not               Operator = "$not"
	Has               Operator = "$has"

	DateTime Format = "date_time"
)

type Rule struct {
	Format   Format
	Operator Operator
	Operands []any
}

func Evaluate(r Rule) (bool, error) {
	switch r.Operator {
	case Equals:
		return evaluateEquals(r)
	case LessThan:
		return evaluateLessThan(r)
	case GreaterThan:
		return evaluateGreaterThan(r)
	case LessThanEquals:
		return evaluateLessThanEquals(r)
	case GreaterThanEquals:
		return evaluateGreaterThanEquals(r)
	case Contains:
		return evaluateContains(r)
	case And:
		return evaluateAnd(r)
	case Or:
		return evaluateOr(r)
	case Xor:
		return evaluateXor(r)
	case Not:
		return evaluateNot(r)
	case Has:
		return evaluateHas(r)
	default:
		return false, fmt.Errorf("unknown operator \"%s\"", r.Operator)
	}
}

// evaluateEquals evaluates the quality operation. This function expects
// a binary expression. If this is not the case, this function returns an error.
func evaluateEquals(r Rule) (bool, error) {
	if len(r.Operands) != 2 {
		return false, fmt.Errorf("invalid number of operands (%d) for operator %s", len(r.Operands), r.Operator)
	}
	return r.Operands[0] == r.Operands[1], nil
}

// evaluateAnd accepts a Rule and checks if all of its operands are also of type Rule.
// If so, then this function evaluates logical AND over all operands. If not,
// this function returns an error.
func evaluateAnd(r Rule) (bool, error) {
	for i, subrule := range r.Operands {
		valueType := reflect.TypeOf(subrule)
		ruleType := reflect.TypeOf(Rule{})
		if !valueType.AssignableTo(ruleType) {
			return false, fmt.Errorf("invalid comparison; cannot evaluate rule at position %d", i)
		}

		result, err := Evaluate(subrule.(Rule))
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}

	return true, nil
}

// evaluateOr accepts a Rule and checks if all of its operands are also of type Rule.
// If so, then this function evaluates logical OR over all operands. If not,
// this function returns an error.
func evaluateOr(r Rule) (bool, error) {
	for i, subrule := range r.Operands {
		valueType := reflect.TypeOf(subrule)
		ruleType := reflect.TypeOf(Rule{})
		if !valueType.AssignableTo(ruleType) {
			return false, fmt.Errorf("invalid comparison; cannot evaluate rule at position %d", i)
		}

		result, err := Evaluate(subrule.(Rule))
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}

	return false, nil
}

// evaluateXor accepts a Rule and checks if all of its operands are also of type Rule.
// If so, then this function evaluates logical XOR over all operands. If not,
// this function returns an error.
func evaluateXor(r Rule) (bool, error) {
	ok := false
	for i, subrule := range r.Operands {
		valueType := reflect.TypeOf(subrule)
		ruleType := reflect.TypeOf(Rule{})
		if !valueType.AssignableTo(ruleType) {
			return false, fmt.Errorf("invalid comparison; cannot evaluate rule at position %d", i)
		}

		result, err := Evaluate(subrule.(Rule))
		if err != nil {
			return false, fmt.Errorf("invalid comparison")
		}
		if result {
			if ok {
				return false, nil
			}
		}
		ok = result
	}

	return ok, nil
}

// evaluateNot accepts a Rule and checks if it has a single Operand of type Rule.
// If so, then this function negates the Operand. If not, this function returns an error.
func evaluateNot(r Rule) (bool, error) {
	if len(r.Operands) != 1 {
		return false, fmt.Errorf("invalid number of operands (%d) for operator %s", len(r.Operands), r.Operator)
	}

	result, err := Evaluate(r)
	if err != nil {
		return false, err
	}

	return result, nil
}

// evaluateLessThan evaluates the less than operation. This function expects
// a binary expression. If this is not the case, this function returns an error.
//
// If operands are of incompatible type, the evaluation returns false with no error.
// If operating on strings of format "date_time", then this function compares the timestamps
// represented by both values.
func evaluateLessThan(r Rule) (bool, error) {
	if len(r.Operands) != 2 {
		return false, fmt.Errorf("invalid number of operands (%d) for operator %s", len(r.Operands), r.Operator)
	}

	lh := reflect.ValueOf(r.Operands[0])
	rh := reflect.ValueOf(r.Operands[1])

	if lh.Kind() == reflect.Int && rh.Kind() == reflect.Int {
		return lh.Int() < rh.Int(), nil
	}
	if lh.Kind() == reflect.Float64 && rh.Kind() == reflect.Float64 {
		return lh.Float() < rh.Float(), nil
	}
	if lh.Kind() == reflect.String && rh.Kind() == reflect.String {
		if r.Format == "date_time" {
			timeA, err := time.Parse(time.RFC3339, lh.String())
			if err != nil {
				return false, err
			}
			timeB, err := time.Parse(time.RFC3339, rh.String())
			if err != nil {
				return false, err
			}

			return timeA.Before(timeB), nil
		}
		return lh.String() < rh.String(), nil
	}
	return false, nil
}

// evaluateGreaterThan evaluates the greater than operation. This function expects
// a binary expression. If this is not the case, this function returns an error.
//
// If operands are of incompatible type, the evaluation returns false with no error.
// If operating on strings of format "date_time", then this function compares the timestamps
// represented by both values.
func evaluateGreaterThan(r Rule) (bool, error) {
	if len(r.Operands) != 2 {
		return false, fmt.Errorf("invalid number of operands (%d) for operator %s", len(r.Operands), r.Operator)
	}

	lh := reflect.ValueOf(r.Operands[0])
	rh := reflect.ValueOf(r.Operands[1])
	if lh.Kind() == reflect.Int && rh.Kind() == reflect.Int {
		return lh.Int() > rh.Int(), nil
	}
	if lh.Kind() == reflect.Float64 && rh.Kind() == reflect.Float64 {
		return lh.Float() > rh.Float(), nil
	}
	if lh.Kind() == reflect.String && rh.Kind() == reflect.String {
		if r.Format == "date_time" {
			timeA, err := time.Parse(time.RFC3339, lh.String())
			if err != nil {
				return false, err
			}
			timeB, err := time.Parse(time.RFC3339, rh.String())
			if err != nil {
				return false, err
			}

			return timeA.After(timeB), nil
		}
		return lh.String() > rh.String(), nil
	}
	return false, nil
}

// evaluateLessThanEquals evaluates the less than or equals operation by negating
// the result of evaluating the greater than operation.
func evaluateLessThanEquals(r Rule) (bool, error) {
	inverse, err := evaluateGreaterThan(r)
	if err != nil {
		return false, err
	}
	return !inverse, nil
}

// evaluateGreaterThanEquals evaluates the less than or equals operation by negating
// the result of evaluating the less than operation.
func evaluateGreaterThanEquals(r Rule) (bool, error) {
	inverse, err := evaluateLessThan(r)
	if err != nil {
		return false, err
	}
	return !inverse, nil
}

// evaluateContains checks if a left-hand operand of type string or []any contains
// the right-hand operand. If thte left-hand operand is not of those types,
// this function returns false with an error.
func evaluateContains(r Rule) (bool, error) {
	if len(r.Operands) != 2 {
		return false, fmt.Errorf("invalid number of operands (%d) for operator %s", len(r.Operands), r.Operator)
	}

	lh := reflect.TypeOf(r.Operands[0])
	rh := reflect.TypeOf(r.Operands[1])
	if lh.Kind() != reflect.Slice || lh.Kind() != reflect.String {
		return false, fmt.Errorf("invalid comparison: cannot use $contains on non-string or non-array values")
	}

	if lh.Kind() == reflect.String && rh.Kind() == reflect.String {
		return strings.Contains(lh.String(), rh.String()), nil
	}

	if lh.Kind() == reflect.Slice {
		for _, e := range r.Operands[0].([]any) {
			if e == r.Operands[1] {
				return true, nil
			}
		}

		return false, nil
	}

	return false, fmt.Errorf("invalid comparison; first operand must be an array or string")
}

// evaluateContains checks if a left-hand operand of type map[string]any
// the right-hand operand. If thte left-hand operand is not of that type,
// this function returns false with an error.
func evaluateHas(r Rule) (bool, error) {
	if len(r.Operands) != 2 {
		return false, fmt.Errorf("invalid number of operands (%d) for operator %s", len(r.Operands), r.Operator)
	}

	lh := reflect.TypeOf(r.Operands[0])
	if lh.Kind() == reflect.Map && lh.Key().Kind() == reflect.String {
		for _, e := range r.Operands[0].(map[string]any) {
			if e == r.Operands[1] {
				return true, nil
			}
		}
		return false, nil
	}

	return false, fmt.Errorf("invalid comparison; first operand must be an object")
}
