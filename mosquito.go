package mosquito

import (
	"net"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc"
	"github.com/chrisho/mosquito/helper"
	"os"
	"github.com/joho/godotenv"
	"google.golang.org/grpc/credentials"
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

func GetServer() (*grpc.Server){
	var opts []grpc.ServerOption
	creds, err := credentials.NewServerTLSFromFile(path + "/config/server.crt", path + "/config/server.key")
	if err != nil {
		grpclog.Errorf("Failed to generate credentials %v", err)
	}
	opts = []grpc.ServerOption{grpc.Creds(creds)}

	server = grpc.NewServer(opts...)

	return server
}

func RunServer() {
	listen_addr := ":" + helper.GetEnv("ServerPort")
	listen, err := net.Listen("tcp", listen_addr)
	if err != nil {
		grpclog.Error(err)
	}
	server.Serve(listen)
}

func GetClientConn() (*grpc.ClientConn) {
	address := helper.GetEnv("ServerAddress") + ":" + helper.GetEnv("ServerPort")
	var opts []grpc.DialOption
	var creds credentials.TransportCredentials


	creds, err := credentials.NewClientTLSFromFile("config/server.crt",
		"sude365.com")

	opts = append(opts, grpc.WithTransportCredentials(creds),
		grpc.WithInsecure())

	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		grpclog.Error(err)
	}

	return conn
}


/*
func RunServer(serverAddr, serverPort string, processor thrift.TProcessor) (err error){
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	address := serverAddr + ":" + serverPort
	ssl := helper.GetEnv("SSL")

	var transport thrift.TServerTransport
	if ssl == "true" {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair("config/server.crt", "config/server.key"); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return err
		}
		transport, err = thrift.NewTSSLServerSocket(address, cfg)
	} else {
		transport, err = thrift.NewTServerSocket(address)
	}

	if err != nil {
		return err
	}
	//handler := &serviceHandler.UserHandler{}
	//processor := thriftHandler.NewSdUserServiceProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	println("Starting the simple server... on ", address)
	return server.Serve()
}
*/