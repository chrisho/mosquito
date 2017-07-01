package mongodb

import (
	"os"
	"strings"
	"gopkg.in/mgo.v2"
)

var (
	connStr string
	db      string
)

type Person struct {
	Name  string
	Phone string
}

func GetMongodbConn() (conn *mgo.Session) {
	return
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

	// set default redis host
	if hostLen < 1 {
		println("mongodb host is empty, use default address 127.0.0.1:6379 ")
		firstHost = "127.0.0.1:27017"
	} else {
		firstHost = hostList[0]
	}

	// has many redis address ?
	if hostLen > 1 {
		for i := 1; i < hostLen; i++ {
			anotherHost += "," + hostList[i]
		}
	}

	// join connStr
	if user == "" {
		connStr = firstHost
	} else {
		connStr = user + ":" + password + "@" + firstHost
	}

	// join connStr
	if anotherHost != "" {
		connStr += anotherHost
	}

	// join connStr and set use db
	if db != "" {
		connStr += "/" + db
	}

	// join connStr and set use options
	if options != "" {
		connStr += "?" + options
	}
}
