package api

import (
	"fmt"
	"gin_exercise/controller"
	"gin_exercise/dao"
	"gin_exercise/jwtauth"
	"gin_exercise/mydb"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func UsersignupHandler(g *gin.Context) {
	user1 := controller.User{}
	g.Bind(&user1)
	user2, _ := dao.Userinfobyname(user1.Username)
	if user2.Username == user1.Username {
		g.JSON(200, gin.H{
			"code":     1,
			"msg":      "aready exist",
			"username": user1.Username,
		})
		return
	}
	user1.Password = dao.Sha1password(user1.Password)
	err := dao.Adduser(&user1)
	if err != nil {
		g.JSON(200, gin.H{
			"code": 2,
			"msg":  "add user failed",
			"用户名":  user1.Username,
		})
	} else {
		g.JSON(200, gin.H{
			"code": 0,
			"msg":  "add user success",
			"用户名":  user1.Username,
		})
	}
}

func UsersigninHandler(g *gin.Context) {
	// 用户发送用户名和密码过来
	// 校验用户名和密码是否正确
	user1 := controller.User{}
	g.Bind(&user1)
	user2, _ := dao.Userinfobyname(user1.Username)
	if user2 == nil {
		g.JSON(200, gin.H{
			"code":     1,
			"msg":      "name not found",
			"username": user1.Username,
		})
	}
	pass := dao.Sha1password(user1.Password)
	if user2.Password != pass {
		g.JSON(200, gin.H{
			"code":     1,
			"msg":      "password wrong",
			"username": user1.Username,
		})
		return
	}
	tokenString, err := jwtauth.GenToken(user1.Username)
	if err != nil {
		g.JSON(200, gin.H{
			"code":     1,
			"msg":      "token wrong",
			"username": user1.Username,
		})
	}
	g.Set("name", user1.Username)
	g.JSON(200, gin.H{
		"code":  0,
		"msg":   "signin success",
		"token": tokenString,
	})
}

func UserinfobynameHandler(g *gin.Context) {
	username := g.Query("username")
	user, _ := dao.Userinfobyname(username)
	g.JSON(200, gin.H{
		"code": 0,
		"msg":  "find user success",
		"user": user,
	})

}

func ParseJwtHandler(g *gin.Context) {
	token := controller.Mytoken{}
	g.Bind(&token)
	claims, err := jwtauth.ParseToken(token.Token)
	if err != nil {
		g.JSON(200, gin.H{
			"code":  1,
			"error": err,
		})
	} else {
		g.JSON(200, gin.H{
			"code":   0,
			"claims": claims,
		})
	}

}

func IndexHandler(g *gin.Context) {
	g.HTML(200, "index.html", nil)
}

func InitdatabaseHandler(g *gin.Context) {
	database := controller.User{}
	g.Bind(&database)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		database.DatabaseUsername, database.DatabasePassword, database.Ip, database.Port, database.DatabaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	mydb.UserDB.Model(&controller.User{}).Where("username=?", database.Username).Updates(database)
	if err != nil {
		seelog.Error(err)
		g.JSON(200, gin.H{
			"code": 1,
			"msg":  err,
		})
	} else {
		mydb.InfoDB = db
		g.JSON(200, gin.H{
			"code": 0,
		})
	}
}

func UsersInfoHandler(g *gin.Context) {
	users, err := dao.Allusers()
	if err != nil {
		g.JSON(200, gin.H{
			"code": 1,
			"msg":  err,
		})
	} else {
		g.JSON(200, gin.H{
			"code":  0,
			"users": users,
		})
	}
}

func DeleteUserHandler(g *gin.Context) {
	user := controller.User{}
	g.Bind(&user)
	err := dao.Deleteuser(user.Username)
	if err != nil {
		g.JSON(200, gin.H{
			"code": 1,
			"msg":  err,
		})
	} else {
		g.JSON(200, gin.H{
			"code": 0,
			"msg":  "delete success",
		})
	}
}
