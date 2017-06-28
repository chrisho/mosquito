package helper

import (
	"os"
	"log"
	"strings"

	"github.com/chrisho/mosquito/utils"
	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc/grpclog"
)

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