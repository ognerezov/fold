package security

var (
	secretKey     = []byte("secret-key")
	guestPassword = "1234"
)

const (
	Query = "query"
	Body  = "body"
)

type CredentialsUsage struct {
	Required         bool `json:"required"`
	Save             bool `json:"save"`
	ExchangeRequired bool `json:"exchangeRequired"`
}

type Credentials struct {
	Secrets map[string]string `json:"secrets"`
	Usage   CredentialsUsage  `json:"usage"`
	In      string            `json:"in"`
}

func (c *Credentials) SecretInputRequired() bool {
	if !c.Usage.Required {
		return false
	}
	for _, v := range c.Secrets {
		if v == "" {
			return true
		}
	}
	return false
}
