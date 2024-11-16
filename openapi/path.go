package openapi

func (a *ApiDescription) Describe(path string, obj Path) *Path {
	a.Paths[path] = obj
	return &obj
}

func (a *ApiDescription) Path(path string) *Path {
	p := Path{}
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

func (p *Path) GetJson(summary string, schema Schema) *Path {
	return p.Get(Method{
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
	})
}

func (p *Path) ReceiveJson(m string, summary string, inSchema, outSchema Schema) *Path {
	return p.Method(m,
		Method{
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
		})
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
