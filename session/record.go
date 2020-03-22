package session

import "orm/clause"

// Insert 将对象的值转换为SQL，然后将数据插入数据库表中
// u1 := &User{Name: "Sam", Age: 25} s.Insert(u1, u2, ...)
func (s *Session) Insert(objs ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, obj := range objs {
		table := s.Model(obj).RefTable()
		// 构造INSERT语句
		// _insert(tableName string fields []string)
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(obj))
	}

	s.clause.Set(clause.VALUES, recordValues...)

	// 构建 INSERT 语句
	sql, sqlVars := s.clause.BuildSQL(clause.INSERT, clause.VALUES)

	// Raw 接收SQL和参数
	result, err := s.Raw(sql, sqlVars...).Exec()
	// 返回受影响的行数，和错误
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
