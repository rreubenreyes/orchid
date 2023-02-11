package rules

import (
	"errors"
	"reflect"
	"time"
)

type Operator string
type Scalar string

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

	Number   Scalar = "number"
	String   Scalar = "string"
	Boolean  Scalar = "bool"
	DateTime Scalar = "date_time"
)

type Rule struct {
	Scalar    Scalar
	Operator  Operator
	LHOperand any
	RHOperand any
	Children  []Rule
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
	}
}

func evaluateEquals(r Rule) (bool, error) {
	return r.LHOperand == r.RHOperand, nil
}

func evaluateNot(r Rule) (bool, error) {
	if len(r.Children) == 1 {
		result, err := Evaluate(r)
		if err != nil {
			return false, err
		}

		return result, nil
	}

	return false, errors.New("invalid comparison")
}

func evaluateLessThan(r Rule) (bool, error) {
	lh := reflect.ValueOf(r.LHOperand)
	rh := reflect.ValueOf(r.RHOperand)
	if lh.Kind() == reflect.Int && rh.Kind() == reflect.Int {
		return lh.Int() < rh.Int(), nil
	}
	if lh.Kind() == reflect.Float64 && rh.Kind() == reflect.Float64 {
		return lh.Float() < rh.Float(), nil
	}
	if lh.Kind() == reflect.String && rh.Kind() == reflect.String {
		if r.Scalar == "date_time" {
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
	return false, errors.New("invalid comparison")
}

func evaluateGreaterThan(r Rule) (bool, error) {
	lh := reflect.ValueOf(r.LHOperand)
	rh := reflect.ValueOf(r.RHOperand)
	if lh.Kind() == reflect.Int && rh.Kind() == reflect.Int {
		return lh.Int() > rh.Int(), nil
	}
	if lh.Kind() == reflect.Float64 && rh.Kind() == reflect.Float64 {
		return lh.Float() > rh.Float(), nil
	}
	if lh.Kind() == reflect.String && rh.Kind() == reflect.String {
		if r.Scalar == "date_time" {
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
	return false, errors.New("invalid comparison")
}
