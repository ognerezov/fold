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
	AString = Schema{
		Type: "string",
	}
)
