package state

import (
	"encoding/json"
	"fmt"
)

type Field struct {
	Type     string `json:"type"`
	Value    any    `json:"value,omitempty"`
	Nullable bool   `json:"nullable,omitempty"`
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
	if f.Type, ok = m["type"].(string); !ok {
		return fmt.Errorf("invalid type: %s", f.Type)
	}

	// if `nullable` is missing or not boolean, assume it is false
	if f.Nullable, ok = m["nullable"].(bool); !ok {
		f.Nullable = false
	}

	switch f.Type {
	case "number":
		if f.Value, ok = m["value"].(float64); !ok {
			if f.Nullable && f.Value == nil {
				f.Value = 0
				return nil
			}

			f.Value = 0 // set value to expected type
			return fmt.Errorf("expected number, got %T", m["value"])
		}
		return nil
	case "string":
		if f.Value, ok = m["initial"].(string); !ok {
			if f.Nullable && f.Value == nil {
				f.Value = 0
				return nil
			}

			f.Value = ""
			return fmt.Errorf("expected string, got %T", m["value"])
		}
		return nil
	case "bool":
		if f.Value, ok = m["initial"].(bool); !ok {
			if f.Nullable && f.Value == nil {
				f.Value = 0
				return nil
			}

			f.Value = false
			return fmt.Errorf("expected bool, got %T", m["value"])
		}
		return nil

		// TODO: recursively unmarshal complex fields
	case "list":
		if f.Value, ok = m["initial"].([]any); !ok {
			if f.Nullable && f.Value == nil {
				f.Value = 0
				return nil
			}

			f.Value = nil
			return fmt.Errorf("expected list, got %T", m["value"])
		}
		return nil
	case "map":
		if f.Value, ok = m["initial"].(map[string]any); !ok {
			if f.Nullable && f.Value == nil {
				f.Value = 0
				return nil
			}

			f.Value = nil
			return fmt.Errorf("expected map, got %T", m["value"])
		}
		return nil
	default:
		return fmt.Errorf("invalid type: %s", f.Type)
	}
}
