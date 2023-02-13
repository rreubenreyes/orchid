package state

import (
	"testing"
)

func TestField_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"invalid type",
			args{[]byte(`{"type":"invalid","value":1,"nullable":false}`)}, true},

		{"number (non-nullable, not null)",
			args{[]byte(`{"type":"number","value":1,"nullable":false}`)}, false},
		{"number (non-nullable, is null)",
			args{[]byte(`{"type":"number","value":null,"nullable":false}`)}, true},
		{"number (nullable, not null)",
			args{[]byte(`{"type":"number","value":1,"nullable":true}`)}, false},
		{"number (nullable, is null)",
			args{[]byte(`{"type":"number","value":null,"nullable":true}`)}, false},
		{"number (invalid type)",
			args{[]byte(`{"type":"number","value":"invalid","nullable":false}`)}, true},
		{"number (empty)",
			args{[]byte(`{"type":"number","nullable":false}`)}, true},

		{"string (non-nullable, not null)",
			args{[]byte(`{"type":"string","value":"valid","nullable":false}`)}, false},
		{"string (non-nullable, is null)",
			args{[]byte(`{"type":"string","value":null,"nullable":false}`)}, true},
		{"string (nullable, not null)",
			args{[]byte(`{"type":"string","value":"valid","nullable":true}`)}, false},
		{"string (nullable, is null)",
			args{[]byte(`{"type":"string","value":null,"nullable":true}`)}, false},
		{"string (invalid type)",
			args{[]byte(`{"type":"string","value":1,"nullable":false}`)}, true},
		{"string (empty)",
			args{[]byte(`{"type":"string","nullable":false}`)}, true},

		{"bool (non-nullable, not null)",
			args{[]byte(`{"type":"bool","value":true,"nullable":false}`)}, false},
		{"bool (non-nullable, is null)",
			args{[]byte(`{"type":"bool","value":null,"nullable":false}`)}, true},
		{"bool (nullable, not null)",
			args{[]byte(`{"type":"bool","value":true,"nullable":true}`)}, false},
		{"bool (nullable, is null)",
			args{[]byte(`{"type":"bool","value":null,"nullable":true}`)}, false},
		{"bool (invalid type)",
			args{[]byte(`{"type":"bool","value":1,"nullable":false}`)}, true},
		{"bool (empty)",
			args{[]byte(`{"type":"bool","nullable":false}`)}, true},

		{"list (non-nullable, not null)",
			args{[]byte(`{"type":"list","value":[],"nullable":false}`)}, false},
		{"list (non-nullable, is null)",
			args{[]byte(`{"type":"list","value":null,"nullable":false}`)}, true},
		{"list (nullable, not null)",
			args{[]byte(`{"type":"list","value":[],"nullable":true}`)}, false},
		{"list (nullable, is null)",
			args{[]byte(`{"type":"list","value":null,"nullable":true}`)}, false},
		{"list (invalid type)",
			args{[]byte(`{"type":"list","value":1,"nullable":false}`)}, true},
		{"list (empty)",
			args{[]byte(`{"type":"bool","nullable":false}`)}, true},
		{"list (element type is Field)", args{[]byte(`
			{"type":"list","value":[{"type":"bool","value":true,"nullable":true}],"nullable":false}
			`)}, false},
		{"list (element type is not Field)", args{[]byte(`
			{"type":"list","value":[1],"nullable":false}
			`)}, true},

		{"map (non-nullable, not null)",
			args{[]byte(`{"type":"map","value":{},"nullable":false}`)}, false},
		{"map (non-nullable, is null)",
			args{[]byte(`{"type":"map","value":null,"nullable":false}`)}, true},
		{"map (nullable, not null)",
			args{[]byte(`{"type":"map","value":{},"nullable":true}`)}, false},
		{"map (nullable, is null)",
			args{[]byte(`{"type":"map","value":null,"nullable":true}`)}, false},
		{"map (invalid type)",
			args{[]byte(`{"type":"map","value":1,"nullable":false}`)}, true},
		{"map (empty)",
			args{[]byte(`{"type":"bool","nullable":false}`)}, true},
		{"map (underlying value type is Field)", args{[]byte(`
			{"type":"map","value":{"nested":{"type":"bool","value":true,"nullable":true}},"nullable":false}
			`)}, false},
		{"map (underlying value type is not Field)", args{[]byte(`
			{"type":"map","value":{"nested":1},"nullable":false}
			`)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Field{}
			if err := f.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Field.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
