package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/linkedin/goavro/v2"
)

type State struct {
	codec *goavro.Codec
	value map[string]any
}

type pathPart struct {
	Value any
}

func parseJSONPath(path string) ([]string, error) {
	// remove leading dot (.)
	if len(path) > 0 && path[0] == '.' {
		path = path[1:]
	}

	// split the path by dot (.)
	var ppath []string
	sp, err := splitPath(path)
	if err != nil {
		return nil, err
	}
	ppath = append(ppath, sp...)

	return ppath, nil
}

func splitPath(path string) ([]string, error) {
	var parts []string

	current := ""
	inBracket := false
	for _, c := range path {
		switch c {
		case '[':
			if !inBracket {
				if current != "" {
					parts = append(parts, current)
					current = ""
				}
				inBracket = true
			} else {
				current += string(c)
			}
		case ']':
			if inBracket {
				inBracket = false
				if current != "" {
					index, err := strconv.Atoi(current)
					if err == nil {
						parts = append(parts, fmt.Sprintf("%d", index))
					} else {
						return nil, err
					}
					current = ""
				}
			} else {
				current += string(c)
			}
		case '.':
			if !inBracket {
				if current != "" {
					parts = append(parts, current)
					current = ""
				}
			} else {
				current += string(c)
			}
		default:
			current += string(c)
		}
	}

	if current != "" {
		parts = append(parts, current)
	}

	return parts, nil
}

func New(data []byte) (State, error) {
	// validate schema represents a map or record type
	var s State
	var m map[string]any
	err := json.Unmarshal(data, &m)
	if err != nil {
		return s, err
	}

	t, ok := m["type"]
	if !ok {
		return s, errors.New("schema is invalid: could not determine top-level field type")
	}

	switch t := t.(type) {
	case string:
		if t != "record" {
			return s, errors.New("schema is invalid: top-level field must be a map")
		}
	default:
		return s, errors.New("schema is invalid: could not determine top-level field type")
	}

	// marshal into Avro codec
	codec, err := goavro.NewCodec(string(data))
	if err != nil {
		return s, err
	}

	err = s.SetCodec(codec)
	if err != nil {
		return s, err
	}

	return s, nil
}

func (s *State) SetCodec(c *goavro.Codec) error {
	if s.codec != nil {
		return errors.New("codec already initialized")
	}

	s.codec = c

	return nil
}

func (s State) CanonicalSchema() string {
	return s.codec.CanonicalSchema()
}

func (s State) Value() map[string]any {
	return s.value
}

func (s *State) ValueAtPath(path string) (any, error) {
	ppath, err := parseJSONPath(path)
	if err != nil {
		return nil, err
	}

	result := any(s.value)
	var ok bool
	for _, p := range ppath {
		switch s := result.(type) {
		case map[string]any:
			result, ok = s[p]
			if !ok {
				return nil, fmt.Errorf("path '%s' does not exist", path)
			}
		case []any:
			i, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}

			if len(s) < i+1 {
				return nil, fmt.Errorf("path '%s' does not exist", path)
			}
			result = s[i]
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