// package state
package main

import (
	"encoding/json"
	"fmt"
)

type Field struct {
	Type    string `json:"type,omitempty"`
	Initial any    `json:"initial,omitempty"`
	Value   any    `json:"value,omitempty"`
}

type State map[string]Field

func (f *Field) UnmarshalJSON(data []byte) error {
	var m map[string]any
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	var ok bool

	// check declared type is valid
	f.Type, ok = m["type"].(string)
	if !ok {
		return fmt.Errorf("invalid type: %s", f.Type)
	}

	switch f.Type {
	case "number":
		if f.Initial, ok = m["initial"].(float64); ok {
			return nil
		}
		f.Initial = 0
		return fmt.Errorf("expected number, got %T", m["initial"])
	case "string":
		if f.Initial, ok = m["initial"].(string); ok {
			return nil
		}
		f.Initial = ""
		return fmt.Errorf("expected string, got %T", m["initial"])
	case "bool":
		if f.Initial, ok = m["initial"].(bool); ok {
			return nil
		}
		f.Initial = false
		return fmt.Errorf("expected bool, got %T", m["initial"])
	case "list":
		if f.Initial, ok = m["initial"].([]any); ok {
			return nil
		}
		f.Initial = nil
		return fmt.Errorf("expected list, got %T", m["initial"])
	case "map":
		if f.Initial, ok = m["initial"].(map[string]any); ok {
			return nil
		}
		f.Initial = nil
		return fmt.Errorf("expected map, got %T", m["initial"])
	default:
		return fmt.Errorf("invalid type: %s", f.Type)
	}
}

func main() {
	// XXX: below is test code, remove if it's sitll here
	var s State
	jsonStr := `{
		"field0": {
			"type": "number",
			"initial": 0
		}
	}
	`

	err := json.Unmarshal([]byte(jsonStr), &s)
	if err != nil {
		return
	}
	fmt.Printf("e: %+v\n", s)
}
