package zookeeper

import (
	"os"
	"log"
	"time"
	"strings"
	"errors"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/asaskevich/govalidator"
	"github.com/chrisho/mosquito/helper"
)

var (
	zkConn       *zk.Conn
	zkHost       []string
	zkTimeout    = 3 * time.Second
	zkRootPath   string
	zkServerPath string
)

func GetZkConn() (conn *zk.Conn) {
	return
}

// create zookeeper connection
func NewConn() (conn *zk.Conn, event <-chan zk.Event, err error) {

	if len(zkHost) < 1 {
		initZookeeperParams()
	}

	conn, event, err = zk.Connect(zkHost, zkTimeout)

	return
}

// init zookeeper params
func initZookeeperParams() {
	zkHost = strings.Split(os.Getenv("ZookeeperHost"), ",")

	if len(zkHost) == 0 {
		println("zookeeper server path is empty")
	}

	zkRootPath = helper.GetEnv("ZkRootPath")
	zkServerPath = zkRootPath + helper.GetEnv("ServerPath")
}

// register server into zookeeper
func RegMicroServer() (children []string, err error) {
	zkConn, _, _ = NewConn()
	addAuth()
	createRootPath()
	createServerPath()
	createServerAddressPath()

	children, _, err = zkConn.Children(zkServerPath)
	return
}

// addAuth digest credential
func addAuth() {
	digest := helper.GetEnv("ZkAuthScheme")
	credential := helper.GetEnv("ZkAuthCredential")
	err := zkConn.AddAuth(digest, []byte(credential))

	if err != nil {
		log.Println(err)
	}
}

// exist rootPath or create rootPath
func createRootPath() {

	if zkRootPath == "" {
		return
	}

	if ok, _, _ := zkConn.Exists(zkRootPath); !ok {
		_, err := zkConn.Create(zkRootPath, nil, 0, zk.AuthACL(zk.PermAll))

		if err != nil {
			log.Println(err)
		}
	}
}

// exist serverPath or create serverPath
func createServerPath() {

	if ok, _, _ := zkConn.Exists(zkServerPath); !ok {
		_, err := zkConn.Create(zkServerPath, nil, 0, zk.AuthACL(zk.PermAll))

		if err != nil {
			log.Println(err)
		}
	}
}

// create Ephemeral server path : address:port
func createServerAddressPath() {

	serverAddressPath := zkServerPath + "/" + getServerAddress()

	if ok, _, _ := zkConn.Exists(serverAddressPath); !ok {
		_, err := zkConn.Create(serverAddressPath, nil, zk.FlagEphemeral, zk.AuthACL(zk.PermRead))

		if err != nil {
			log.Println(err)
		}
	}
}

func getServerAddress() (ipAddress string) {
	serverAddress := helper.GetEnv("ServerAddress")
	serverPort := helper.GetEnv("ServerPort")

	if serverAddress == "" {
		log.Println(errors.New("serverAddress is empty or not exist"))
		serverAddress = "127.0.0.1"
	}

	if ok := govalidator.IsIP(serverAddress); !ok {

		ipCutSet := strings.TrimRight(serverAddress, ".*")

		serverAddress = helper.ContainsIp(ipCutSet)
	}

	if serverPort == "" {
		ipAddress = serverAddress
	} else {
		ipAddress = serverAddress + ":" + serverPort
	}

	return
}
