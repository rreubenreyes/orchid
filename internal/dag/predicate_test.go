package dag

import (
	"fmt"
	"testing"

	"github.com/rreubenreyes/orchid/internal/state"
)

var (
	TRUE  = true
	FALSE = false
	STR   = "bar"
	FLOAT = 1.0
	ANY   = any("bar")
)

func mustInitState(d []byte, initial []byte) state.State {
	s, err := state.New(d)
	if err != nil {
		panic(err)
	}

	_, err = s.Update(initial)
	if err != nil {
		panic(err)
	}

	return s
}

func fooState(t string) []byte {
	return []byte(fmt.Sprintf(`{
		"title": "MyRecord",
		"type": "object",
		"properties": {
			"foo": {
				"type": "%s"
			}
		}
	}`, t))
}

func TestPredicate_Eval(t *testing.T) {
	type fields struct {
		Variable       string
		BoolEq         *bool
		StrEq          *string
		NumEq          *float64
		NumLT          *float64
		NumLTE         *float64
		NumGT          *float64
		NumGTE         *float64
		ContainsSubstr *string
		IsSubstrOf     *string
		Contains       *any
		IsElementOf    *[]any
		And            *[]Predicate
		Or             *[]Predicate
		Not            *Predicate
	}
	type args struct {
		s state.State
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "BoolEq passing",
			fields: fields{
				Variable: ".foo",
				BoolEq:   &TRUE,
			},
			args: args{
				s: mustInitState(fooState("boolean"), []byte(`{"foo": true}`)),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "StrEq passing",
			fields: fields{
				Variable: ".foo",
				StrEq:    &STR,
			},
			args: args{
				s: mustInitState(fooState("string"), []byte(`{"foo": "bar"}`)),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "IsSubstrOf passing",
			fields: fields{
				Variable:   ".foo",
				IsSubstrOf: &STR,
			},
			args: args{
				s: mustInitState(fooState("string"), []byte(`{"foo": "b"}`)),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "ContainsSubstr passing",
			fields: fields{
				Variable:       ".foo",
				ContainsSubstr: &STR,
			},
			args: args{
				s: mustInitState(fooState("string"), []byte(`{"foo": "barrrr"}`)),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "NumEq (float) passing",
			fields: fields{
				Variable: ".foo",
				NumEq:    &FLOAT,
			},
			args: args{
				s: mustInitState(fooState("number"), []byte(`{"foo": 1.0}`)),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "NumEq (long) passing",
			fields: fields{
				Variable: ".foo",
				NumEq:    &FLOAT,
			},
			args: args{
				s: mustInitState(fooState("number"), []byte(`{"foo": 1}`)),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "IsElementOf passing for numeric",
			fields: fields{
				Variable:    ".foo",
				IsElementOf: &[]any{1.0, "foo", nil, true},
			},
			args: args{
				s: mustInitState(fooState("number"), []byte(`{"foo": 1}`)),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "IsElementOf passing for string",
			fields: fields{
				Variable:    ".foo",
				IsElementOf: &[]any{1.0, "foo", nil, true},
			},
			args: args{
				s: mustInitState(fooState("string"), []byte(`{"foo": "foo"}`)),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "IsElementOf passing for bool",
			fields: fields{
				Variable:    ".foo",
				IsElementOf: &[]any{1.0, "foo", nil, true},
			},
			args: args{
				s: mustInitState(fooState("boolean"), []byte(`{"foo": true}`)),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "IsElementOf passing for nil",
			fields: fields{
				Variable:    ".foo",
				IsElementOf: &[]any{1, "foo", nil, true},
			},
			args: args{
				s: mustInitState(fooState("null"), []byte(`{"foo": null}`)),
			},
			want:    true,
			wantErr: false,
		},
		// TODO: this test fails; need to implement Contains
		{
			name: "Contains",
			fields: fields{
				Variable: ".foo",
				Contains: &ANY,
			},
			args: args{
				s: mustInitState([]byte(`{
					"title": "MyRecord",
					"type": "object",
					"properties": {
						"foo": {
							"type": "array",
							"items": {
								"type": "string"
							}
						}
					}
				}`), []byte(`{"foo": ["hi"]}`)),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Predicate{
				Variable:       tt.fields.Variable,
				BoolEq:         tt.fields.BoolEq,
				StrEq:          tt.fields.StrEq,
				NumEq:          tt.fields.NumEq,
				NumLT:          tt.fields.NumLT,
				NumLTE:         tt.fields.NumLTE,
				NumGT:          tt.fields.NumGT,
				NumGTE:         tt.fields.NumGTE,
				IsSubstrOf:     tt.fields.IsSubstrOf,
				ContainsSubstr: tt.fields.ContainsSubstr,
				IsElementOf:    tt.fields.IsElementOf,
				Contains:       tt.fields.Contains,
				And:            tt.fields.And,
				Or:             tt.fields.Or,
				Not:            tt.fields.Not,
			}
			got, err := p.Eval(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: Predicate.Eval() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("%s: Predicate.Eval() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
