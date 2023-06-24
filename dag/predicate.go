package dag

import (
	"errors"
	"fmt"
	"strings"
)

type SimplePredicate struct {
	Variable    string             `json:"variable"`
	BoolEq      *bool              `json:"bool_eq"`
	StrEq       *string            `json:"str_eq"`
	NumEq       *float64           `json:"num_eq"`
	NumLT       *float64           `json:"num_lt"`
	NumLTE      *float64           `json:"num_lte"`
	NumGT       *float64           `json:"num_gt"`
	NumGTE      *float64           `json:"num_gte"`
	StrContains *string            `json:"str_contains"`
	And         *[]SimplePredicate `json:"and"`
	Or          *[]SimplePredicate `json:"or"`
	Not         *SimplePredicate   `json:"not"`
}

type SimplePredicate2 SimplePredicate

func (p SimplePredicate) Eval(s map[string]any) (bool, error) {
	if p.BoolEq != nil {
		v, ok := s[p.Variable]
		if !ok {
			return false, fmt.Errorf("variable %s is not present in state", v)
		}
		switch v := v.(type) {
		case bool:
			return v == *p.BoolEq, nil
		}

		return false, nil
	}

	if p.StrEq != nil {
		v, ok := s[p.Variable]
		if !ok {
			return false, fmt.Errorf("variable %s is not present in state", v)
		}
		switch v := v.(type) {
		case string:
			return v == *p.StrEq, nil
		}

		return false, nil
	}

	if p.NumEq != nil {
		v, ok := s[p.Variable]
		if !ok {
			return false, fmt.Errorf("variable %s is not present in state", v)
		}
		switch v := v.(type) {
		case float64:
			return v == *p.NumEq, nil
		}

		return false, nil
	}

	if p.NumLT != nil {
		v, ok := s[p.Variable]
		if !ok {
			return false, fmt.Errorf("variable %s is not present in state", v)
		}
		switch v := v.(type) {
		case float64:
			return v < *p.NumEq, nil
		}

		return false, nil
	}

	if p.NumLTE != nil {
		v, ok := s[p.Variable]
		if !ok {
			return false, fmt.Errorf("variable %s is not present in state", v)
		}
		switch v := v.(type) {
		case float64:
			return v <= *p.NumEq, nil
		}

		return false, nil
	}

	if p.NumGT != nil {
		v, ok := s[p.Variable]
		if !ok {
			return false, fmt.Errorf("variable %s is not present in state", v)
		}
		switch v := v.(type) {
		case float64:
			return v > *p.NumEq, nil
		}

		return false, nil
	}

	if p.NumGTE != nil {
		v, ok := s[p.Variable]
		if !ok {
			return false, fmt.Errorf("variable %s is not present in state", v)
		}
		switch v := v.(type) {
		case float64:
			return v >= *p.NumEq, nil
		}

		return false, nil
	}

	if p.StrContains != nil {
		v, ok := s[p.Variable]
		if !ok {
			return false, fmt.Errorf("variable %s is not present in state", v)
		}
		switch v := v.(type) {
		case string:
			return strings.Contains(v, *p.StrContains), nil
		}

		return false, nil
	}

	return false, errors.New("invalid predicate")
}