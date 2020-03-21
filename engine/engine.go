// Engine 负责 用户与ORM进行交互
package engine

import (
	"database/sql"
	"orm/ormlog"
	"orm/session"
)

type Engine struct {
	db *sql.DB
}

// NewEngine
func NewEngine(driver string, source string) (e *Engine, err error) {
	var db *sql.DB
	db, err = sql.Open(driver, source)
	if err != nil {
		ormlog.Error(err)
		return
	}
	// 测试连接是否正常
	if err = db.Ping(); err != nil {
		ormlog.Error(err)
		return
	}

	e = &Engine{db: db}

	ormlog.InfoF("Connect database %s success", driver)

	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		ormlog.ErrorF("Close failed %v", err)
		return
	}
	ormlog.Info("Close database success")
}

// NewSession 创建一个与数据库交互的Session
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db)
}
