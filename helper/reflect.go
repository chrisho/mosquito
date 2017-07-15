package helper

import "reflect"

// model、dest结构体的指针
// 两个相似的结构体赋值
func Reflect(model interface{}, dest interface{}) {

	destValue := reflect.ValueOf(dest).Elem()
	destType := destValue.Type()

	modelValue := reflect.ValueOf(model).Elem()

	for i := 0; i < destValue.NumField(); i++ {
		name := destType.Field(i).Name
		destValue.Field(i).Set(modelValue.FieldByName(name))
	}
}

//panic if s is not a slice
func ReverseSlice(s interface{}) {
	size := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}