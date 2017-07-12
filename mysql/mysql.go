package mysql

import (
	"errors"
	"github.com/chrisho/mosquito/helper"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"
	"strings"
)

var db *gorm.DB

func init() {
	var err error
	db, err = newConn()
	if err != nil {
		log.Println(err)
	}
}

func newConn() (*gorm.DB, error) {

	if db != nil {
		return db, nil
	}

	setTablePrefix()
	connection, err := gorm.Open("mysql", getDataSource())
	if err != nil {
		return nil, err
	}

	db = connection

	db.DB().SetConnMaxLifetime(30 * time.Second)

	return db, err
}

func GetConn() (*gorm.DB, error) {
	if db == nil {
		return nil, errors.New("connection is not exist")
	}
	return db, nil
}

func setTablePrefix() {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		if !strings.HasPrefix(defaultTableName, helper.GetEnv("MysqlPrefix")) {
			return helper.GetEnv("MysqlPrefix") + defaultTableName
		}
		return defaultTableName
	}
}

func getDataSource() string {
	dataSourceName := ""
	dataSourceName += helper.GetEnv("MysqlUser") + ":"
	dataSourceName += helper.GetEnv("MysqlPassword") + "@"
	dataSourceName += "tcp(" + helper.GetEnv("MysqlHost") + ")/"
	dataSourceName += helper.GetEnv("MysqlName") + "?"
	dataSourceName += helper.GetEnv("MysqlParameters")

	return dataSourceName
}
