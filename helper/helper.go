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
	"github.com/asaskevich/govalidator"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

const (
	YYYYMMDDHHIISS = "2006-01-02 15:04:05"
	YYYYMMDDHHII   = "2006-01-02 15:04"
	YYYYMMDDHH     = "2006-01-02 15"
	YYYYMMDD       = "2006-01-02"
)

func GetEnv(key string) (value string) {
	value = os.Getenv(key)

	if value == "" {
		log.Println(key, " value is empty or not exist")
	}

	return
}

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

// get time now unixtime == php time()
func UnixTime() int32 {
	return int32(time.Now().Unix())
}

// get time now date format YYYY-MM-DD HH:II:SS
func DateTime() string {
	return time.Now().Format(YYYYMMDDHHIISS)
}

// RFC3339 time string to format Y-m-d H:i,Y-m-d H ...
func RFC3339Time2Date(timeStr, dateLayout string) string {
	timestamp, _ := time.ParseInLocation(time.RFC3339, timeStr, time.Local)

	return timestamp.Format(dateLayout)
}

// change date to format Y-m-d H:i,Y-m-d H ...
func Date2Date(inputDateLayout, inputDate, resultDateLayout string) string {
	timestamp, _ := time.ParseInLocation(inputDateLayout, inputDate, time.Local)

	return timestamp.Format(resultDateLayout)
}

//  date to timestamp
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

// map map[int32]string
// map 的 key 必须从 0 开始，顺序排列
// sep is placed between elements in the resulting string.
func ImplodeMapInt32String(maps map[int32]string, sep string) (str string) {
	mapLen := int32(len(maps))

	for i := int32(0); i < mapLen; i++ {
		str += maps[i] + sep
	}

	str = strings.TrimRight(str, sep)

	return
}

// map map[int]string
// map 的 key 必须从 0 开始，顺序排列
// sep is placed between elements in the resulting string.
func ImplodeMapIntString(maps map[int]string, sep string) (str string) {
	mapLen := len(maps)

	for i := 0; i < mapLen; i++ {
		str += maps[i] + sep
	}

	str = strings.TrimRight(str, sep)

	return
}

//panic if s is not a slice
func ReverseSlice(s interface{}) {
	size := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

// Trim string space
func TrimStringSpace(str string) string {
	return strings.Trim(str, " ")
}

func GrpcError(c codes.Code, format string) error {
	return status.Errorf(c, format)
}

// table-表名，Struct-结构体的指针
// 构建 select 字段
func SqlSelectField(table string, Struct interface{}) string {
	field := ""

	structValue := reflect.ValueOf(Struct).Elem()
	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		name := structType.Field(i).Name
		field += table + "." + utils.SnakeString(name) + ","
	}

	return strings.TrimRight(field, ",")
}
