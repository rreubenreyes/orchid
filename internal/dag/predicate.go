package dag

import (
	"errors"
	"strings"

	"github.com/rreubenreyes/orchid/internal/state"
)

type Predicate struct {
	Variable    string       `json:"variable"`
	BoolEq      *bool        `json:"bool_eq"`
	StrEq       *string      `json:"str_eq"`
	NumEq       *float64     `json:"num_eq"`
	NumLT       *float64     `json:"num_lt"`
	NumLTE      *float64     `json:"num_lte"`
	NumGT       *float64     `json:"num_gt"`
	NumGTE      *float64     `json:"num_gte"`
	StrContains *string      `json:"str_contains"`
	And         *[]Predicate `json:"and"`
	Or          *[]Predicate `json:"or"`
	Not         *Predicate   `json:"not"`
}

func (p Predicate) Eval(s state.State) (bool, error) {
	v, err := s.ValueAtPath(p.Variable)
	if err != nil {
		return false, err
	}

	if p.BoolEq != nil {
		switch v := v.(type) {
		case bool:
			return v == *p.BoolEq, nil
		}

		return false, nil
	}

	if p.StrEq != nil {
		switch v := v.(type) {
		case string:
			return v == *p.StrEq, nil
		}

		return false, nil
	}

	if p.NumEq != nil {
		switch v := v.(type) {
		case float64:
			return v == *p.NumEq, nil
		}

		return false, nil
	}

	if p.NumLT != nil {
		switch v := v.(type) {
		case float64:
			return v < *p.NumEq, nil
		}

		return false, nil
	}

	if p.NumLTE != nil {
		switch v := v.(type) {
		case float64:
			return v <= *p.NumEq, nil
		}

		return false, nil
	}

	if p.NumGT != nil {
		switch v := v.(type) {
		case float64:
			return v > *p.NumEq, nil
		}

		return false, nil
	}

	if p.NumGTE != nil {
		switch v := v.(type) {
		case float64:
			return v >= *p.NumEq, nil
		}

		return false, nil
	}

	if p.StrContains != nil {
		switch v := v.(type) {
		case string:
			return strings.Contains(v, *p.StrContains), nil
		}

		return false, nil
	}

	return false, errors.New("invalid predicate")
}