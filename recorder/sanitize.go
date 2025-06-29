package recorder

const (
	Randomize = "randomize"
	Erase     = "erase"
)

type Sanitizer struct {
	Method  string   `json:"method"`
	Combine int      `json:"combine"`
	Values  [][]any  `json:"values"`
	Parents []string `json:"parents"`
}
