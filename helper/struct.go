package helper

import "reflect"

// structPointer 必须是 struct 的 指针
func StructFieldByNameIsValid(structPointer interface{}, fieldName string) bool {

	sElem := reflect.ValueOf(structPointer).Elem()

	return sElem.FieldByName(fieldName).IsValid()
}

// structPointer 必须是 struct 的 指针
func StructFieldByNameSetStringValue(structPointer interface{}, fieldName, newValue string) (oldValue string) {

	sElem := reflect.ValueOf(structPointer).Elem()
	reflectValue := sElem.FieldByName(fieldName)

	if ! reflectValue.IsValid() {
		return oldValue
	}
	oldValueInterface := reflectValue.Interface()

	if oldValueStr, ok := oldValueInterface.(string); ok {
		oldValue = oldValueStr
	} else {
		return oldValue
	}
	reflectValue.SetString(newValue)

	return oldValue
}
