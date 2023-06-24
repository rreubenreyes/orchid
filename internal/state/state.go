package main

import (
	"fmt"

	"github.com/linkedin/goavro/v2"
)

type State struct {
	schema  *goavro.Codec
	history []map[string]any // keep state history in memory for now, will move to db or something later
	value   map[string]any
}

func (s State) Schema() *goavro.Codec {
	return s.schema
}

func (s State) Value() map[string]any {
	return s.value
}

func (s *State) LoadSchema(data []byte) error {
	codec, err := goavro.NewCodec(string(data))
	if err != nil {
		return err
	}

	s.schema = codec

	return nil
}

func (s *State) Update(data []byte) (map[string]any, error) {
	native, _, err := s.schema.NativeFromTextual(data)
	if err != nil {
		return nil, err
	}

	switch native := native.(type) {
	case map[string]any:
		s.value = native
		s.history = append(s.history, native)

		return s.value, nil
	default:
		return nil, fmt.Errorf("data unmarshaled into unexpected type")
	}
}