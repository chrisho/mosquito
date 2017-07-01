package helper

import (
	"log"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/chrisho/mosquito/utils"
	"google.golang.org/grpc/grpclog"
	"reflect"
)

func GetEnv(key string) (value string) {
	value = os.Getenv(key)

	if value == "" {
		log.Println(key, " value is empty or not exist")
	}

	return
}

// get ip if it contains ipCutSet
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

// two struct be similar, reflect model struct to destination struct
func Reflect(model interface{}, dest interface{}) {

	destValue := reflect.ValueOf(dest).Elem()
	destType := destValue.Type()

	modelValue := reflect.ValueOf(model).Elem()

	for i := 0; i < destValue.NumField(); i++ {
		name := destType.Field(i).Name
		destValue.Field(i).Set(modelValue.FieldByName(name))
	}
}

// get server address from config file *.env
func GetServerAddress() (ipAddress string) {
	serverAddress := GetEnv("ServerAddress")
	serverPort := GetEnv("ServerPort")

	if serverAddress == "" {
		grpclog.Info("serverAddress is empty or not exist")
		serverAddress = "127.0.0.1"
	}

	if ok := govalidator.IsIP(serverAddress); !ok {

		ipCutSet := strings.TrimRight(serverAddress, ".*")

		serverAddress = ContainsIp(ipCutSet)
	}

	if serverPort == "" {
		ipAddress = serverAddress
	} else {
		ipAddress = serverAddress + ":" + serverPort
	}

	return
}
