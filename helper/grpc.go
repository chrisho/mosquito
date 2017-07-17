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

	file = fmt.Sprintf("error file : %v ( line : %d )", file, line)

	logrus.Error(file)

	return status.Errorf(c, format)
}
