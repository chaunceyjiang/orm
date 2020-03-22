// 用于实现与数据库的交互
package session

import (
	"database/sql"
	"orm/clause"
	"orm/dialect"
	"orm/ormlog"
	"orm/schema"
	"strings"
)

type Session struct {
	// db 原生数据库
	db *sql.DB
	// sql 存储sql语句
	sql strings.Builder
	// sqlVars 存储SQL语句中的变量
	sqlVars []interface{}

	dialect dialect.Dialect // 数据库dialect,用来屏蔽不同数据库之间的差异

	refTable *schema.Schema // 数据库表与对象Model 的连接

	clause clause.Clause // 存储未编译的SQL

	tx *sql.Tx // 支持事物
}

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// 检查是否实现CommonDB接口
var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

// New 返回一个Session 连接
func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{db: db, dialect: dialect}
}

// DB 获取原生DB
// CommonDB 返回一个公共DB
func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

// Clear 清除Session中存储的SQL
func (s *Session) Clear() {
	s.sqlVars = nil
	s.sql.Reset()
	s.clause = clause.Clause{}
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Exec 执行Session 中存储的SQL 对sql.DB.Exec 的封装
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	ormlog.InfoF("SQL %s Params %v", s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		ormlog.Error(err)
	}
	return
}

func (s *Session) QueryRaw() *sql.Row {
	defer s.Clear()
	ormlog.InfoF("SQL %s Params %v", s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) QueryRaws() (rows *sql.Rows, err error) {
	defer s.Clear()
	ormlog.InfoF("SQL %s Params %v", s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		ormlog.Error(err)
	}
	return
}
