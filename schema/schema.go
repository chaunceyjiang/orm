// ORM核心 对象(object)与数据库表(table)的转化
package schema

import (
	"go/ast"
	"orm/dialect"
	"reflect"
)

// 数据库中有 表名 列名 列类型 约束
// 对象有    结构体    字段   TAG

// Field 字段<---> 列名
type Field struct {
	Name string // 对象的字段名
	Type string // 对象字段的字段类型
	Tag  string // 对象字段的Tag
}

// Schema 代表数据库中的一张表
type Schema struct {
	Model      interface{}       // 被映射对象Model
	Name       string            // 表名
	Fields     []*Field          // 被映射对象的所有字段
	FieldNames []string          // 被映射对象的所有字段名字
	fieldMap   map[string]*Field // 记录字段名与Field的关系，方便以后使用
}

// GetField 根据字段名字直接获取 Field
func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

// RecordValues 将一个obj的值取出来，转化未SQL的参数
func (s *Schema) RecordValues(obj interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(obj))
	var fieldValues []interface{}
	for _, field := range s.Fields {
		// dest 必须是结构体，然后根据字段名找到该字段
		// Interface() 将当前的值转化为一个 接口
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}

// Parse 将任意对象Model解析为数据库中的一张表
// 利用反射
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()

	// 下载这种反射方法，没有办法处理 dest是指针的情况，反射指针的Type是空值
	//modelType := reflect.ValueOf(dest).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}
	// 依次获取该对象的每一个字段
	for i := 0; i < modelType.NumField(); i++ {
		f := modelType.Field(i)

		if !f.Anonymous && ast.IsExported(f.Name) {
			// 不是匿名字段且可导出
			field := &Field{
				Name: f.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(f.Type))),
			}
			// 获取tag 中orm 标示
			if v, ok := f.Tag.Lookup("orm"); ok {
				field.Tag = v
			}

			schema.FieldNames = append(schema.FieldNames, f.Name)
			schema.Fields = append(schema.Fields, field)
			schema.fieldMap[f.Name] = field
		}
	}
	return schema
}
