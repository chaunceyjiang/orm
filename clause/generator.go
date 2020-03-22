// 实现SELECT中各个字句的生成器
package clause

import (
	"fmt"
	"strings"
)

/*
SELECT col1, col2, ...
    FROM table_name
    WHERE [ conditions ]
    GROUP BY col1
    HAVING [ conditions ]

*/

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

// _select 生成 SELECT语句 SELECT $field FROM $tableName
func _select(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %s FROM %s ", fields, tableName), []interface{}{}

}

// _insert 生成insert 语句  INSERT INTO $tableName($fields)
func _insert(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%s) ", tableName, fields), []interface{}{}
}

// _values 生成values 语句  VALUES ($value1),($value2),...
func _values(values ...interface{}) (string, []interface{}) {
	var sql strings.Builder
	var sqlVars []interface{}

	sql.WriteString("VALUES ")

	var buildStr string
	for i, value := range values {
		v := value.([]interface{})
		if buildStr == "" {
			// 生成VALUES 的占位符
			buildStr = genBindVars(len(v))
		}
		// 生成 $(value1)
		sql.WriteString(fmt.Sprintf("(%v)", buildStr))

		if i+1 != len(values) {
			// 如果不是最后一个，则增加 ',' 分割
			// 生成 $(value1),
			sql.WriteString(", ")
		}

		sqlVars = append(sqlVars, v...)
	}

	return sql.String(), sqlVars
}

// genBindVars 构造参数占位符
func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

// _where 生成where 语句
func _where(values ...interface{}) (string, []interface{}) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars

}

// _limit 生成LIMIT 语句
func _limit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

// _delete 删除
func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

// _update 更新  第一个参数表名，第二个是map ，待更新的键值对
func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	kv := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}
	for k, v := range kv {
		keys = append(keys, k+" = ? ")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ",")), vars
}

func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"COUNT(*)"})
}
func init() {
	generators = make(map[Type]generator)
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[SELECT] = _select
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[ORDERBY] = _orderBy
	generators[UPDATE] = _update
	generators[COUNT] = _count
	generators[DELETE] = _delete
}
