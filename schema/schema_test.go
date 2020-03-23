package schema

import (
	"orm/dialect"
	"orm/ormlog"
	"testing"
)

// 表名 User
type User struct {
	// 主键约束
	Name string `orm:"PRIMARY KEY"`
	Age  int
}

func (u *User) TableName() string {
	return "ns_user"
}
func TestParse(t *testing.T) {

	d, _ := dialect.GetDialect("sqlite3")
	s := Parse(&User{"foo", 18}, d)
	if s.Name == "User" {
		ormlog.Error("表名映射错误")
		t.Fail()
	}
	if s.Name != "ns_user"{
		ormlog.Error("表名映射错误")
		t.Fail()
	}
	if s.Fields[1].Name != "Age" {
		ormlog.Error("列名映射错误")
		t.Fail()
	}
	if s.GetField("Name").Tag != "PRIMARY KEY" {
		ormlog.Error("约束解析错误")
		t.Fail()
	}
}
