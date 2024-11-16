package openapi

import "regexp"

var (
	PathRegex = regexp.MustCompile("\\{([a-zA-Z0-9_]+)}")
)

func (a *ApiDescription) Describe(path string, obj Path) *Path {
	a.Paths[path] = obj
	return &obj
}

func (a *ApiDescription) Path(path string) *Path {
	p := Path{
		"description": path,
	}
	a.Paths[path] = p
	return &p
}

func (p *Path) Method(m string, method Method) *Path {
	(*p)[m] = method
	return p
}

func (p *Path) Get(method Method) *Path {
	return p.Method(Get, method)
}

func (p *Path) Post(method Method) *Path {
	return p.Method(Post, method)
}

func (p *Path) Patch(method Method) *Path {
	return p.Method(Patch, method)
}

func (p *Path) Delete(method Method) *Path {
	return p.Method(Delete, method)
}

func (p *Path) Put(method Method) *Path {
	return p.Method(Put, method)
}

func (p *Path) Path() string {
	val, ok := (*p)["description"]
	if !ok {
		return ""
	}
	return val.(string)
}

func (p *Path) PathParams() []Parameter {
	res := make([]Parameter, 0)
	params := PathRegex.FindAllStringSubmatch(p.Path(), -1)
	for _, param := range params {
		res = append(res, Parameter{
			Name:     param[1],
			Required: true,
			In:       "path",
			Style:    "simple",
			Schema:   AString,
		})
	}

	return res
}

func (p *Path) GetJson(summary string, schema Schema) *Path {
	return p.Get(*(&Method{
		Summary: summary,
		Responses: map[string]Response{
			"200": {
				Description: "Json data",
				Content: map[string]Content{
					ApplicationJson: {
						Schema: schema,
					},
				},
			},
		},
	}).WithParams(p.PathParams()))
}

func (p *Path) ReceiveJson(m string, summary string, inSchema, outSchema Schema) *Path {
	return p.Method(m,
		*(&Method{
			Summary: summary,
			RequestBody: &RequestBody{
				Required: true,
				Content: map[string]Content{
					ApplicationJson: {Schema: inSchema},
				},
			},
			Responses: map[string]Response{
				"200": {
					Description: "Json data",
					Content: map[string]Content{
						ApplicationJson: {
							Schema: outSchema,
						},
					},
				},
			},
		}).WithParams(p.PathParams()))
}

func (p *Path) PostJson(summary string, inSchema, outSchema Schema) *Path {
	return p.ReceiveJson(Post, summary, inSchema, outSchema)
}

func (p *Path) PatchJson(summary string, inSchema, outSchema Schema) *Path {
	return p.ReceiveJson(Patch, summary, inSchema, outSchema)
}

func (p *Path) PutJson(summary string, inSchema, outSchema Schema) *Path {
	return p.ReceiveJson(Put, summary, inSchema, outSchema)
}
