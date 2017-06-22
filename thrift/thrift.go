package thrift

import (
	"errors"
	"crypto/tls"
	"git.apache.org/thrift.git/lib/go/thrift"
)

// processor ：thrift generate interfaces -> server implements interfaces -> thrift server new process
// secure : use or no SSL
// certFile : ssl cert
// keyFile : ssl key
// --------------------------------------------------
// example : get a SdUserServiceProcessor
// --------------------------------------------------
// handler := &UserServiceImpl{}
// processor := userThrift.NewSdUserServiceProcessor(handler)
// --------------------------------------------------
func RunServer(processor thrift.TProcessor, serverAddress string, secure bool, certFile, keyFile string) error {

	if secure && (certFile == "" || keyFile == "") {
		return errors.New("certFile or keyFile is empty")
	}

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	var transport thrift.TServerTransport
	var err error
	if secure {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair("server.crt", "server.key"); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return err
		}
		transport, err = thrift.NewTSSLServerSocket(serverAddress, cfg)
	} else {
		transport, err = thrift.NewTServerSocket(serverAddress)
	}

	if err != nil {
		return err
	}

	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	println("Starting the server on ", serverAddress)

	return server.Serve()
}

// secure : use or no SSL
// --------------------------------------------------
// example : get a SdUserServiceClient
// --------------------------------------------------
// client := userThrift.NewSdUserServiceClientFactory(transport, protocolFactory)
// --------------------------------------------------
// you should close transport : defer transport.Close()
// --------------------------------------------------
func RunClient(serverAddress string, secure bool) (transport thrift.TTransport, protocolFactory *thrift.TBinaryProtocolFactory, err error) {

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	if secure {
		cfg := new(tls.Config)
		cfg.InsecureSkipVerify = true
		transport, err = thrift.NewTSSLSocket(serverAddress, cfg)
	} else {
		transport, err = thrift.NewTSocket(serverAddress)
	}

	if err != nil {
		return
	}

	transport = transportFactory.GetTransport(transport)

	err = transport.Open()

	return
}