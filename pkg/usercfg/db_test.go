package usercfg_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/vogtp/rag/internal/testhelper"
	"github.com/vogtp/rag/pkg/usercfg"
)

func TestDataBase_QueryUser(t *testing.T) {
	ctx := t.Context()
	db, _ := testhelper.GetDB(t)
	usrs, err := db.Users(ctx)
	if err != nil {
		t.Fatalf("query users: %v", err)
	}
	if len(usrs) != len(testhelper.AllUsers) {
		t.Errorf("user count not correct: have: %v want %v", len(usrs), len(testhelper.AllUsers))
	}
	testUsrs := []usercfg.User{testhelper.User1, testhelper.User2}
	for _, tu := range testUsrs {
		u, err := db.User(ctx, tu.Name)
		if err != nil {
			t.Errorf("cannot query user %q: %v", tu.Name, err)
		}
		testhelper.CompareUser(t, &tu, u)
	}
}

func TestDataBase_QueryAPIKey(t *testing.T) {
	ctx := t.Context()
	db, _ := testhelper.GetDB(t)
	usrs, err := db.Users(ctx)
	if err != nil {
		t.Fatalf("query users: %v", err)
	}
	if len(usrs) != len(testhelper.AllUsers) {
		t.Errorf("user count not correct: have: %v want %v", len(usrs), len(testhelper.AllUsers))
	}
	testUsrs := []usercfg.User{testhelper.User1, testhelper.User2}
	for _, tu := range testUsrs {
		u, err := db.UserByAPIKey(ctx, tu.APIKey)
		if err != nil {
			t.Errorf("cannot query user %q: %v", tu.Name, err)
		}
		testhelper.CompareUser(t, &tu, &u[0])
	}
}

func TestDataBase_Add(t *testing.T) {
	ctx := t.Context()
	db, _ := testhelper.GetDB(t)
	tu, err := db.User(ctx, testhelper.User1.Name)
	if err != nil {
		t.Errorf("cannot query user %q: %v", testhelper.User1.Name, err)
	}
	partsFmt := "newParts%v,newParts%v"
	dnFmt := "newDisplayName%v"
	for i := range tu.Collections {
		tu.Collections[i].Displayname = fmt.Sprintf(dnFmt, i)
		tu.Collections[i].Source.Parts = fmt.Sprintf(partsFmt, i, i)
	}
	if err := db.Add(ctx, tu); err != nil {
		t.Fatalf("cannot update user %q: %v", tu.Name, err)
	}
	u, err := db.User(ctx, tu.Name)
	if err != nil {
		t.Errorf("cannot query user %q: %v", tu.Name, err)
	}
	testhelper.CompareUser(t, tu, u)
	for i, c := range tu.Collections {
		if c.Displayname != fmt.Sprintf(dnFmt, i) {
			t.Errorf("collection display name not correct: %v -> %s", i, c.Displayname)
		}
		if c.Source.Parts != fmt.Sprintf(partsFmt, i, i) {
			t.Errorf("collection source parts not correct: %v -> %s", i, c.Source.Parts)
		}
	}
}

func TestDataBase_CollectionNextUpdate(t *testing.T) {
	ctx := t.Context()
	db, _ := testhelper.GetDB(t)
	cols, err := db.CollectionsToUpdate(ctx, time.Now())
	if err != nil {
		t.Fatalf("query users: %v", err)
	}
	if len(cols) != 1 {
		t.Errorf("Wrong number of collections to update found: %d want:1", len(cols))
	}
	if cols[0].Displayname != "TestDisplayNameUser1" {
		t.Errorf("Wrong collection found: %q", cols[0].Collectionname)
	}
}
