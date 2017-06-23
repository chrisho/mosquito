package zookeeper

import (
	"os"
	"strings"
	"time"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	zkHost       []string
	ZkTimeout    = 3 * time.Second
	ZkRootPath   string
	ZkServerPath string
)

func GetZkConn() (conn *zk.Conn) {
	return
}

// create zookeeper connection
func NewConn() (conn *zk.Conn, event <-chan zk.Event, err error) {

	if len(zkHost) < 1 {
		initZookeeperParams()
	}

	conn, event, err = zk.Connect(zkHost, ZkTimeout)

	return
}

// init zookeeper params
func initZookeeperParams() {
	zkHost = strings.Split(os.Getenv("ZookeeperHost"), ",")

	ZkRootPath = os.Getenv("ZkRootPath")

	ZkServerPath = ZkRootPath + os.Getenv("ServerPath")

	if ZkServerPath == "" {
		println("zookeeper server path is empty")
	}
}
