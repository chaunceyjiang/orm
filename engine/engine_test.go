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

func TestNewEngine2(t *testing.T) {
	engine, err := NewEngine("sqlite3", "../orm.db")
	if err != nil {
		ormlog.Error(err)
		t.FailNow()
	}
	defer engine.Close()

	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS Student;").Exec()
	_, _ = s.Raw("CREATE TABLE Student(Name text);").Exec()
	_, _ = s.Raw("CREATE TABLE Student(Name text);").Exec()

	result, _ := s.Raw("INSERT INTO Student(`Name`) values (?), (?)", "Foo", "Bar").Exec()
	count, _ := result.RowsAffected()
	ormlog.Debug(count)
}

func TestEngine_Transaction(t *testing.T) {

}