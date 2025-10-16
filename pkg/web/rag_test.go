package web

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vogtp/rag/internal/testhelper"
	"github.com/vogtp/rag/pkg/web/bearer"
	"github.com/vogtp/rag/pkg/web/oidc"
)

func TestServer_BearerToken_AllFromRequest(t *testing.T) {

	db, slog := testhelper.GetDB(t)
	srv, err := New(t.Context(), slog, db)
	if err != nil {
		t.Fatalf("could not create web server: %v", err)
	}
	tu := testhelper.User1
	testRequest := httptest.NewRequest(http.MethodGet, "/ui", nil)
	testRequest.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tu.APIKey))
	h := bearer.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testRequest = r
	})
	h.ServeHTTP(httptest.NewRecorder(), testRequest)
	bt, ok := bearer.Get(testRequest)
	if !ok {
		t.Errorf("no bearer token found in request")
	}
	if bt != tu.APIKey {
		t.Errorf("wrong bearer token found: got %q want %q", bt, tu.APIKey)
	}
	rag := srv.ragMgr.FromRequest(t.Context(), testRequest)
	models := rag.Models(t.Context())
	if len(models) < 1 {
		t.Errorf("found no models")
	}
	// FIXME what is the difference between Models() and ListCollections()
	// -> cleanup interface
	// for _, c := range models {
	// 	t.Error(c.GetName())
	// }
}

func TestServer_OIDCUser_AllFromRequest(t *testing.T) {

	db, slog := testhelper.GetDB(t)
	srv, err := New(t.Context(), slog, db)
	if err != nil {
		t.Fatalf("could not create web server: %v", err)
	}
	tu := testhelper.User1
	oidc.TESTINGUsernameDoNotUse(tu.Name)
	// passing nil as request allows test usernames
	un, err := oidc.UserName(nil)
	if err != nil {
		t.Errorf("cannot get username from requst: %v", err)
	}
	if un != tu.Name {
		t.Errorf("wrong username found: got %q want %q", un, tu.Name)
	}
	rag := srv.ragMgr.FromRequest(t.Context(), nil)
	models := rag.Models(t.Context())
	if len(models) < 1 {
		t.Errorf("found no models")
	}
	// FIXME what is the difference between Models() and ListCollections()
	// -> cleanup interface
	// for _, c := range models {
	// 	t.Error(c.GetName())
	// }
}
