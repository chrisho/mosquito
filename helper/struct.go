package helper

import "reflect"

// structPointer 必须是 struct 的 指针
func StructFieldByNameIsValid(structPointer interface{}, fieldName string) bool {

	sElem := reflect.ValueOf(structPointer).Elem()

	return sElem.FieldByName(fieldName).IsValid()
}
