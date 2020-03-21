package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type sqlite3 struct {
}

// DataTypeOf 返回sqlite3数据库列席 参考 https://www.runoob.com/sqlite/sqlite-data-types.html
func (s *sqlite3) DataTypeOf(typ reflect.Value) string {
	// 反射，判断该值的类型
	switch typ.Kind() {
	case reflect.Bool:
		// Sqlite没有单独的布尔存储类型，它使用INTEGER作为存储类型，0为false，1为true
		return "INTEGER" // 返回sqlite3数据库类型
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "INTEGER"
	case reflect.Int64, reflect.Uint64:
		return "BIGINT"
	case reflect.String:
		return "TEXT"
	case reflect.Float32, reflect.Float64:
		return "REAL"
	case reflect.Array, reflect.Slice:
		return "BLOB"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "DATETIME"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%d)", typ.Type().Name(), typ.Kind()))
}

// TableExistSQL 返回判断表是否存在的SQL
func (s *sqlite3) TableExistSQL(tableName string) (sql string, sqlVars []interface{}) {
	sqlVars = []interface{}{tableName}
	sql = "SELECT tbl_name FROM sqlite_master WHERE type='table' AND tbl_name=?"
	return
}

// 判断该数据库是否实现了Dialect接口
var _ Dialect = (*sqlite3)(nil)


// 自动注册
func init()  {
	RegisterDialect("sqlite3", &sqlite3{})
}