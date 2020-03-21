package session

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"orm/dialect"
	"orm/ormlog"
	"testing"
)

func TestSession_CreateTable(t *testing.T) {
	// 表名 User
	type UserTest struct {
		// 主键约束
		Name string `orm:"PRIMARY KEY"`
		Age  int
	}
	db, err := sql.Open("sqlite3", "../orm.db")
	if err != nil {
		t.FailNow()
	}
	// sess := New(db)
	defer db.Close()

	d, _ := dialect.GetDialect("sqlite3")
	sess := New(db, d)

	sess.Model(&UserTest{
		Name: "foo",
		Age:  18,
	})
	if sess.HasTable() == true {
		ormlog.Error("UserTest EXISTS!")
		t.FailNow()
	}

	sess.CreateTable()

	if sess.HasTable() == true {
		ormlog.Debug("create table success")
	}

	sess.DropTable()

	if sess.HasTable() == true {
		ormlog.Error(err)
		t.FailNow()
	}
}
