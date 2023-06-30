package dag

import (
	"errors"
	"reflect"
	"strings"

	"github.com/rreubenreyes/orchid/internal/state"
)

type Predicate struct {
	Variable       string       `json:"variable"`
	BoolEq         *bool        `json:"bool_eq"`
	StrEq          *string      `json:"str_eq"`
	NumEq          *float64     `json:"num_eq"`
	NumLT          *float64     `json:"num_lt"`
	NumLTE         *float64     `json:"num_lte"`
	NumGT          *float64     `json:"num_gt"`
	NumGTE         *float64     `json:"num_gte"`
	ContainsSubstr *string      `json:"contains_substr"`
	IsSubstrOf     *string      `json:"is_substr_of"`
	Contains       *any         `json:"contains"`
	IsElementOf    *[]any       `json:"is_element_of"`
	And            *[]Predicate `json:"and"`
	Or             *[]Predicate `json:"or"`
	Not            *Predicate   `json:"not"`
}

func safeIsNil(i any) bool {
	switch v := i.(type) {
	case bool, string, int, int8,
		int16, int32, int64, float32, float64:
		return false
	default:
		return reflect.ValueOf(v).IsNil()
	}
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
		case int:
			n := int(*p.NumEq)
			return v == n, nil
		case int32:
			n := int32(*p.NumEq)
			return v == n, nil
		case int64:
			n := int64(*p.NumEq)
			return v == n, nil
		case float32:
			n := float32(*p.NumEq)
			return v == n, nil
		case float64:
			n := *p.NumEq
			return v == n, nil
		}

		return false, nil
	}

	if p.NumLT != nil {
		switch v := v.(type) {
		case int:
			n := int(*p.NumLT)
			return v == n, nil
		case int32:
			n := int32(*p.NumLT)
			return v == n, nil
		case int64:
			n := int64(*p.NumLT)
			return v == n, nil
		case float32:
			n := float32(*p.NumLT)
			return v == n, nil
		case float64:
			n := *p.NumLT
			return v == n, nil
		}

		return false, nil
	}

	if p.NumLTE != nil {
		switch v := v.(type) {
		case int:
			n := int(*p.NumLTE)
			return v == n, nil
		case int32:
			n := int32(*p.NumLTE)
			return v == n, nil
		case int64:
			n := int64(*p.NumLTE)
			return v == n, nil
		case float32:
			n := float32(*p.NumLTE)
			return v == n, nil
		case float64:
			n := *p.NumLTE
			return v == n, nil
		}

		return false, nil
	}

	if p.NumGT != nil {
		switch v := v.(type) {
		case int:
			n := int(*p.NumGT)
			return v == n, nil
		case int32:
			n := int32(*p.NumGT)
			return v == n, nil
		case int64:
			n := int64(*p.NumGT)
			return v == n, nil
		case float32:
			n := float32(*p.NumGT)
			return v == n, nil
		case float64:
			n := *p.NumGT
			return v == n, nil
		}

		return false, nil
	}

	if p.NumGTE != nil {
		switch v := v.(type) {
		case int:
			n := int(*p.NumGTE)
			return v == n, nil
		case int32:
			n := int32(*p.NumGTE)
			return v == n, nil
		case int64:
			n := int64(*p.NumGTE)
			return v == n, nil
		case float32:
			n := float32(*p.NumGTE)
			return v == n, nil
		case float64:
			n := *p.NumGTE
			return v == n, nil
		}

		return false, nil
	}

	if p.IsSubstrOf != nil {
		switch v := v.(type) {
		case string:
			return strings.Contains(*p.IsSubstrOf, v), nil
		}

		return false, nil
	}

	if p.ContainsSubstr != nil {
		switch v := v.(type) {
		case string:
			return strings.Contains(v, *p.ContainsSubstr), nil
		}

		return false, nil
	}

	if p.IsElementOf != nil {
		vv := reflect.ValueOf(v)
		for _, pp := range *p.IsElementOf {
			pv := reflect.ValueOf(pp)
			if safeIsNil(v) && safeIsNil(pp) {
				return true, nil
			}
			if vv.CanInt() && pv.CanInt() {
				return vv.Int() == pv.Int(), nil
			}
			if vv.CanFloat() && pv.CanFloat() {
				return vv.Float() == pv.Float(), nil
			}
			if vv.Kind() == reflect.String && pv.Kind() == reflect.String {
				return vv.String() == pv.String(), nil
			}
			if vv.Kind() == reflect.Bool && pv.Kind() == reflect.Bool {
				return vv.Bool() == pv.Bool(), nil
			}
		}
		return false, nil
	}

	return false, errors.New("invalid predicate")
}