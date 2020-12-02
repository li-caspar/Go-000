package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func GetDB() (*gorm.DB, error) {
	if DB != nil {
		return DB, nil
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
		"root",
		"root",
		"127.0.0.1",
		"test",
		true,
		//"Asia/Shanghai"),
		"Local")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Discard})

	return db, err
}
