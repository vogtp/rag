package cfg

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

const (
	// CfgFile
	CfgFile = "config.file"
	// CfgSave triggers periodic config saves
	CfgSave = "config.save"

	// WebListen set the address the webserver should listen on
	WebListen = "web.listen"

	// HTTPProxy enables SOCKS5 proxies for http requests
	HTTPProxy = "http.proxy"

	// LogLevel error warn info debug
	LogLevel = "log.level"
	// LogSource should we log the source
	LogSource = "log.source"
	// LogJson log in json
	LogJson = "log.json"

	// CacheDir the subdirectory to use as cache, subdir of modeldir if it does not start with /
	CacheDir = "cache.dir"

	// ModelDefault is the default model when no model is given
	ModelDefault = "model.default"
	// ModelEmbedding is the default model used for embeddings
	ModelEmbedding = "model.embedding"

	// OllamaHosts is an URL
	ollamaHosts = "ollama.hosts"

	// ChromaContainerName the name of the chroma container
	ChromaContainerName = "chroma.container.name"
	// ChromaUrl the URL where chroma can be reached
	chromaUrl = "chroma.url"
	// ChromaContainerPort is the port chroma should be started on (0: disable)
	ChromaContainerPort = "chroma.container.port"
	// ChromaContainerImage chroma container to pull
	ChromaContainerImage = "chroma.container.image"

	// ConfluenceKey is the confluence access token
	ConfluenceKey = "confluence.key"
	// ConfluenceBaseURL is the base URL of the confluence instance
	ConfluenceBaseURL = "confluence.baseURL"
	// ConfluenceSpaces defines the spaces to scrap
	ConfluenceSpaces = "confluence.spaces"
	// ConfluenceMaxAge is the maximum age a confluence page can have to be included in
	ConfluenceMaxAge = "confluence.maxAge"
	//VecDBUpdateIntervall is the intervall the vectorDB is updated
	VecDBUpdateIntervall = "vecdb.update_intervall"
	//VecDBColName is the name prefix of the vectorDB collections
	VecDBColName = "vecdb.collection_name"

	// OIDCClientID OIDC Client ID
	OIDCClientID = "oidc.client_id"
	// OIDCClientSecret OIDC Secret
	OIDCClientSecret = "oidc.client_secret"
	// OIDCIssuer OIDC issuer (AKA auth endpoint)
	OIDCIssuer = "oidc.issuer"
	// OIDCRedirectURI OIDC redirect URI (AKA local auth callback)
	OIDCRedirectURI = "oidc.redirect_uri"

	// HTTPUserAgent sets the user agent of the http request
	HTTPUserAgent = "http.user_agent"

	// PProf enable pprof
	PProf = "pprof.file"

	// FIXME
	ApiBearerToken = "web.api_token"
)

var (
	// DefaultConfluenceMaxAge is the max age of a confluence page to be included
	DefaultConfluenceMaxAge = 7 * 356 * 24 * time.Hour
)

func init() {
	pflag.Bool(CfgSave, false, "Should the configs be written to file periodically")
	pflag.String(CfgFile, fmt.Sprintf("%s.yml", appName), "File with the config to load")
	pflag.String(LogLevel, "warn", "Set the loglevel: error warn info debug trace off")
	pflag.Bool(LogSource, false, "Log the source line")
	pflag.Bool(LogJson, false, "Log in json")

	pflag.String(HTTPProxy, "", "enables SOCKS5 proxies for http requests, eg. localhost:1928")
	pflag.String(WebListen, ":8080", "Address the webserver should listen on")

	pflag.String(ModelDefault, "gpt-oss", "The default model when no model is given")
	pflag.String(ModelEmbedding, "bge-m3", "The default model used for embeddings")
	// pflag.String(ModelEmbedding, "mxbai-embed-large", "The default model used for embeddings")
	//pflag.String(ModelEmbedding, "nomic-embed-text", "The default model used for embeddings")

	pflag.String(ConfluenceKey, "", "The confluence access token")
	pflag.String(ConfluenceBaseURL, "", "The confluence access token")
	pflag.StringSlice(ConfluenceSpaces, nil, "The confluence spaces to scrap")
	pflag.Duration(ConfluenceMaxAge, DefaultConfluenceMaxAge, "The maximum age a confluence page can have to be included in")
	pflag.Duration(VecDBUpdateIntervall, DefaultVecDBUpdateIntervall, "the intervall the vectorDB is updated")
	pflag.String(VecDBColName, "go-rag", "the name prefix of the vectorDB collections")
	pflag.String(chromaUrl, chromaDefaultURL, "the URL where chroma can be reached")
	pflag.String(ChromaContainerPort, chromaDefaultPort, "the port chroma should be started on (0: disable)")
	pflag.String(ChromaContainerName, "chroma", "the name of the chroma container")
	pflag.String(ChromaContainerImage, chromaDefaultContainerImage, "chroma container to pull")
	pflag.String(OIDCClientID, "", "OIDCClientID OIDC Client ID")
	pflag.String(OIDCClientSecret, "", "OIDC Secret")
	pflag.String(OIDCIssuer, "", "OIDC issuer (AKA auth endpoint)")
	pflag.String(OIDCRedirectURI, "", "OIDC redirect URI (AKA local auth callback)")
	pflag.String(HTTPUserAgent, "go-rag", "sets the user agent of the http request")

	pflag.String(PProf, "", "enable pprof")

	pflag.String(ApiBearerToken, "", "bearer token for webAPI")
}
