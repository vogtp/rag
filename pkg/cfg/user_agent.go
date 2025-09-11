package cfg

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	defaultUa = "Go-http-client/1.1"
)

func UserAgent() string {
	ua := defaultUa
	name := viper.GetString(HTTPUserAgent)
	if len(name) < 1 {
		return ua
	}
	ua = fmt.Sprintf("%s %s", name, defaultUa)
	return ua
}
