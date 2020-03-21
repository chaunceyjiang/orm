package session

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestNew(t *testing.T) {
	db, err := sql.Open("sqlite3", "../orm.db")
	if err != nil {
		t.FailNow()
	}
	sess := New(db)
	raw := sess.Raw("SELECT name from User limit 1;").QueryRaw()
	var name string
	raw.Scan(&name)
	if name == "" {
		t.FailNow()
	}
}
