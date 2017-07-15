package helper

import (
	"reflect"
	"strings"
	"github.com/chrisho/mosquito/utils"
)

// table-表名，Struct-结构体的指针
// 构建 select 字段
func SqlSelectField(table string, Struct interface{}) string {
	field := ""

	structValue := reflect.ValueOf(Struct).Elem()
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		name := structType.Field(i).Name
		field += table + "." + utils.SnakeString(name) + ","
	}

	return strings.TrimRight(field, ",")
}
