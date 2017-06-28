package mosquito

import (
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
	creds, err := credentials.NewServerTLSFromFile(path+"/config/server.crt", path+"/config/server.key")
	if err != nil {
		grpclog.Errorf("Failed to generate credentials %v", err)
	}
	opts = []grpc.ServerOption{grpc.Creds(creds)}

	server = grpc.NewServer(opts...)

	return server
}

func RunServer() {
	listen_addr := zookeeper.GetServerAddress()
	grpclog.Info("server address is ", listen_addr)

	children, err := zookeeper.RegMicroServer()

	if err != nil {
		grpclog.Fatal(err)
	}
	grpclog.Info("zookeeper servers is ", children)

	listen, err := net.Listen("tcp", listen_addr)

	if err != nil {
		grpclog.Error(err)
	}

	err = server.Serve(listen)

	if err != nil {
		grpclog.Fatal(err)
	}
}

func GetClientConn() *grpc.ClientConn {
	address := zookeeper.GetServerAddress()
	grpclog.Info("serverAddress is ", address)
	var opts []grpc.DialOption
	var creds credentials.TransportCredentials

	creds, err := credentials.NewClientTLSFromFile("config/server.crt", "sude365.com")

	opts = append(opts, grpc.WithTransportCredentials(creds), grpc.WithInsecure())

	conn, err := grpc.Dial(address, opts...)

	if err != nil {
		grpclog.Error(err)
	}

	return conn
}
