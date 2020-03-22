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

	// 将链式调用前面的一起build
	sql, sqlVars := s.clause.BuildSQL(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)

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

// Update 支持两张传参方式 1. map[string]interface{} 2. kv list: "Name","for","Age",16
// 依赖where语法
func (s *Session) Update(kv ...interface{}) (int64, error) {
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		// 如果不是map方式，
		m = make(map[string]interface{})

		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}

	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)

	sql, sqlVars := s.clause.BuildSQL(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, sqlVars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete 删除语句
// 依赖where语法
func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.BuildSQL(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Count records with where clause
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.BuildSQL(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRaw()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

// Limit 结合其他语法一起使用
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

// Where 结合其他语法一起使用
func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

// OrderBy 结合其他语法一起使用
func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	// 根据传入的对象，创建一个对应的S
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return errors.New("Not Found ")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
