package session

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"orm/dialect"
	"orm/ormlog"
	"testing"
)

type UserTest struct {
	// 主键约束
	Name string `orm:"PRIMARY KEY"`
	Age  int
}

func testHelper(t *testing.T) *Session {
	db, err := sql.Open("sqlite3", "../orm.db")
	if err != nil {
		t.FailNow()
	}

	d, _ := dialect.GetDialect("sqlite3")
	sess := New(db, d)

	sess = sess.Model(&UserTest{})

	if sess.HasTable() {
		sess.DropTable()
	}
	sess.CreateTable()
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

	return sess
}
func TestSession_Insert(t *testing.T) {
	testHelper(t)
}

func TestSession_Find(t *testing.T) {
	sess := testHelper(t)

	var users []UserTest
	err := sess.Find(&users)

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

func TestSession_Update(t *testing.T) {
	sess := testHelper(t)
	sess.Where("Age = ?", 25).Update("Name", "barr")
	var user []UserTest
	sess.Where("Age = ?", 25).Find(&user)
	if user[0].Name != "barr" {
		ormlog.Error("Update error")
		t.FailNow()
	}
}

func TestSession_First(t *testing.T) {
	sess := testHelper(t)
	var user UserTest
	sess.Where("Age = ?",25).First(&user)
	if user.Name != "bar" {
		ormlog.Error("First error")
		t.FailNow()
	}
}