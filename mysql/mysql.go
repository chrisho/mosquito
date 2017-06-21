package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type Mysql struct {
	db *gorm.DB
}

var (
	mysql *Mysql
)

func NewConn(dataSource string) (*gorm.DB, error) {

	if mysql.db != nil {
		return mysql.db, nil
	}

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

func getConn() (*gorm.DB, error){
	if mysql.db == nil {
		return nil, error("not exist connection")
	}
	return mysql.db, nil
}