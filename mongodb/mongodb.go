package mongodb

import (
	"os"
	"strings"
	"gopkg.in/mgo.v2"
)

var (
	Conn *mgo.Session
	connStr string
	db      string
)

type Person struct {
	Name  string
	Phone string
}

// create mongodb connection
func NewConn() (conn *mgo.Session, err error) {

	if connStr == "" {
		initMongodbParams()
	}

	conn, err = mgo.Dial(connStr)

	return
}

// init redis params
func initMongodbParams() {
	db = os.Getenv("MongodbDb")

	hosts := os.Getenv("MongodbHost")
	user := os.Getenv("MongodbUser")
	password := os.Getenv("MongodbPassword")
	options := os.Getenv("MongodbOptions")

	hostList := strings.Split(hosts, ",")

	hostLen, firstHost, anotherHost := len(hostList), "", ""

	if hostLen < 1 {
		println("mongodb host is empty")
	} else {
		firstHost = hostList[0]
	}

	if hostLen > 1 {
		for i := 1; i < hostLen; i++ {
			anotherHost += "," + hostList[i]
		}
	}

	if user == "" {
		connStr = firstHost
	} else {
		connStr = user + ":" + password + "@" + firstHost
	}

	if anotherHost != "" {
		connStr += anotherHost
	}

	if db != "" {
		connStr += "/" + db
	}

	if options != "" {
		connStr += "?" + options
	}
}
