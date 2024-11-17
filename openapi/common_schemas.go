package openapi

const (
	Get    = "get"
	Put    = "put"
	Delete = "delete"
	Post   = "post"
	Patch  = "patch"

	ApplicationJson = "application/json"
)

var (
	AnArray = Schema{
		Type: "array",
	}
	IdObject = Schema{
		Type: "object",
		Properties: map[string]Schema{
			"id": {
				Type: "string",
			},
		},
		Required: []string{"id"},
	}
	CountObject = Schema{
		Type: "object",
		Properties: map[string]Schema{
			"count": {
				Type: "number",
			},
		},
		Required: []string{"count"},
	}
	AString = Schema{
		Type: "string",
	}
	AnObject = Schema{
		Type: "object",
	}
)

func (s *Schema) ToQueryParams() []Parameter {
	res := make([]Parameter, len(s.Properties))
	count := 0
	for name, schema := range s.Properties {
		res[count] = Parameter{
			Name:   name,
			Schema: schema,
			Style:  "simple",
			In:     "query",
		}
		count++
	}

	return res
}
