package rules

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

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

type Operator string
type Format string

type Int int
type Float float64
type String string
type Bool bool
type Complex []Operable
type Map map[string]Operable

type Operable interface {
	isOperable()
}

type Rule struct {
	Format   Format
	Operator Operator
	Operands []Operable
}

func (Int) isOperable()     {}
func (Float) isOperable()   {}
func (String) isOperable()  {}
func (Bool) isOperable()    {}
func (Complex) isOperable() {}
func (Map) isOperable()     {}
func (Rule) isOperable()    {}

func Evaluate(r Rule) (bool, error) {
	switch r.Operator {
	case Equals:
		return eq(r)
	case LessThan:
		return lt(r)
	case GreaterThan:
		return gt(r)
	case LessThanEquals:
		return lte(r)
	case GreaterThanEquals:
		return gte(r)
	case Contains:
		return contains(r)
	case And:
		return and(r)
	case Or:
		return or(r)
	case Xor:
		return xor(r)
	case Not:
		return not(r)
	case Has:
		return has(r)
	default:
		return false, fmt.Errorf("unknown operator \"%s\"", r.Operator)
	}
}

// eq evaluates the equality operation. This function expects
// a binary expression. If this is not the case, this function returns an error.
func eq(r Rule) (bool, error) {
	if len(r.Operands) != 2 {
		return false, fmt.Errorf("invalid number of operands (%d) for operator %s", len(r.Operands), r.Operator)
	}
	return r.Operands[0] == r.Operands[1], nil
}

// lt evaluates the "less than" operation. This function expects
// a binary expression. If this is not the case, this function returns an error.
func lt(r Rule) (bool, error) {
	if len(r.Operands) != 2 {
		return false, fmt.Errorf("invalid number of operands (%d) for operator %s", len(r.Operands), r.Operator)
	}

	a := r.Operands[0]
	b := r.Operands[1]
	switch a.(type) {
	case Int:
		return int(a.(Int)) < int(b.(Int)), nil
	case Float:
		return float64(a.(Float)) < float64(b.(Float)), nil
	case String:
		if r.Format == "date_time" {
			t1, err1 := time.Parse(time.RFC3339, string(a.(String)))
			t2, err2 := time.Parse(time.RFC3339, string(b.(String)))
			if err1 == nil && err2 == nil {
				return t1.Before(t2), nil
			}
		}
		return false, fmt.Errorf("invalid comparison: cannot compare strings which are not date_time")
	default:
		return false, fmt.Errorf("invalid comparison: unsupported type: %T", a)
	}
}

// lt evaluates the "greater than" operation. This function expects
// a binary expression. If this is not the case, this function returns an error.
func gt(r Rule) (bool, error) {
	if len(r.Operands) != 2 {
		return false, fmt.Errorf("invalid number of operands (%d) for operator %s", len(r.Operands), r.Operator)
	}

	a := r.Operands[0]
	b := r.Operands[1]
	switch a.(type) {
	case Int:
		return int(a.(Int)) > int(b.(Int)), nil
	case Float:
		return float64(a.(Float)) > float64(b.(Float)), nil
	case String:
		if r.Format == "date_time" {
			if reflect.ValueOf(a).Kind() == reflect.String {
				t1, err1 := time.Parse(time.RFC3339, string(a.(String)))
				t2, err2 := time.Parse(time.RFC3339, string(b.(String)))
				if err1 == nil && err2 == nil {
					return t1.After(t2), nil
				}
			}
		}
		return false, fmt.Errorf("invalid comparison: cannot compare strings which are not date_time")
	default:
		return false, fmt.Errorf("invalid comparison: unsupported type: %T", a)
	}
}

// evaluateLessThanEquals evaluates the less than or equals operation by negating
// the result of evaluating the greater than operation.
func lte(r Rule) (bool, error) {
	inverse, err := gt(r)
	if err != nil {
		return false, err
	}
	return !inverse, nil
}

// evaluateGreaterThanEquals evaluates the less than or equals operation by negating
// the result of evaluating the less than operation.
func gte(r Rule) (bool, error) {
	inverse, err := lt(r)
	if err != nil {
		return false, err
	}
	return !inverse, nil
}

// and accepts a Rule and checks if all of its operands are also of type Rule.
// If so, then this function evaluates logical AND over all operands. If not,
// this function returns an error.
func and(r Rule) (bool, error) {
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

// or accepts a Rule and checks if all of its operands are also of type Rule.
// If so, then this function evaluates logical OR over all operands. If not,
// this function returns an error.
func or(r Rule) (bool, error) {
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

// xor accepts a Rule and checks if all of its operands are also of type Rule.
// If so, then this function evaluates logical XOR over all operands. If not,
// this function returns an error.
func xor(r Rule) (bool, error) {
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

// not accepts a Rule and checks if it has a single Operand of type Rule.
// If so, then this function negates the Operand. If not, this function returns an error.
func not(r Rule) (bool, error) {
	if len(r.Operands) != 1 {
		return false, fmt.Errorf("invalid number of operands (%d) for operator %s", len(r.Operands), r.Operator)
	}

	result, err := Evaluate(r)
	if err != nil {
		return false, err
	}

	return result, nil
}

// contains checks if a left-hand operand of type string or []any contains
// the right-hand operand. If thte left-hand operand is not of those types,
// this function returns false with an error.
func contains(r Rule) (bool, error) {
	if len(r.Operands) != 2 {
		return false, fmt.Errorf("invalid number of operands (%d) for operator %s", len(r.Operands), r.Operator)
	}

	a := reflect.TypeOf(r.Operands[0])
	b := reflect.TypeOf(r.Operands[1])
	if a.Kind() != reflect.Slice || a.Kind() != reflect.String {
		return false, fmt.Errorf("invalid comparison: cannot use $contains on non-string or non-array values")
	}

	if a.Kind() == reflect.String && b.Kind() == reflect.String {
		return strings.Contains(a.String(), b.String()), nil
	}

	if a.Kind() == reflect.Slice {
		for _, e := range r.Operands[0].(Complex) {
			if e == r.Operands[1] {
				return true, nil
			}
		}

		return false, nil
	}

	return false, fmt.Errorf("invalid comparison; first operand must be an array or string")
}

// has checks if a left-hand operand of type map[string]any
// the right-hand operand. If thte left-hand operand is not of that type,
// this function returns false with an error.
func has(r Rule) (bool, error) {
	if len(r.Operands) != 2 {
		return false, fmt.Errorf("invalid number of operands (%d) for operator %s", len(r.Operands), r.Operator)
	}

	a := reflect.TypeOf(r.Operands[0])
	if a.Kind() == reflect.Map && a.Key().Kind() == reflect.String {
		for _, e := range r.Operands[0].(Map) {
			if e == r.Operands[1] {
				return true, nil
			}
		}
		return false, nil
	}

	return false, fmt.Errorf("invalid comparison; first operand must be an object")
}
