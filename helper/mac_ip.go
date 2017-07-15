package helper

import (
	"strings"
	"google.golang.org/grpc/grpclog"
	"github.com/asaskevich/govalidator"
	"github.com/chrisho/mosquito/utils"
)

// get ip if it contains ipCutSet
// 对比指向IP是否存在于真实IP中
func ContainsIp(ipCutSet string) (r string) {
	r = "127.0.0.1"

	ips := utils.GetLocalIps()

	for _, ip := range ips {

		if ok := strings.Contains(ip, ipCutSet); ok {
			return ip
		}
	}

	return
}

// get server address from config file *.env
func GetServerAddress() (ipAddress string) {
	serverAddress := GetEnv("ServerAddress")
	serverPort := GetEnv("ServerPort")
	if serverPort == "" {
		grpclog.Fatal("no any server ports")
	}
	/*
	if serverAddress == "" {
		grpclog.Info("serverAddress is empty or not exist")
		serverAddress = "127.0.0.1"
	}
	*/
	if ok := govalidator.IsIP(serverAddress); serverAddress != "" && !ok {

		ipCutSet := strings.TrimRight(serverAddress, ".*")

		serverAddress = ContainsIp(ipCutSet)
	}

	ipAddress = serverAddress + ":" + serverPort

	return
}
