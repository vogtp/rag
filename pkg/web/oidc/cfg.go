package oidc

// Config configures the oidc middleware
type Config struct {
	ClientID     string
	ClientSecret string
	Issuer       string
	RedirectURI  string
	LoginPath    string
	Scopes       []string
	ResponseMode string
	KeyPath      string
}

func (c Config) Valid() bool {
	return len(c.ClientID) > 0 &&len(c.ClientSecret) > 0 &&len(c.Issuer) > 0 &&len(c.RedirectURI) > 0 
}