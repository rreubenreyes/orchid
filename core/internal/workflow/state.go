package workflow

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
	type alias Field
	a := struct {
		*alias                 // embed alias type; will partially unmarshal
		Value  json.RawMessage `json:"value,omitempty"`
	}{
		alias: (*alias)(f), // this alias shares the address of f
	}
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	// handle nullable case where provided Value is null
	if a.Nullable && string(a.Value) == "null" {
		f.Value = nil

		return nil
	}

	// handle non-nullable case where provided Value is null
	if !a.Nullable && string(a.Value) == "null" {
		return fmt.Errorf("field is non-nullable, got null")
	}

	// in all other cases, check type against a.Type accordingly
	switch a.Type {
	case "number":
		var n float64
		if err := json.Unmarshal(a.Value, &n); err != nil {
			return fmt.Errorf("expected number, got %s", string(a.Value))
		}
		f.Value = n

		return nil
	case "string":
		var s string
		if err := json.Unmarshal(a.Value, &s); err != nil {
			return fmt.Errorf("expected string, got %s", string(a.Value))
		}
		f.Value = s

		return nil
	case "bool":
		var b bool
		if err := json.Unmarshal(a.Value, &b); err != nil {
			return fmt.Errorf("expected bool, got %s", string(a.Value))
		}
		f.Value = b

		return nil
	case "list":
		var l []Field
		if err := json.Unmarshal(a.Value, &l); err != nil {
			return err
		}
		f.Value = l

		return nil
	case "map":
		var m map[string]Field
		if err := json.Unmarshal(a.Value, &m); err != nil {
			return err
		}
		f.Value = m

		return nil
	default:
		return fmt.Errorf("invalid type: %s", f.Type)
	}
}
