package api

import (
	"gin_exercise/controller"
	"gin_exercise/dao"
	"gin_exercise/jwtauth"

	"github.com/gin-gonic/gin"
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

func IndexHandler(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}
