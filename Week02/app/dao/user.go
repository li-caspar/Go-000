package dao

import (
	"app/model"
	"github.com/pkg/errors"
)

type User struct {
	Id   int64
	Name string
}

func (User) TableName() string {
	return "user"
}

func (u *User) GetOne(id int) error {
	db, err := model.GetDB()
	if err != nil {
		return errors.Wrap(err, "get db error")
	}
	if err := db.First(u, id).Error; err != nil {
		return errors.Wrap(err, "db first user fail")
	}
	return nil
}
