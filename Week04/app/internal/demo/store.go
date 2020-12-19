package demo

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DataBase struct {
	DB *sql.DB
}

func NewDB() *DataBase {
	database := &DataBase{
		nil,
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.addr"),
		viper.GetString("db.name"),
		viper.GetBool("db.debug"),
		//"Asia/Shanghai"),
		"Local")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		panic(err)
	}
	database.DB, err = db.DB()
	if err != nil {
		panic(err)
	}
	return database
}

func (d *DataBase) Stop() {
	if d != nil && d.DB != nil {
		_ = d.DB.Close()
	}
}
