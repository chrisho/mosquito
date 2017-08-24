package zookeeper

import (
	"github.com/chrisho/mosquito/helper"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"strconv"
)

// register server into zookeeper
func RegMicroServer() (children []string, err error) {

	if zkConn == nil {
		zkConn, _, _ = NewConn()
	}

	createRootPath()
	createServerPath()
	createServerAddressPath()

	children, _, err = zkConn.Children(zkServerPath)
	return
}

// exist rootPath or create rootPath
func createRootPath() {

	if zkRootPath == "" {
		return
	}

	if ok, _, _ := zkConn.Exists(zkRootPath); !ok {
		_, err := zkConn.Create(zkRootPath, nil, 0, getAcl())

		if err != nil {
			log.Println(err)
		}
	}
}

// exist serverPath or create serverPath
func createServerPath() {

	if zkServerPath == "" {
		return
	}

	if ok, _, _ := zkConn.Exists(zkServerPath); !ok {
		_, err := zkConn.Create(zkServerPath, nil, 0, getAcl())

		if err != nil {
			log.Println(err)
		}
	}
}

// create Ephemeral server path : address:port
func createServerAddressPath() {

	index, address := checkServerAddress()

	serverAddressPath := zkServerPath + "/" + index

	if ok, _, _ := zkConn.Exists(serverAddressPath); !ok {
		_, err := zkConn.Create(serverAddressPath, []byte(address), zk.FlagEphemeral, getAcl())

		if err != nil {
			log.Println(err)
		}
	}
}

func checkServerAddress() (index, address string) {

	index = "1"
	address = helper.GetEnv("ServerRegisterIP")
	if address == "" {
		address = helper.GetServerAddress()
	} else {
		address += ":" + helper.GetEnv("ServerPort")
	}

	// micro server path is set ?
	if zkServerPath == "" {
		zkServerPath = "/"
	}

	children, _, err := zkConn.Children(zkServerPath)

	if err != nil {
		log.Println(err)
	}

	childLen := len(children)

	if childLen == 0 {
		return
	} else if childLen == 1 {
		if children[0] == "zookeeper" {
			return
		}
	}

	index = existServerAddress(children, address)

	return
}

func existServerAddress(children []string, address string) (index string) {

	var maxIndex int = 0

	for _, path := range children {
		if path == "zookeeper" {
			continue
		}

		index, _ := strconv.Atoi(path)

		if maxIndex < index {
			maxIndex = index
		}
	}

	maxIndex += 1

	index = strconv.Itoa(maxIndex)

	return
}
