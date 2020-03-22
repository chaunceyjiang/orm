package session

import "orm/ormlog"

// 对原生的Begin进行封装
func (s *Session) Begin() (err error) {
	ormlog.Info("transaction begin")
	if s.tx, err = s.db.Begin(); err != nil {
		ormlog.Error(err)
		return err
	}
	return
}

func (s *Session) Commit() (err error) {
	ormlog.Info("transaction commit")
	if err = s.tx.Commit(); err != nil {
		ormlog.Error(err)
	}
	return
}

func (s *Session) Rollback() (err error) {
	ormlog.Info("transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		ormlog.Error(err)
	}
	return
}
