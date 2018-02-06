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

var (
	db          *gorm.DB
	TablePrefix = helper.GetEnv("MysqlPrefix")
)

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
		return nil, errors.New("connection is not exist")
	}

	db = connection

	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	setMysqlDebug(db)
	db.DB().SetConnMaxLifetime(30 * time.Second)

	return db, err
}

func GetConn() (*gorm.DB, error) {
	db, err := newConn()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func setMysqlDebug(db *gorm.DB) {
	if helper.GetEnv("MysqlDebug") == "TRUE" {
		db.LogMode(true)
	}
}

func setTablePrefix() {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		if !strings.HasPrefix(defaultTableName, TablePrefix) {
			return TablePrefix + defaultTableName
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
