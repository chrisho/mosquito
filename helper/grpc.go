package helper

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime"
	"fmt"
	_ "github.com/chrisho/mosquito/alilog"
	"github.com/sirupsen/logrus"
)

// grpc error
func GrpcError(c codes.Code, format string) error {
	_, file, line, _ := runtime.Caller(1)

	// 错误代码 200 、 422 、自定义验证错误（10000+）， 不记录日志
	if c == 200 || c == 422 || c > 9999 {
		return status.Errorf(c, format)
	}

	file = fmt.Sprintf("error file : %v ( code : %d) ( line : %d )", file, int(c), line)
	logrus.Error(file)
	logrus.Error("error info :" + format)

	return status.Errorf(c, format)
}
