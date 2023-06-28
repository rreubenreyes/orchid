package dag

import "testing"

func TestDAG_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		d       *DAG
		args    args
		wantErr bool
	}{
		{
			name:    "data is empty",
			d:       &DAG{},
			args:    args{data: []byte("")},
			wantErr: true,
		},
		{
			name:    "DAG is empty",
			d:       &DAG{},
			args:    args{data: []byte("{}")},
			wantErr: true,
		},
		{
			name: "simplest possible graph",
			d:    &DAG{},
			args: args{data: []byte(`{
				"start": {
					"rules": [
						{"end": true }
					]
				}
			}`)},
			wantErr: false,
		},
		{
			name: "simplest possible graph with at least one traversal",
			d:    &DAG{},
			args: args{data: []byte(`{
				"start": {
					"rules": [
						{ "next": "A" }
					]
				},
				"A": {
					"rules": [
						{"end": true }
					]
				}
			}`)},
			wantErr: false,
		},
		{
			name:    `some node contains less than one rule`,
			d:       &DAG{},
			args:    args{data: []byte(`{"start": {"rules": []}}`)},
			wantErr: true,
		},
		{
			name: `DAG does not contain a "start" node`,
			d:    &DAG{},
			args: args{data: []byte(`{
				"not_start": {
					"rules": [
						{ "next": "end" }
					]
				}
			}`)},
			wantErr: true,
		},
		{
			name: `DAG specifies "next" and "end" simultaneously`,
			d:    &DAG{},
			args: args{data: []byte(`{
				"start": {
					"rules": [
						{ "next": "end", "end": true }
					]
				}
			}`)},
			wantErr: true,
		},
		{
			name: `DAG specifies "next" and "wait" simultaneously`,
			d:    &DAG{},
			args: args{data: []byte(`{
				"start": {
					"rules": [
						{ "next": "end", "end": true }
					]
				}
			}`)},
			wantErr: true,
		},
		{
			name: `DAG specifies "wait" and "end" simultaneously`,
			d:    &DAG{},
			args: args{data: []byte(`{
				"start": {
					"rules": [
						{ "wait": true, "end": true }
					]
				}
			}`)},
			wantErr: true,
		},
		{
			name: "graph is cyclic",
			d:    &DAG{},
			args: args{data: []byte(`{
				"start": {
					"rules": [
						{ "next": "A" }
					]
				},
				"A": {
					"rules": [
						{ "next": "start" }
					]
				}
			}`)},
			wantErr: true,
		},
		{
			name: "DAG contains isolated nodes",
			d:    &DAG{},
			args: args{data: []byte(`{
				"start": {
					"rules": [
						{ "end": true }
					]
				},
				"A": {
					"rules": [
						{ "next": "start" }
					]
				}
			}`)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("%s: DAG.UnmarshalJSON() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}