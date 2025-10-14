package web

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/vogtp/go-angular"
)

func (srv *Server) routes() error {

	fsys, err := fs.Sub(assetData, "ng/intrasearch/dist/intrasearch/browser")
	if err != nil {
		panic(err)
	}
	ngFS := angular.FileSystem(fsys)
	// srv.oidcMux.Handle("/", http.FileServer(ngFS))
	srv.mux.Handle("/", http.FileServer(ngFS))
	srv.mux.Handle("/static/", http.StripPrefix(srv.baseURL, http.FileServer(http.FS(assetData))))

	srv.openAiAPI("/api/")
	srv.oidcMux.HandleFunc("/vecdb/", srv.vecDBlist)
	srv.oidcMux.HandleFunc("/vecdb/{collection}", srv.vecDBsearch)
	srv.oidcMux.HandleFunc("/summary/{uuid}", srv.handleSummary)
	srv.oidcMux.HandleFunc("/user/", srv.handleUserEnt)
	srv.oidcMux.HandleFunc("/user_new/", srv.handleUser)

	// srv.oidcMux.Handle("/graphql/", usercfg.HttpHandler())
	// srv.oidcMux.Handle("/graphiql/", playground.Handler("RAG", "/graphql/"))

	return nil
}

func (srv *Server) openAiAPI(apiBasePath string) {
	srv.slog.Info("Registering openAI API", "basePath", apiBasePath)

	srv.mux.HandleFunc(fmt.Sprintf("POST %scompletions", apiBasePath), srv.completionHandler)
	srv.mux.HandleFunc(fmt.Sprintf("POST %schat/completions", apiBasePath), srv.chatCompletionHandler)
	srv.mux.HandleFunc(fmt.Sprintf("GET %smodels", apiBasePath), srv.modelsHandler)
	srv.mux.HandleFunc(fmt.Sprintf("GET %smodels/{model}", apiBasePath), srv.modelsHandler)

}
