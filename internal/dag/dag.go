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

func (d *DAG) UnmarshalJSON(data []byte) error {
	type alias DAG
	dag := alias{}
	if err := json.Unmarshal([]byte(data), &dag); err != nil {
		return errors.New("invalid DAG definition")
	}

	dd := DAG(dag)
	if err := validate(dd); err != nil {
		return err
	}

	*d = dd

	return nil
}