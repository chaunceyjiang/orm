// 屏蔽不同数据库的差异
package dialect

import "reflect"

var dialectMap = map[string]Dialect{}

// Dialect 实现该接口的数据库都能被ORM接管
type Dialect interface {
	// DataTypeOf 用于将Go的类型转化为该数据库的数据库类型
	DataTypeOf(typ reflect.Value) string
	// 返回某个表是否存在的SQL ，每个数据库的SQL不一样相同
	TableExistSQL(tableName string) (string, []interface{})
}

// RegisterDialect 注册该数据库到ORM中
func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}
// GetDialect 获取该数据库的实例
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectMap[name]
	return
}
