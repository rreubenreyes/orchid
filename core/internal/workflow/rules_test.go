package workflow

import "testing"

func TestEvaluate(t *testing.T) {
	type args struct {
		r Rule
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"$eq: non-binary expression",
			args{Rule{"", "$eq", []Operable{Float(1.0), Float(2.0), Float(3.0)}}},
			false, true},
		{"$eq: non-scalar expression",
			args{Rule{"", "$eq", []Operable{Complex{Float(1.0)}}}},
			false, true},
		{"$eq: different types",
			args{Rule{"", "$eq", []Operable{Float(1.0), Bool(true)}}},
			false, false},
		{"$eq: Float - valid true result",
			args{Rule{"", "$eq", []Operable{Float(1.0), Float(1.0)}}},
			true, false},
		{"$eq: Float - valid false result",
			args{Rule{"", "$eq", []Operable{Float(1.0), Float(0.0)}}},
			false, false},
		{"$eq: String - valid true result",
			args{Rule{"", "$eq", []Operable{String("1"), String("1")}}},
			true, false},
		{"$eq: String - valid false result",
			args{Rule{"", "$eq", []Operable{String("1"), String("2")}}},
			false, false},
		{"$eq: Bool - valid true result",
			args{Rule{"", "$eq", []Operable{Bool(true), Bool(true)}}},
			true, false},
		{"$eq: Bool - valid false result",
			args{Rule{"", "$eq", []Operable{Bool(true), Bool(false)}}},
			false, false},

		{"$lt: non-binary expression",
			args{Rule{"", "$lt", []Operable{Float(1.0), Float(2.0), Float(3.0)}}},
			false, true},
		{"$lt: invalid type",
			args{Rule{"", "$lt", []Operable{Complex{}}}},
			false, true},
		{"$lt: different types",
			args{Rule{"", "$lt", []Operable{Float(1.0), Bool(true)}}},
			false, false},
		{"$lt: Float - valid true result",
			args{Rule{"", "$lt", []Operable{Float(0.0), Float(1.0)}}},
			true, false},
		{"$lt: Float - valid false result",
			args{Rule{"", "$lt", []Operable{Float(1.0), Float(1.0)}}},
			false, false},
		{"$lt: String - valid true result",
			args{Rule{"date_time", "$lt", []Operable{String("1999-02-13T18:13:07.386Z"), String("2000-02-13T18:13:07.386Z")}}},
			true, false},
		{"$lt: String - valid false result",
			args{Rule{"date_time", "$lt", []Operable{String("1999-02-13T18:13:07.386Z"), String("1999-02-13T18:13:07.386Z")}}},
			false, false},
		{"$lt: String - no format",
			args{Rule{"", "$lt", []Operable{String("1"), String("1")}}},
			false, true},
		{"$lt: String - not ISO-8601",
			args{Rule{"date_time", "$lt", []Operable{String("1"), String("1")}}},
			false, true},

		{"$lte: non-binary expression",
			args{Rule{"", "$lte", []Operable{Float(1.0), Float(2.0), Float(3.0)}}},
			false, true},
		{"$lte: invalid type",
			args{Rule{"", "$lte", []Operable{Complex{}}}},
			false, true},
		{"$lte: different types",
			args{Rule{"", "$lte", []Operable{Float(1.0), Bool(true)}}},
			false, false},
		{"$lte: Float - valid true result (lt)",
			args{Rule{"", "$lte", []Operable{Float(0.0), Float(1.0)}}},
			true, false},
		{"$lte: Float - valid true result (eq)",
			args{Rule{"", "$lte", []Operable{Float(0.0), Float(0.0)}}},
			true, false},
		{"$lte: Float - valid false result",
			args{Rule{"", "$lte", []Operable{Float(1.0), Float(0.0)}}},
			false, false},
		{"$lte: String - valid true result (lt)",
			args{Rule{"date_time", "$lte", []Operable{String("1999-02-13T18:13:07.386Z"), String("2000-02-13T18:13:07.386Z")}}},
			true, false},
		{"$lte: String - valid true result (eq)",
			args{Rule{"date_time", "$lte", []Operable{String("1999-02-13T18:13:07.386Z"), String("1999-02-13T18:13:07.386Z")}}},
			true, false},
		{"$lte: String - valid false result",
			args{Rule{"date_time", "$lte", []Operable{String("2000-02-13T18:13:07.386Z"), String("1999-02-13T18:13:07.386Z")}}},
			false, false},
		{"$lte: String - no format",
			args{Rule{"", "$lte", []Operable{String("1"), String("1")}}},
			false, true},
		{"$lte: String - not ISO-8601",
			args{Rule{"date_time", "$lte", []Operable{String("1"), String("1")}}},
			false, true},

		{"$gt: non-binary expression",
			args{Rule{"", "$gt", []Operable{Float(1.0), Float(2.0), Float(3.0)}}},
			false, true},
		{"$gt: invalid type",
			args{Rule{"", "$gt", []Operable{Complex{}}}},
			false, true},
		{"$gt: different types",
			args{Rule{"", "$gt", []Operable{Float(1.0), Bool(true)}}},
			false, false},
		{"$gt: Float - valid true result",
			args{Rule{"", "$gt", []Operable{Float(1.0), Float(0.0)}}},
			true, false},
		{"$gt: Float - valid false result",
			args{Rule{"", "$gt", []Operable{Float(1.0), Float(1.0)}}},
			false, false},
		{"$gt: String - valid true result",
			args{Rule{"date_time", "$gt", []Operable{String("2000-02-13T18:13:07.386Z"), String("1999-02-13T18:13:07.386Z")}}},
			true, false},
		{"$gt: String - valid false result",
			args{Rule{"date_time", "$gt", []Operable{String("1999-02-13T18:13:07.386Z"), String("1999-02-13T18:13:07.386Z")}}},
			false, false},
		{"$gt: String - no format",
			args{Rule{"", "$gt", []Operable{String("1"), String("1")}}},
			false, true},
		{"$gt: String - not ISO-8601",
			args{Rule{"date_time", "$gt", []Operable{String("1"), String("1")}}},
			false, true},

		{"$gte: non-binary expression",
			args{Rule{"", "$gte", []Operable{Float(1.0), Float(2.0), Float(3.0)}}},
			false, true},
		{"$gte: invalid type",
			args{Rule{"", "$gte", []Operable{Complex{}}}},
			false, true},
		{"$gte: different types",
			args{Rule{"", "$gte", []Operable{Float(1.0), Bool(true)}}},
			false, false},
		{"$gte: Float - valid true result (gt)",
			args{Rule{"", "$gte", []Operable{Float(1.0), Float(0.0)}}},
			true, false},
		{"$gte: Float - valid true result (eq)",
			args{Rule{"", "$gte", []Operable{Float(0.0), Float(0.0)}}},
			true, false},
		{"$gte: Float - valid false result",
			args{Rule{"", "$gte", []Operable{Float(0.0), Float(1.0)}}},
			false, false},
		{"$gte: String - valid true result (gt)",
			args{Rule{"date_time", "$gte", []Operable{String("2000-02-13T18:13:07.386Z"), String("1999-02-13T18:13:07.386Z")}}},
			true, false},
		{"$gte: String - valid true result (eq)",
			args{Rule{"date_time", "$gte", []Operable{String("1999-02-13T18:13:07.386Z"), String("1999-02-13T18:13:07.386Z")}}},
			true, false},
		{"$gte: String - valid false result",
			args{Rule{"date_time", "$gte", []Operable{String("1999-02-13T18:13:07.386Z"), String("2000-02-13T18:13:07.386Z")}}},
			false, false},
		{"$gte: String - no format",
			args{Rule{"", "$gte", []Operable{String("1"), String("1")}}},
			false, true},
		{"$gte: String - not ISO-8601",
			args{Rule{"date_time", "$gte", []Operable{String("1"), String("1")}}},
			false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Evaluate(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}