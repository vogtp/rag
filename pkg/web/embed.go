package web

import (
	"embed"
	"html/template"
)

var (
	// if ng-build fails ng/intrasearch/dist/intrasearch/browser may contain no files
	//go:embed templates static ng/intrasearch/dist/intrasearch/browser
	assetData embed.FS
	templates = template.Must(template.ParseFS(assetData, "templates/*.gohtml", "templates/common/*.gohtml"))
)
