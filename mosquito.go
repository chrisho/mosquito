package mosquito

import (
	"github.com/chrisho/mosquito/helper"
	"github.com/chrisho/mosquito/zookeeper"
	"github.com/chrisho/mosquito/redis"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"net"
	"os"
	"strconv"
	"time"
	"log"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"github.com/sirupsen/logrus"
)

const envFile = "/config/conf.env"

var server *grpc.Server
var path string

func init() {
	path, _ = os.Getwd()

	err := godotenv.Load(path + envFile)
	if err != nil {
		grpclog.Error(err)
	}
}

func GetServer() *grpc.Server {
	// grpc 选项
	var opts []grpc.ServerOption
	if helper.GetEnv("SSL") == "true" {
		creds, err := credentials.NewServerTLSFromFile(path+"/config/server.crt", path+"/config/server.key")
		if err != nil {
			grpclog.Errorf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	// 注册interceptor
	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		grpcInterceptor(ctx, req, info)
		// 继续处理请求
		return handler(ctx, req)
	}
	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	// 实例化服务
	server = grpc.NewServer(opts...)

	return server
}

func RunServer() {

	_, err := zookeeper.RegMicroServer()

	if err != nil {
		grpclog.Fatal("reg server fail, ", err)
	}

	listenAddr := helper.GetServerAddress()

	log.Print("server address is ", listenAddr)

	listen, err := net.Listen("tcp", ":"+helper.GetEnv("ServerPort"))
	if err != nil {
		grpclog.Error(err)
	}

	err = server.Serve(listen)
	if err != nil {
		grpclog.Fatal(err)
	}
}

var prefixKey = "zk:"

func GetClientConn(serviceName string, userCredential ...*UserCredential) (client *grpc.ClientConn, err error) {

	db, _ := strconv.Atoi(helper.GetEnv("ZkRedisDb"))
	redisClient, err := redis.NewConnDB(db)
	if err != nil {
		return
	}
	defer redisClient.Close()
	redisServiceName := prefixKey + helper.GetEnv("ZkRootPath") + "/" + serviceName

	addr, err := redisClient.Get(redisServiceName).Result()
	if err != nil {
		return
	}

	var opts []grpc.DialOption
	var creds credentials.TransportCredentials

	if helper.GetEnv("SSL") == "true" {
		creds, err = credentials.NewClientTLSFromFile("config/server.crt", serviceName+".local")
		if err != nil {
			return
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	// 用户信息
	if len(userCredential) == 1{
		opts = append(opts, grpc.WithPerRPCCredentials(userCredential[0]))
	}

	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithTimeout(5*time.Second))
	client, err = grpc.Dial(addr, opts...)
	return
}

func GetLocalClientConn(serviceName string, userCredential ...*UserCredential) (conn *grpc.ClientConn, err error) {

	address := helper.GetServerAddress()

	var opts []grpc.DialOption
	var creds credentials.TransportCredentials

	if helper.GetEnv("SSL") == "true" {
		creds, err = credentials.NewClientTLSFromFile("config/server.crt", serviceName+".local")
		if err != nil {
			panic(err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	// 用户信息
	if len(userCredential) == 1{
		opts = append(opts, grpc.WithPerRPCCredentials(userCredential[0]))
	}

	grpclog.Info("get client server address is ", address)
	conn, err = grpc.Dial(address, opts...)
	if err != nil {
		grpclog.Error(err)
	}

	return
}

// 拦截器
func grpcInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo) {
	// 拦截消息
	var interceptorMessage string
	// metadata
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// authority
		if data, ok := md[":authority"]; ok {
			interceptorMessage += fmt.Sprintf("authority : %v ,\n ", data[0])
		}
		// user_id
		if data, ok := md["user_id"]; ok {
			interceptorMessage += fmt.Sprintf("user_id : %v ,\n ", data[0])
		}
		// username
		if data, ok := md["username"]; ok {
			interceptorMessage += fmt.Sprintf("user_name : %v ,\n ", data[0])
		}
	}
	// grpc method
	interceptorMessage += fmt.Sprintf("grpc_method : %v ,\n ", info.FullMethod)
	// grpc param
	interceptorMessage += fmt.Sprintf("grpc_method : %v ,\n ", req)
	// 日志
	logrus.Println(interceptorMessage)
}

// userCredential 用户认证
type UserCredential struct {
	User map[string]string
}

func (s UserCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	if s.User != nil {
		return s.User, nil
	}
	return map[string]string{
		"user_id":  "0",    // user_id(小写)
		"username": "test", // username(小写)
	}, nil
}

func (s UserCredential) RequireTransportSecurity() bool {
	return helper.GetEnv("SSL") == "true"
}

// 用户凭证
func NewUserCredential() *UserCredential {
	return new(UserCredential)
}
