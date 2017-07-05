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
	var opts []grpc.ServerOption
	if helper.GetEnv("SSL") == "true" {
		creds, err := credentials.NewServerTLSFromFile(path+"/config/server.crt", path+"/config/server.key")
		if err != nil {
			grpclog.Errorf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	server = grpc.NewServer(opts...)

	return server
}

func RunServer() {

	_, err := zookeeper.RegMicroServer()

	if err != nil {
		grpclog.Fatal("reg server fail, ", err)
	}

	listen_addr := helper.GetServerAddress()

	grpclog.Info("server address is ", listen_addr)

	listen, err := net.Listen("tcp", listen_addr)
	if err != nil {
		grpclog.Error(err)
	}

	err = server.Serve(listen)
	if err != nil {
		grpclog.Fatal(err)
	}
}

var prefixKey = "zk:"

func GetClientConn(service_name string) (client *grpc.ClientConn, err error) {

	db, _ := strconv.Atoi(helper.GetEnv("ZkRedisDb"))
	redisClient, err := redis.NewConnDB(db)
	if err != nil {
		return
	}
	defer redisClient.Close()
	service_name = prefixKey + helper.GetEnv("ZkRootPath") + "/" + service_name

	addr, err := redisClient.Get(service_name).Result()
	if err != nil {
		return
	}

	var opts []grpc.DialOption
	var creds credentials.TransportCredentials

	if helper.GetEnv("SSL") == "true" {
		creds, err = credentials.NewClientTLSFromFile("config/server.crt", "192.168.0.193")
		if err != nil {
			return
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithTimeout(5 * time.Second))
	client, err = grpc.Dial(addr, opts...)
	return
}

func GetLocalClientConn() (conn *grpc.ClientConn, err error) {

	address := helper.GetServerAddress()

	var opts []grpc.DialOption
	var creds credentials.TransportCredentials

	if helper.GetEnv("SSL") == "true" {
		creds, err = credentials.NewClientTLSFromFile("config/server.crt", "sude365.com")
		if err != nil {
			panic(err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err = grpc.Dial(address, opts...)
	if err != nil {
		grpclog.Error(err)
	}

	return
}