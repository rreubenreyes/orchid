package dag

import (
	"encoding/json"
	"errors"
)

type DAG map[string]Node

type Node struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Next string `json:"next"`
	Wait bool   `json:"wait"`
	End  bool   `json:"end"`
}

func FromJSON(data string) (DAG, error) {
	var dag DAG
	if err := json.Unmarshal([]byte(data), &dag); err != nil {
		return nil, errors.New("invalid DAG definition")
	}
	if err := validate(dag); err != nil {
		return nil, err
	}

	return dag, nil
}