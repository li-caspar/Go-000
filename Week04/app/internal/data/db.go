package data

import (
	"app/internal/config"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.DB.UserName,
		cfg.DB.Password,
		cfg.DB.Addr,
		cfg.DB.Name,
		cfg.DB.Charset,
		true,
		"Local",
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		return db, errors.Wrap(err, "can not open db connection")
	}
	return db, nil
}
