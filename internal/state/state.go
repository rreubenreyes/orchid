package state

import (
	"encoding/json"
	"fmt"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/jmoiron/jsonq"
	"github.com/xeipuuv/gojsonschema"
)

// Client functionality
type State struct {
	schema *gojsonschema.Schema
	raw    string
}

func NewState(s string, initRaw []byte) (*State, error) {
	// validate the schema
	loader := gojsonschema.NewStringLoader(s)
	schema, err := gojsonschema.NewSchema(loader)
	if err != nil {
		return nil, err
	}

	// ensure the root is an object
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(initRaw), &data); err != nil {
		return nil, err
	}

	// ensure the schema is strict
	patch, err := jsonpatch.DecodePatch([]byte(`{
		"op": "add",
		"path", "/additionalProperties",
		"value": false
	}`))
	if err != nil {
		panic(err)
	}

	init, err := patch.Apply([]byte(initRaw))
	if err != nil {
		panic(err)
	}

	// validate the initial state object
	initLoader := gojsonschema.NewStringLoader(string(init))
	result, err := schema.Validate(initLoader)
	if err != nil {
		fmt.Printf("error validating schema")
		return nil, err
	}

	if result.Valid() {
		fmt.Printf("initial state is valid")
	} else {
		fmt.Printf("initial state is not valid. see errors:\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}

	state := &State{
		schema: schema,
		raw:    string(init),
	}

	return state, nil
}

// TEST: State values can be retrieved.
func (s State) Get(path ...string) (any, error) {
	data := map[string]any{}
	dec := json.NewDecoder(strings.NewReader(s.raw))
	dec.Decode(&data)
	jq := jsonq.NewQuery(data)

	val, err := jq.Interface(path...)
	if err != nil {
		return nil, err
	}

	return val, nil
}

// TEST: State updates are expressed as jsonpatch-compatible operations.
// TEST: State updates are validated after performing a jsonpatch.
func (s *State) Set(val any, path ...string) error {
	// all calls to Set specify a replace jsonpatch op
	patchSpec := make(map[string]any)
	patchSpec["op"] = "replace"
	patchSpec["path"] = "/" + strings.Join(path, "/")
	patchSpec["value"] = val

	patchData, err := json.Marshal(patchSpec)
	if err != nil {
		return err
	}

	patch, err := jsonpatch.DecodePatch(patchData)
	if err != nil {
		return err
	}

	result, err := patch.Apply([]byte(s.raw))
	if err != nil {
		return err
	}

	// make sure updated state is valid
	resultLoader := gojsonschema.NewStringLoader(string(result))
	_, err = s.schema.Validate(resultLoader)
	if err != nil {
		return err
	}

	s.raw = string(result)

	return nil
}
