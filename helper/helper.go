package helper

import "reflect"

// model User 映射到 Thrift User
func ReflectThrift(model interface{}, thrift interface{}) {

	thriftValue := reflect.ValueOf(thrift).Elem() // tU 的映射 Value
	thirftType := thriftValue.Type()

	modelValue := reflect.ValueOf(model).Elem()

	for i := 0; i < thriftValue.NumField(); i++ {
		name := thirftType.Field(i).Name
		thriftValue.Field(i).Set(modelValue.FieldByName(name))
	}
}
