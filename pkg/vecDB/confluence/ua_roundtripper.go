package confluence

import (
	"net/http"

	"github.com/vogtp/rag/pkg/cfg"
)

type uaRT struct {
	http.RoundTripper
	ua string
}

func (ur *uaRT) getUa() string {
	if len(ur.ua) > 0 {
		return ur.ua
	}
	ur.ua = cfg.UserAgent()
	return ur.ua
}

func (ur *uaRT) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", ur.getUa())
	return ur.RoundTripper.RoundTrip(r)
}
