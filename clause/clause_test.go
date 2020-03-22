package clause

import (
	"orm/ormlog"
	"reflect"
	"testing"
)

func TestClause_BuildSQL(t *testing.T) {
	var clause Clause
	clause.Set(WHERE, "Name = ?", "foo")
	clause.Set(ORDERBY, "Age DESC")
	clause.Set(LIMIT, 3)
	clause.Set(SELECT, "User", []string{"Name", "Age"})

	sql, sqlVars := clause.BuildSQL(SELECT, WHERE, LIMIT, ORDERBY)
	if sql != `SELECT Name,Age FROM User 
WHERE Name = ?
LIMIT ?
ORDER BY Age DESC` {
		ormlog.Debug("sql build error")
		t.FailNow()
	}
	if !reflect.DeepEqual(sqlVars, []interface{}{"foo", 3}) {
		ormlog.Debug("sqlVars error")
		t.FailNow()
	}
}
