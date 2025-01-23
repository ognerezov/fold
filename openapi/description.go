package openapi

type ResponseDescription interface {
	Describe() map[string]Response
}
