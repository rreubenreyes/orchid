package main

import (
	"reflect"
	"testing"

	"github.com/linkedin/goavro/v2"
)

func TestState_ValueAtPath(t *testing.T) {
	type fields struct {
		codec *goavro.Codec
		value map[string]any
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
		// TODO: Add test cases.
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &State{
				codec: tt.fields.codec,
				value: tt.fields.value,
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