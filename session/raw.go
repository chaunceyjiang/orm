// 用于实现与数据库的交互
package session

import (
	"database/sql"
	"orm/ormlog"
	"strings"
)

type Session struct {
	// db 原生数据库
	db *sql.DB
	// sql 存储sql语句
	sql strings.Builder
	// sqlVars 存储SQL语句中的变量
	sqlVars []interface{}
}

// New 返回一个Session 连接
func New(db *sql.DB) *Session {
	return &Session{db: db}
}

// DB 获取原生DB
func (s *Session) DB() *sql.DB {
	return s.db
}

// Clear 清除Session中存储的SQL
func (s *Session) Clear() {
	s.sqlVars = nil
	s.sql.Reset()
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
	return s.DB().QueryRow(s.sql.String(),s.sqlVars...)
}

func (s *Session) QueryRaws()(rows *sql.Rows,err error) {
	defer s.Clear()
	ormlog.InfoF("SQL %s Params %v", s.sql.String(), s.sqlVars)
	if rows,err = s.DB().Query(s.sql.String(),s.sqlVars...);err!=nil{
		ormlog.Error(err)
	}
	return
}