package mosquito

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"crypto/tls"
	"github.com/chrisho/mosquito/helper"
)


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