package middleware

import (
	"fmt"
	"gin_exercise/dao"
	"gin_exercise/jwtauth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("authorization")
		if token == "" {
			c.Next()
		}
		// 按空格分割
		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 1,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}

		claim, err := jwtauth.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 1,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		name := claim.Name
		_, err = dao.Userinfobyname(name)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2006,
				"msg":  fmt.Sprintf("无效的姓名:%v", err),
			})
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("name", name)
		c.Next() // 后续的处理函数可以用过c.Get("name")来获取当前请求的用户信息
	}
}
