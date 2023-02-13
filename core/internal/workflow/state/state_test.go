package state

import (
	"testing"
)

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
		{"number (non-nullable, not null)", fields{},
			args{[]byte(`{"type":"number","value":1,"nullable":false}`)}, false},

		{"number (non-nullable, is null)", fields{},
			args{[]byte(`{"type":"number","value":null,"nullable":false}`)}, true},

		{"number (non-nullable, is empty)", fields{},
			args{[]byte(`{"type":"number","nullable":false}`)}, true},

		{"number (nullable, not null)", fields{},
			args{[]byte(`{"type":"number","value":1,"nullable":true}`)}, false},

		{"number (nullable, is null)", fields{},
			args{[]byte(`{"type":"number","value":null,"nullable":true}`)}, false},

		{"number (nullable, is empty)", fields{},
			args{[]byte(`{"type":"number","nullable":true}`)}, false},
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
