package helper

import (
	"reflect"
	"os"
	"log"
)

func ReflectThrift(model interface{}, thrift interface{}) {

	thriftValue := reflect.ValueOf(thrift).Elem()
	thirftType := thriftValue.Type()

	modelValue := reflect.ValueOf(model).Elem()

	for i := 0; i < thriftValue.NumField(); i++ {
		name := thirftType.Field(i).Name
		thriftValue.Field(i).Set(modelValue.FieldByName(name))
	}
}

func GetEnv(key string) (value string) {
	value = os.Getenv(key)

	if value == "" {
		log.Println(key, " value is empty or not exist")
	}

	return
}