package helper

import (
	"reflect"
	"os"
	"log"
	"strings"
	"github.com/chrisho/mosquito/utils"
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

func ContainsIp(ipCutSet string) (r string) {
	r = "127.0.0.1"

	if ipCutSet == "" {
		return
	}

	ips := utils.GetLocalIps()

	for _, ip := range ips {

		if ok := strings.Contains(ip, ipCutSet); ok {
			return ip
		}
	}

	return
}