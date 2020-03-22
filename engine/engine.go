// Engine 负责 用户与ORM进行交互
package engine

import (
	"database/sql"
	"orm/dialect"
	"orm/ormlog"
	"orm/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
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
	d, ok := dialect.GetDialect(driver)
	if !ok {
		ormlog.ErrorF("dialect %s Not Found", driver)
		return
	}
	e = &Engine{db: db, dialect: d}

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
	return session.New(e.db, e.dialect)
}

type TxFunc func(*session.Session) (interface{}, error)


// Transaction 用户只需要将所有操作都放在同一个回调函数中，就可以简单使用这个事物
func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if er := s.Begin(); er != nil {
		return
	}

	// 捕获异常，进行rollback
	defer func() {
		if p := recover(); p != nil {
			// 回滚数据
			_ = s.Rollback()

			// 异常继续传递
			panic(p)
		} else if err != nil {
			// 这里的 err 是回调函数产生的err,因此只回滚数据，不修改err
			_ = s.Rollback()
		} else {
			// 进入这个分之，代表err为nil ，因此可以复用err用来接收commit的err
			err = s.Commit()
		}
	}()
	// 回调函数
	return f(s)
}
