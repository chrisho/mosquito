package zookeeper

import (
	"github.com/chrisho/mosquito/helper"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
	"os"
	"log"
	"time"
)

var (
	zkConn       *zk.Conn
	zkHost       []string
	acl          []zk.ACL
	zkTimeout    = 3 * time.Second
	zkRootPath   string
	zkServerPath string
)

func GetZkConn() (conn *zk.Conn) {
	if zkConn != nil {
		return zkConn
	}
	zkConn, _, _ = NewConn()
	return zkConn
}

// create zookeeper connection
func NewConn() (conn *zk.Conn, event <-chan zk.Event, err error) {

	if len(zkHost) < 1 {
		initZookeeperParams()
	}

	conn, event, err = zk.Connect(zkHost, zkTimeout)

	zkConn = conn
	addAuth()

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

// addAuth digest credential
func addAuth() {
	digest := helper.GetEnv("ZkAuthScheme")
	credential := helper.GetEnv("ZkAuthCredential")
	err := zkConn.AddAuth(digest, []byte(credential))

	if err != nil {
		log.Println(err)
	}
}

// get acl : zk will be set this acl
func getAcl() []zk.ACL {

	if len(acl) > 0 {
		return acl
	}

	credential := strings.Split(helper.GetEnv("ZkAuthCredential"), ":")

	if len(credential) <= 1 {
		acl = zk.AuthACL(zk.PermAll)
	} else {
		acl = zk.DigestACL(zk.PermAll, credential[0], credential[1])
	}

	return acl
}
