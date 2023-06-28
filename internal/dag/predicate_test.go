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
		"type": "record",
		"name": "MyRecord",
		"fields": [
			{
				"name": "foo",
				"type": "%s"
			}
		]
	}`, t))
}

func TestPredicate_Eval(t *testing.T) {
	type fields struct {
		Variable    string
		BoolEq      *bool
		StrEq       *string
		NumEq       *float64
		NumLT       *float64
		NumLTE      *float64
		NumGT       *float64
		NumGTE      *float64
		StrContains *string
		And         *[]Predicate
		Or          *[]Predicate
		Not         *Predicate
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
			name: "NumEq passing",
			fields: fields{
				Variable: ".foo",
				NumEq:    &FLOAT,
			},
			args: args{
				s: mustInitState(fooState("long"), []byte(`{"foo": 1}`)),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Predicate{
				Variable:    tt.fields.Variable,
				BoolEq:      tt.fields.BoolEq,
				StrEq:       tt.fields.StrEq,
				NumEq:       tt.fields.NumEq,
				NumLT:       tt.fields.NumLT,
				NumLTE:      tt.fields.NumLTE,
				NumGT:       tt.fields.NumGT,
				NumGTE:      tt.fields.NumGTE,
				StrContains: tt.fields.StrContains,
				And:         tt.fields.And,
				Or:          tt.fields.Or,
				Not:         tt.fields.Not,
			}
			got, err := p.Eval(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Predicate.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Predicate.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}