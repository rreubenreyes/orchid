package state

import (
	"reflect"
	"testing"

	"github.com/xeipuuv/gojsonschema"
)

func TestState_ValueAtPath(t *testing.T) {
	type fields struct {
		schema *gojsonschema.Schema
		value  map[string]any
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "basic example",
			fields: fields{
				value: map[string]any{
					"foo": "bar",
				},
			},
			args:    args{path: ".foo"},
			want:    "bar",
			wantErr: false,
		},
		{
			name: "basic example with square brackets",
			fields: fields{
				value: map[string]any{
					"foo": []any{"bar"},
				},
			},
			args:    args{path: ".foo[0]"},
			want:    "bar",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &State{
				schema: tt.fields.schema,
				value:  tt.fields.value,
			}
			got, err := s.ValueAtPath(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("State.ValueAtPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("State.ValueAtPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
