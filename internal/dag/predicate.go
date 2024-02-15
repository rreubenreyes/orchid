package dag

// "github.com/rreubenreyes/orchid/internal/state"

type Predicate struct {
	Variable       string       `json:"variable"`
	BoolEq         *bool        `json:"bool_eq"`
	StrEq          *string      `json:"str_eq"`
	NumEq          *float64     `json:"num_eq"`
	NumLT          *float64     `json:"num_lt"`
	NumLTE         *float64     `json:"num_lte"`
	NumGT          *float64     `json:"num_gt"`
	NumGTE         *float64     `json:"num_gte"`
	ContainsSubstr *string      `json:"contains_substr"`
	IsSubstrOf     *string      `json:"is_substr_of"`
	Contains       *any         `json:"contains"`
	IsElementOf    *[]any       `json:"is_element_of"`
	And            *[]Predicate `json:"and"`
	Or             *[]Predicate `json:"or"`
	Not            *Predicate   `json:"not"`
}
