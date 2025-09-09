package confluence

import (
	"fmt"
	"net/http"
)

const (
	defaultUa = "Go-http-client/1.1"
)

type uaRT struct {
	http.RoundTripper
	Name string
	ua   string
}

func (ur *uaRT) getUa() string {
	if ur.ua != "" {
		return ur.ua
	}
	if ur.Name == "" {
		ur.ua = defaultUa
		return ur.ua
	}
	ur.ua = fmt.Sprintf("%s %s", ur.Name, defaultUa)
	return ur.ua
}
func (ur *uaRT) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", ur.getUa())
	return ur.RoundTripper.RoundTrip(r)
}
