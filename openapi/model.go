package openapi

type Api struct {
	Openapi    string          `json:"openapi"`
	Info       Info            `json:"info"`
	Servers    []Server        `json:"servers"`
	Paths      map[string]Path `json:"paths"`
	Components Components      `json:"components"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type Server struct {
	Url         string         `json:"url"`
	Description string         `json:"description"`
	Variables   map[string]any `json:"variables"`
}

type Example struct {
	Summary       string         `json:"summary"`
	Value         map[string]any `json:"value"`
	ExternalValue string         `json:"externalValue"`
	Ref           string         `json:"$ref"`
}

type Parameter struct {
	Name     string             `json:"name"`
	In       string             `json:"in"`
	Required bool               `json:"required"`
	Style    string             `json:"style"`
	Schema   Schema             `json:"schema"`
	Examples map[string]Example `json:"examples"`
}

type Schema struct {
	Type                 string            `json:"type"`
	Properties           map[string]Schema `json:"properties"`
	Description          string            `json:"description"`
	Format               string            `json:"format"`
	OneOf                []Schema          `json:"oneOf"`
	Ref                  string            `json:"$ref"`
	Required             []string          `json:"required"`
	AdditionalProperties bool              `json:"additionalProperties"`
	Encoding             map[string]any    `json:"encoding"`
}

type RequestBody struct {
	Ref         string             `json:"$ref"`
	Description string             `json:"description"`
	Content     map[string]Content `json:"content"`
	Required    bool               `json:"required"`
}

type Content struct {
	Schema      Schema             `json:"schema"`
	Examples    map[string]Example `json:"examples"`
	Description string             `json:"description"`
}

type Path struct {
	Get     Method `json:"get"`
	Post    Method `json:"post"`
	Put     Method `json:"put"`
	Patch   Method `json:"patch"`
	Delete  Method `json:"delete"`
	Head    Method `json:"head"`
	Options Method `json:"options"`
}

type Method struct {
	Summary     string             `json:"summary"`
	Parameters  []Parameter        `json:"parameters"`
	RequestBody RequestBody        `json:"requestBody"`
	Responses   map[string]Content `json:"responses"`
}

type Components struct {
	RequestBodies map[string]RequestBody `json:"requestBodies"`
	Schemas       map[string]Schema      `json:"schemas"`
	Examples      map[string]Example     `json:"examples"`
}
