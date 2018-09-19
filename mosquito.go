package mosquito

import (
	"fmt"
	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/chrisho/mosquito/alilog"
	"github.com/chrisho/mosquito/helper"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const envFile = "/config/conf.env"

var (
	server    *grpc.Server
	path      string
	port      = ":50051"
	debugPort = ""
	runPort   = ":50051"
)

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
	if strings.ToLower(helper.GetEnv("SSL")) == "true" {
		certFile := helper.GetEnv("SSLCertFile")
		keyFile := helper.GetEnv("SSLKeyFile")
		creds, err := credentials.NewServerTLSFromFile(path+"/"+certFile, path+"/"+keyFile)
		if err != nil {
			grpclog.Errorf("Failed to generate credentials %v", err)
		}
		log.Println("set grpc with credentials")
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

	isDebug()

	listen, err := net.Listen("tcp", runPort)
	if err != nil {
		grpclog.Error(err)
	}

	log.Println("listening TCP: " + runPort)
	err = server.Serve(listen)
	if err != nil {
		grpclog.Fatal(err)
	}
}

func GetClientConn(serviceName string, userCredential ...*UserCredential) (client *grpc.ClientConn, err error) {

	runPort = port

	serviceName = helper.ConvertUnderlineToWhippletree(serviceName)
	host := serviceName + helper.GetEnv("SSLSuffixServerName")
	address := serviceName + helper.GetEnv("ClusterSuffixDomain")
	return setClientConn(host, address, userCredential)
}

func GetLocalClientConn(serviceName string, userCredential ...*UserCredential) (conn *grpc.ClientConn, err error) {

	isDebug()

	host := helper.ConvertUnderlineToWhippletree(serviceName) + helper.GetEnv("SSLSuffixServerName")
	address := "127.0.0.1"
	return setClientConn(host, address, userCredential)
}

func setClientConn(host string, address string, userCredential []*UserCredential) (conn *grpc.ClientConn, err error) {
	var opts []grpc.DialOption
	var optsCallOption []grpc.CallOption
	var creds credentials.TransportCredentials

	// 设置接收最大条数
	optsCallOption = append(optsCallOption, grpc.MaxCallRecvMsgSize(100*1024*1024))
	opts = append(opts, grpc.WithDefaultCallOptions(optsCallOption...))

	// client to server
	if strings.ToLower(helper.GetEnv("SSL")) == "true" {
		// k8s-k8s
		certFile := helper.GetEnv("SSLCACertFile")
		creds, err = credentials.NewClientTLSFromFile(path+"/"+certFile, host)
		if err != nil {
			panic(err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		// 本机-k8s
		opts = append(opts, grpc.WithInsecure())
	}
	// 本机-局域网
	// 本机-本机

	// 用户信息
	if len(userCredential) == 1 && userCredential[0] != nil {
		opts = append(opts, grpc.WithPerRPCCredentials(userCredential[0]))
	}

	log.Println(address+runPort, host)
	grpclog.Info("Certificate Host: ", host)
	grpclog.Info("Connect Server: ", address+runPort)
	//conn, err = grpc.Dial(address+port, opts...)
	conn, err = grpc.Dial(address+runPort, opts...)
	if err != nil {
		grpclog.Error(err)
	}

	return
}

// 拦截器
func grpcInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo) {
	// 不启动阿里云
	if alilog.LogOff || alilog.LogStore == nil {
		return
	}
	// 阿里云日志内容
	var contents []*sls.LogContent
	// metadata
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// authority
		if data, ok := md[":authority"]; ok {
			key := "authority"
			contents = append(contents, &sls.LogContent{Key: &key, Value: &data[0]})
		}
		// user_id
		if data, ok := md["user_id"]; ok {
			key := "user_id"
			contents = append(contents, &sls.LogContent{Key: &key, Value: &data[0]})
		}
		// username
		if data, ok := md["username"]; ok {
			key := "username"
			contents = append(contents, &sls.LogContent{Key: &key, Value: &data[0]})
		}
	}
	// grpc method
	methodKey := "grpc_method"
	methodValue := info.FullMethod
	contents = append(contents, &sls.LogContent{Key: &methodKey, Value: &methodValue})
	// grpc param
	paramKey := "grpc_param"
	paramValue := fmt.Sprint(req)
	contents = append(contents, &sls.LogContent{Key: &paramKey, Value: &paramValue})
	// 发送日志
	pushAliLog(contents)
}

// userCredential 用户认证
type UserCredential struct {
	User map[string]string
}

// 用户凭证
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
	return strings.ToLower(helper.GetEnv("SSL")) == "true"
}

// 用户凭证
func NewUserCredential() *UserCredential {
	return new(UserCredential)
}

func pushAliLog(contents []*sls.LogContent) {
	// 日志
	var slsLogs []*sls.Log
	// timeNowUnix
	timeNowUnix := uint32(time.Now().Unix())
	slsLogs = append(slsLogs, &sls.Log{
		Time:     &timeNowUnix,
		Contents: contents,
	})
	// 日志组
	logGroup := &sls.LogGroup{
		Topic:  &alilog.LogTopic,
		Source: &alilog.IpSource,
		Logs:   slsLogs,
	}
	// 发送日志
	err := alilog.LogStore.PutLogs(logGroup)
	if err != nil {
		logrus.Error("logStore.PutLogs error : " + err.Error())
	}
}

// Debug设置
// Debug=true 开启调试
// DebugServerListenPort=:50051 自定义端口
func isDebug() {
	if strings.ToLower(helper.GetEnv("Debug")) == "true" {
		if helper.GetEnv("DebugServerListenPort") != "" {
			debugPort = helper.GetEnv("DebugServerListenPort")
		}
	}
	if debugPort != "" {
		runPort = debugPort
	} else {
		runPort = port
	}
}
