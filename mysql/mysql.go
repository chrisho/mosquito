package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
	"errors"
)

var db *gorm.DB

func NewConn(dataSource, prefix string) (*gorm.DB, error) {

	if db != nil {
		return db, nil
	}

	setTablePrefix(prefix)
	db, err := gorm.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}

	db.DB().SetConnMaxLifetime(30 * time.Second)

	return db, err
}

func GetConn() (*gorm.DB, error){
	if db == nil {
		return nil, errors.New("connection is not exist")
	}
	return db, nil
}

func setTablePrefix(prefix string) {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return prefix + defaultTableName
	}
}