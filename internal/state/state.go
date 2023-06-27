package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/linkedin/goavro/v2"
)

type State struct {
	codec *goavro.Codec
	value map[string]any
}

func parseJSONPath(path string) ([]string, error) {
	// remove leading dot (.)
	if len(path) > 0 && path[0] == '.' {
		path = path[1:]
	}

	// split the path by dot (.)
	ppath := make([]string, 0)
	ppath = append(ppath, splitPath(path)...)

	return ppath, nil
}

func splitPath(path string) []string {
	var parts []string

	current := ""
	inBracket := false
	for _, c := range path {
		switch c {
		case '[':
			inBracket = true
		case ']':
			inBracket = false
		}

		if c == '.' && !inBracket {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}

	parts = append(parts, current)

	return parts
}

func (s *State) New(data []byte) error {
	if s.codec != nil {
		fmt.Printf("state already initialized with message type")
	}

	// validate schema represents a map or record type
	var m map[string]any
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	t, ok := m["type"]
	if !ok {
		return errors.New("schema is invalid: could not determine top-level field type")
	}

	switch t := t.(type) {
	case string:
		if t != "map" && t != "record" {
			return errors.New("schema is invalid: top-level field must be a map")
		}
	default:
		return errors.New("schema is invalid: could not determine top-level field type")
	}

	// marshal into Avro codec
	codec, err := goavro.NewCodec(string(data))
	if err != nil {
		return err
	}

	s.codec = codec

	return nil
}

func (s State) CanonicalSchema() string {
	return s.codec.CanonicalSchema()
}

func (s State) Value() map[string]any {
	return s.value
}

func (s *State) ValueAtPath(path string) (any, error) {
	var jsonData []byte
	var err error

	jsonData, err = json.Marshal(s.value)
	if err != nil {
		return nil, err
	}

	var result any
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}

	ppath, err := parseJSONPath(path)
	if err != nil {
		return nil, err
	}

	for _, p := range ppath {
		switch m := result.(type) {
		case map[string]any:
			result, ok := m[p]
			if !ok {
				return nil, fmt.Errorf("path '%s' does not exist", path)
			}
			return result, nil
		default:
			return nil, fmt.Errorf("path '%s' does not exist", path)
		}
	}

	return result, nil
}

func (s *State) Update(data []byte) (map[string]any, error) {
	if s.codec == nil {
		return nil, errors.New("state not initialized")
	}

	native, _, err := s.codec.NativeFromTextual(data)
	if err != nil {
		return nil, err
	}

	switch native := native.(type) {
	case map[string]any:
		s.value = native
		return s.value, nil
	default:
		return nil, errors.New("provided value did not unmarshal into a Go map")
	}
}