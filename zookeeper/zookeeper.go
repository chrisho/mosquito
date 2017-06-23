package zookeeper

import (
	"os"
	"strings"
	"time"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	zkHost       []string
	zkTimeout    = 3 * time.Second
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
}
