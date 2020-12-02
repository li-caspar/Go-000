package service

import "app/dao"

func GetUser(id int) (dao.User, error) {
	user := dao.User{}
	err := user.GetOne(id)
	return user, err
}