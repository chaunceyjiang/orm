package session

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"orm/dialect"
	"orm/ormlog"
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

func TestSession_Find(t *testing.T) {
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
	defer sess.DropTable()
	sess.Insert(&UserTest{
		Name: "foo",
		Age:  18,
	}, &UserTest{
		Name: "bar",
		Age:  25,
	})

	var users []UserTest
	err = sess.Find(&users)

	if err != nil {
		ormlog.Error(err)
		t.FailNow()
	}
	if users[0].Name != "foo" {
		ormlog.Error("find name error")
		t.FailNow()
	}

	if users[1].Age != 25 {
		ormlog.Error("find age error")
		t.FailNow()
	}

}
