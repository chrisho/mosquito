package helper

import (
	"reflect"
	"strings"
	"github.com/chrisho/mosquito/utils"
)

// 构建 select 字段
func SqlSelectField(tableName string, structPointer interface{}) string {
	field := ""

	structValue := reflect.ValueOf(structPointer).Elem()
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		name := structType.Field(i).Name
		field += tableName + "." + utils.SnakeString(name) + ","
	}

	return strings.TrimRight(field, ",")
}

// 构建 insert 字段
func SqlInsertField(tableName string, structPointer interface{}, withoutField map[string]bool) (insertField, prepareField string) {
	structValue := reflect.ValueOf(structPointer).Elem()
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		// 无需的字段
		if withoutField[structType.Field(i).Name] {
			continue
		}
		name := structType.Field(i).Name
		insertField += tableName + "." + utils.SnakeString(name) + ","
		prepareField += "?,"
	}
	insertField = strings.TrimRight(insertField, ",")
	prepareField = strings.TrimRight(prepareField, ",")

	return
}
