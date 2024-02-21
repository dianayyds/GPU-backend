package dao

import (
	"crypto/sha1"
	"encoding/hex"
	"gin_exercise/controller"
	"gin_exercise/mydb"
)

func Sha1password(str string) string {
	o := sha1.New()
	o.Write([]byte(str))
	return hex.EncodeToString(o.Sum(nil))
}

func Userinfobyname(name string) (*controller.User, error) {
	user2 := controller.User{}
	result := mydb.UserDB.Model(&controller.User{}).Where("username=?", name).Find(&user2)
	return &user2, result.Error
}

func Adduser(user *controller.User) error {
	result := mydb.UserDB.Create(user)
	if result == nil {
		return result.Error
	} else {
		return nil
	}
}

func Allusers() (*[]controller.User, error) {
	var users = make([]controller.User, 0)
	result := mydb.UserDB.Find(&users)
	return &users, result.Error
}