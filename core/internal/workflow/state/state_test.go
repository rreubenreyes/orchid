package state

import "testing"

func TestField_UnmarshalJSON(t *testing.T) {
	type fields struct {
		Type     string
		Value    any
		Nullable bool
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Field{
				Type:     tt.fields.Type,
				Value:    tt.fields.Value,
				Nullable: tt.fields.Nullable,
			}
			if err := f.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Field.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
