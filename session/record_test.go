package session

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"orm/dialect"
	"testing"
)

func TestSession_Insert(t *testing.T) {
	type UserTest struct {
		// 主键约束
		Name string `orm:"PRIMARY KEY"`
		Age  int
	}
	db, err := sql.Open("sqlite3", "../orm.db")
	if err != nil {
		t.FailNow()
	}
	defer db.Close()
	d, _ := dialect.GetDialect("sqlite3")
	sess := New(db, d)

	sess = sess.Model(&UserTest{})


	if !sess.HasTable() {
		sess.CreateTable()
	}

	n, err := sess.Insert(&UserTest{
		Name: "foo",
		Age:  18,
	}, &UserTest{
		Name: "bar",
		Age:  25,
	})
	if n != 2 {
		t.FailNow()
	}
	if err != nil {
		t.FailNow()
	}

	sess.DropTable()
}
