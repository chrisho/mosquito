package helper

import (
	"log"
	"os"
	"strings"

	"github.com/chrisho/mosquito/utils"
	"google.golang.org/grpc/grpclog"
	"reflect"
	"strconv"
	"time"
)

const (
	YYYYMMDDHHIISS = "2006-01-02 15:04:05"
	YYYYMMDDHHII = "2006-01-02 15:04"
	YYYYMMDDHH = "2006-01-02 15"
	YYYYMMDD = "2006-01-02"
)

func GetEnv(key string) (value string) {
	value = os.Getenv(key)

	if value == "" {
		log.Println(key, " value is empty or not exist")
	}

	return
}

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
// 返回监听IP和网卡IP，假设监听IP不存在则网卡IP为127.0.0.1
func GetServerAddress() (ipAddress string, inetIP string) {
	serverAddress := GetEnv("ServerAddress")
	serverPort := GetEnv("ServerPort")

	if serverAddress == "" {
		grpclog.Info("serverAddress is empty or not exist")
		serverAddress = "127.0.0.1"
	}

	ipCutSet := strings.TrimRight(serverAddress, ".*")  // IP前叙匹配方式
	inetIP = ContainsIp(ipCutSet)

	if serverPort == "" {
		ipAddress = serverAddress
	} else {
		ipAddress = serverAddress + ":" + serverPort
		inetIP = inetIP + ":" + serverPort
	}

	return
}

// RFC3339 time string to format Y-m-d H:i
func RFC3339Time2Date(timeStr, dateLayout string) string {
	timestamp, _ := time.ParseInLocation(time.RFC3339, timeStr, time.Local)

	return timestamp.Format(dateLayout)
}

// timestamp to format Y-m-d H:i:s
func Date2Timestamp(dateLayout, date string) int64 {
	localTime, _ := time.ParseInLocation(dateLayout, date, time.Local)

	return localTime.Unix()
}

// timestamp to format Y-m-d H:i:s
func Timestamp2Date(timestamp int64, dateLayout string) string {
	return time.Unix(timestamp, 0).Format(dateLayout)
}

// string to int32
func String2Int32(str string) int32 {
	int64Str, _ := strconv.Atoi(str)

	return int32(int64Str)
}
