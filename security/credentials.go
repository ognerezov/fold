package security

import "fold/console"

var (
	secretKey     []byte
	guestPassword string
)

func SetSecretKey(key string) {
	if key == "" {
		console.RedPrintln("secret key is empty. Authentication is blocked")
	}
	secretKey = []byte(key)
}

func SetGuestPassword(password string) {
	guestPassword = password
}

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
