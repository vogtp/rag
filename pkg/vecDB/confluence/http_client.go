package confluence

import (
	"net/http"

	"github.com/graze/go-throttled"
	"github.com/vogtp/rag/pkg/cfg"
	"golang.org/x/time/rate"
)

var rateLimiter *rate.Limiter

func getRateLimitHttpClient(client *http.Client) *http.Client {
	if rateLimiter == nil {
		rateLimiter = rate.NewLimiter(rate.Limit(0.4), 1)
	}
	client.Transport = &userAgentRoundTripper{
		RoundTripper: client.Transport,
	}
	return throttled.WrapClient(client, rateLimiter)
}

type userAgentRoundTripper struct {
	http.RoundTripper
	ua string
}

func (ur *userAgentRoundTripper) getUserAgent() string {
	if len(ur.ua) > 0 {
		return ur.ua
	}
	ur.ua = cfg.UserAgent()
	return ur.ua
}

func (ur *userAgentRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", ur.getUserAgent())
	return ur.RoundTripper.RoundTrip(r)
}
