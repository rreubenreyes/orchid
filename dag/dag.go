package dag

import (
	"encoding/json"
	"errors"
)

type Definition struct {
	Nodes []Node `json:"nodes"`
}

type Node struct {
	ID string `json:"id"`
}

func Validate(data string) error {
	var def Definition

	err := json.Unmarshal([]byte(data), &def)
	if err != nil {
		return errors.New("invalid DAG definition")
	}

	return nil
}