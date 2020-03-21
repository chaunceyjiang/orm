package session

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"orm/dialect"
	"testing"
)

func TestNew(t *testing.T) {
	db, err := sql.Open("sqlite3", "../orm.db")
	if err != nil {
		t.FailNow()
	}
	// sess := New(db)
	defer db.Close()
	d, _ := dialect.GetDialect("sqlite3")
	sess := New(db, d)
	raw := sess.Raw("SELECT name from User limit 1;").QueryRaw()
	var name string
	raw.Scan(&name)
	if name == "" {
		t.FailNow()
	}
}
