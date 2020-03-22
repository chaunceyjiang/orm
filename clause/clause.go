package clause

import "strings"

type Clause struct {
	// 存在生成后的SQL语句
	sql map[Type]string
	// 存储SQL中参数
	sqlVars map[Type][]interface{}
}
// Set 存储SQL的子句
func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	// 根据传入的 SQL 类型，去调用相应SQL语句生成器
	sql, sqlVars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = sqlVars

}
// BuildSQL 构建完成SQL语句
func (c *Clause) BuildSQL(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}

	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			// 将各个子句，拼接成一个完整的SQL语句
			sqls = append(sqls, sql)
			// 取出每个子句对应的参数
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	// 生成一个完整的SQL语句，这里有回车符分割，方便日志观察
	return strings.Join(sqls, "\n"), vars
}
