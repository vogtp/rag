package testhelper

import (
	"log/slog"
	"testing"

	"github.com/spf13/viper"
	"github.com/vogtp/rag/pkg/cfg"
	"github.com/vogtp/rag/pkg/logger"
	"github.com/vogtp/rag/pkg/usercfg"
)

// GetDB returns a test db with test users
func GetDB(t *testing.T) (*usercfg.DataBase, *slog.Logger) {
	viper.AddConfigPath("../..")
	cfg.Parse()
	viper.Set(cfg.LogJson, true)
	slog := logger.Create(slog.LevelWarn).With("type", "testLogger")
	d, err := usercfg.Create(t.Context(), slog, "")
	if err != nil {
		t.Fatalf("could not construct dtatabase: %v", err)
	}
	if err := d.Add(t.Context(), &User1); err != nil {
		t.Fatalf("cannot add inital test user1: %v", err)
	}
	if err := d.Add(t.Context(), &User2); err != nil {
		t.Fatalf("cannot add inital test user2: %v", err)
	}
	AllUsers = make([]usercfg.User, 0, 2)
	AllUsers = append(AllUsers, User1, User2)
	return d, slog
}

// CompareUser compairs users with its dependencies
// triggers t.Errorf
func CompareUser(t *testing.T, tu *usercfg.User, u *usercfg.User) {
	if tu.Name != u.Name {
		t.Errorf("user not the correct name: want %s have %s", tu.Name, u.Name)
	}
	if tu.APIKey != u.APIKey {
		t.Errorf("user not the correct APIKey: want %s have %s", tu.APIKey, u.APIKey)
	}
	if len(tu.Collections) != len(u.Collections) {
		t.Errorf("user not the correct number of collections: want %v have %v", len(tu.Collections), len(u.Collections))
	}
	for i, tc := range tu.Collections {
		CompareCollection(t, &tc, &u.Collections[i])
	}
}

func CompareCollection(t *testing.T, tc *usercfg.Collection, c *usercfg.Collection) {
	if tc.DisplayName != c.DisplayName {
		t.Errorf("collection not the correct DisplayName: want %s have %s", tc.DisplayName, c.DisplayName)
	}
	if tc.CollectionName != c.CollectionName {
		t.Errorf("collection not the correct CollectionName: want %s have %s", tc.CollectionName, c.CollectionName)
	}
	if tc.APIKey != c.APIKey {
		t.Errorf("collection not the correct APIKey: want %s have %s", tc.APIKey, c.APIKey)
	}
	CompareSource(t, &tc.Source, &c.Source)
}
func CompareSource(t *testing.T, ts *usercfg.SourceSystem, s *usercfg.SourceSystem) {
	if ts.Name != s.Name {
		t.Errorf("sourcesystem not the correct Name: want %s have %s", ts.Name, s.Name)
	}
	if ts.Type != s.Type {
		t.Errorf("sourcesystem not the correct Type: want %s have %s", ts.Type, s.Type)
	}
	if ts.URL != s.URL {
		t.Errorf("sourcesystem not the correct URL: want %s have %s", ts.URL, s.URL)
	}
	if ts.Key != s.Key {
		t.Errorf("sourcesystem not the correct Key: want %s have %s", ts.Key, s.Key)
	}
	if ts.Parts != s.Parts {
		t.Errorf("sourcesystem not the correct Parts: want %s have %s", ts.Parts, s.Parts)
	}
}

var (
	AllUsers []usercfg.User
	User1    = usercfg.User{
		Name:        "TestUserName1",
		APIKey:      "TestAPIKey1",
		Collections: CollectionsUser1,
	}
	CollectionsUser1 = []usercfg.Collection{
		{
			DisplayName:    "TestDisplayNameUser1",
			CollectionName: "TestCollectionNameUser1",
			APIKey:         "TestAPIKeyUser1",
			Source:         SourceUser1,
		},
	}
	SourceUser1 = usercfg.SourceSystem{
		Name:  "TestSourceNameUser1",
		Type:  usercfg.SourceConfluence,
		URL:   "TestURLUser1",
		Key:   "TestKeyUser1",
		Parts: "TestPartsUser1",
	}
	User2 = usercfg.User{
		Name:        "TestUserName2",
		APIKey:      "TestAPIKey2",
		Collections: CollectionsUser2,
	}
	CollectionsUser2 = []usercfg.Collection{
		{
			DisplayName:    "TestDisplayNameUser2",
			CollectionName: "TestCollectionNameUser2",
			APIKey:         "TestAPIKeyUser2",
			Source:         SourceUser2,
		},
	}
	SourceUser2 = usercfg.SourceSystem{
		Name:  "TestSourceNameUser2",
		Type:  usercfg.SourceConfluence,
		URL:   "TestURLUser2",
		Key:   "TestKeyUser2",
		Parts: "TestPart1User2,TestPart2User2",
	}
)
