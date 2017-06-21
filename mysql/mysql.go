package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
	"errors"
)

type Mysql struct {
	db *gorm.DB
}

var (
	mysql *Mysql
)

func NewConn(dataSource, prefix string) (*gorm.DB, error) {

	if mysql.db != nil {
		return mysql.db, nil
	}

	setTablePrefix(prefix)
	db, err := gorm.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}

	db.DB().SetConnMaxLifetime(30 * time.Second)
	mysql = &Mysql{
		db,
	}

	return mysql.db, err
}

func GetConn() (*gorm.DB, error){
	if mysql.db == nil {
		return nil, errors.New("connection is not exist")
	}
	return mysql.db, nil
}

func setTablePrefix(prefix string) {
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return prefix + defaultTableName
	}
}