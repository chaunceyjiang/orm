package session

import (
	"fmt"
	"orm/ormlog"
	"orm/schema"
	"reflect"
	"strings"
)

// Model 用来更新Session中的refTable ， 实现同一个session操作不同的对象
// 将一个对象转化为一张表
func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {

		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

// RefTable 用来获取Session中refTable的值
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		ormlog.Error("Model is not set")
	}
	return s.refTable
}

// CreateTable
func (s *Session) CreateTable() error {
	// 创建表的SQL在主流数据库中都适用，因此不需要封装成dialect
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		// 拼接SQL
		// 列名 类型 约束
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}

	createSql := fmt.Sprintf(
		//  拼接创建SQL语句
		"CREATE TABLE %s (%s);",
		table.Name,
		strings.Join(columns, ","),
	)
	ormlog.Debug(createSql)
	// 执行创建语句
	_, err := s.Raw(createSql).Exec()
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s;", s.RefTable().Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {

	sql, sqlVars := s.dialect.TableExistSQL(s.RefTable().Name)
	raw := s.Raw(sql, sqlVars...).QueryRaw()
	var tableName string
	_ = raw.Scan(&tableName)
	return tableName == s.RefTable().Name
}
