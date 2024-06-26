package dao

import (
	"crypto/sha1"
	"encoding/hex"
	"gin_exercise/controller"
	"gin_exercise/mydb"

	"golang.org/x/crypto/ssh"
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

func Adduserssh(user *controller.Userssh) error {
	userssh:=controller.Userssh{}
	result := mydb.UserDB.Model(&controller.Userssh{}).Where("username=?", user.Username).Find(&userssh)
	if result.RowsAffected > 0{
		result = mydb.UserDB.Model(&controller.Userssh{}).Where("username=?", user.Username).Updates(user)
		return result.Error
	}else{
		result = mydb.UserDB.Model(&controller.Userssh{}).Create(user)
		return result.Error
	}
}

func Allusers() (*[]controller.User, error) {
	var users = make([]controller.User, 0)
	result := mydb.UserDB.Find(&users)
	return &users, result.Error
}

func AllsshInfo() (*[]controller.Userssh, error) {
	var users = make([]controller.Userssh, 0)
	result := mydb.UserDB.Find(&users)
	return &users, result.Error
}


func Deleteuser(name string) error {
	result := mydb.UserDB.Where("username=?", name).Delete(&controller.User{})
	return result.Error
}

func RunCommand(client *ssh.Client, cmd string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", err
	}
	return string(output), nil
}
