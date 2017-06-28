package mosquito

import (
	"github.com/chrisho/mosquito/helper"
	"github.com/chrisho/mosquito/zookeeper"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"net"
	"os"
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
	listen_addr := helper.GetServerAddress()

	grpclog.Info("server address is ", listen_addr)

	_, err := zookeeper.RegMicroServer()

	if err != nil {
		grpclog.Fatal("reg server fail, ", err)
	}

	listen, err := net.Listen("tcp", listen_addr)
	if err != nil {
		grpclog.Error(err)
	}

	err = server.Serve(listen)
	if err != nil {
		grpclog.Fatal(err)
	}
}

func GetClientConn() (conn *grpc.ClientConn, err error) {
	address := helper.GetServerAddress()
	var opts []grpc.DialOption
	var creds credentials.TransportCredentials

	if helper.GetEnv("SSL") == "true" {
		creds, err = credentials.NewClientTLSFromFile("config/server.crt", "sude365.com")
		if err != nil {
			panic(err)
		}
	}
	opts = append(opts, grpc.WithTransportCredentials(creds), grpc.WithInsecure())

	conn, err = grpc.Dial(address, opts...)
	if err != nil {
		grpclog.Error(err)
	}

	return
}
