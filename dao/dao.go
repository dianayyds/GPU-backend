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
	result := mydb.DB.Model(&controller.User{}).Where("username=?", name).Find(&user2)
	return &user2, result.Error
}

func Adduser(user *controller.User) error {
	result := mydb.DB.Create(user)
	if result == nil {
		return result.Error
	} else {
		return nil
	}
}

func GetUserList(pagenum int, pagesize int) (*[]controller.User, int64, error) {
	var users = make([]controller.User, 0)
	query := mydb.DB.Model(&controller.User{})
	var total int64
	query.Count(&total)
	result := query.Find(&users)
	//result := query.Limit(pagesize).Offset((pagenum - 1) * pagesize).Find(&users)

	return &users, total, result.Error
}
