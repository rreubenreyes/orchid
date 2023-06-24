package dag

import "testing"

func TestFromJSON(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "data is empty",
			args:    args{data: ""},
			wantErr: true,
		},
		{
			name:    "DAG is empty",
			args:    args{data: "{}"},
			wantErr: true,
		},
		{
			name: "simplest possible graph",
			args: args{data: `{
				"start": {
					"rules": [
						{"end": true }
					]
				}
			}`},
			wantErr: false,
		},
		{
			name: "simplest possible graph with at least one traversal",
			args: args{data: `{
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
			}`},
			wantErr: false,
		},
		{
			name:    `some node contains less than one rule`,
			args:    args{data: `{"start": {"rules": []}}`},
			wantErr: true,
		},
		{
			name: `DAG does not contain a "start" node`,
			args: args{data: `{
				"not_start": {
					"rules": [
						{ "next": "end" }
					]
				}
			}`},
			wantErr: true,
		},
		{
			name: `DAG specifies "next" and "end" simultaneously`,
			args: args{data: `{
				"start": {
					"rules": [
						{ "next": "end", "end": true }
					]
				}
			}`},
			wantErr: true,
		},
		{
			name: `DAG specifies "next" and "wait" simultaneously`,
			args: args{data: `{
				"start": {
					"rules": [
						{ "next": "end", "end": true }
					]
				}
			}`},
			wantErr: true,
		},
		{
			name: `DAG specifies "wait" and "end" simultaneously`,
			args: args{data: `{
				"start": {
					"rules": [
						{ "wait": true, "end": true }
					]
				}
			}`},
			wantErr: true,
		},
		{
			name: "graph is cyclic",
			args: args{data: `{
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
			}`},
			wantErr: true,
		},
		{
			name: "DAG contains isolated nodes",
			args: args{data: `{
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
			}`},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := FromJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("%s: FromJSON() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}