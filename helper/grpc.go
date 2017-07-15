package helper

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// grpc error
func GrpcError(c codes.Code, format string) error {
	return status.Errorf(c, format)
}
