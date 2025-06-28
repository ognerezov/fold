package recorder

const (
	Randomize = "randomize"
)

type Sanitizer struct {
	Method  string  `json:"method"`
	Combine int     `json:"combine"`
	Values  [][]any `json:"values"`
}
