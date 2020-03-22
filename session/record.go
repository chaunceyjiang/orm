package session

import (
	"errors"
	"orm/clause"
	"reflect"
)

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

// Find 根据传入的对象，查找相应的表，然后将数据依次填充到对应字段
// Find([]User)
func (s *Session) Find(objs interface{}) error {
	if reflect.Indirect(reflect.ValueOf(objs)).Kind() != reflect.Slice {
		return errors.New("Unsupported! expect Slice ")
	}
	objSlice := reflect.Indirect(reflect.ValueOf(objs))
	// Elem 获取Slice 中的元素
	objType := objSlice.Type().Elem()

	// 根据objs列表中反射出列表中的类型,然后创建相应的schema
	table := s.Model(reflect.New(objType).Elem().Interface()).RefTable()

	// _select(tableName string, fields []string)
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)

	sql, sqlVars := s.clause.BuildSQL(clause.SELECT)

	raws, err := s.Raw(sql, sqlVars...).QueryRaws()
	if err != nil {
		return err
	}

	for raws.Next() {
		// 创建一个全部是空值的新obj
		obj := reflect.New(objType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			// 获取新obj的每个字段的地址
			values = append(values, obj.FieldByName(name).Addr().Interface())
		}

		// 查询数据库，同时通过地址给每个字段赋值
		if err := raws.Scan(values...); err != nil {
			return err
		}
		// 将新的Obj 重新追加到 传递的Slice中
		objSlice.Set(reflect.Append(objSlice, obj))
	}

	return raws.Close()
}
