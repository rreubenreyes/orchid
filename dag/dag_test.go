package dag

import "testing"

func TestValidate(t *testing.T) {
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
			name: `DAG contains a "_wait" node`,
			args: args{data: `{
				"_wait": {
					"rules": [
						{ "next": "end" }
					]
				}
			}`},
			wantErr: true,
		},
		{
			name: `DAG contains an "_end" node`,
			args: args{data: `{
				"_end": {
					"rules": [
						{ "next": "end" }
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}