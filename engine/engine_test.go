package engine

import (
	_ "github.com/mattn/go-sqlite3"
	"orm/ormlog"
	"testing"
)

func TestNewEngine(t *testing.T) {
	engine, err := NewEngine("sqlite3", "../orm.db")
	if err != nil {
		t.FailNow()
	}
	raw := engine.NewSession().Raw("SELECT name from User limit ?", 1).QueryRaw()
	var name string
	raw.Scan(&name)
	if name == "" {
		t.FailNow()
	}
	ormlog.Debug(name)
}
