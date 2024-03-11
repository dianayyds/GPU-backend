package httpserver

import (
	"gin_exercise/api"
	"gin_exercise/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Router struct {
}

func Initroute() {
	//html页面位置
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	//静态文件位置
	r.Static("/static", "./static")
	router := Router{}
	r.Use(router.corsMiddleware)
	r.NoRoute(router.notFound)
	r.NoMethod(router.noMethod)
	//限制上传最大尺寸
	r.MaxMultipartMemory = 8 << 20                               // 8 MiB
	r.Use(middleware.RateLimitMiddleware(time.Second, 100, 100)) //初始100，每秒放出100
	v1 := r.Group("/v1")
	{
		v1.POST("/sign_up", api.UsersignupHandler)
		v1.POST("/sign_in", api.UsersigninHandler)
		v1.POST("/parse_jwt", api.ParseJwtHandler)
		// v1.POST("/init_database", api.InitdatabaseHandler)
		v1.POST("/delete_user", api.DeleteUserHandler)
		v1.POST("/gpt", api.GptHandler)
		v1.POST("/ssh_connect", api.SshConnectHandler)
		v1.GET("/cpu_info", api.CpuinfoHandler)
		v1.GET("/gpu_info", api.GpuinfoHandler)
		v1.GET("/detailed_gpu_info", api.DetailedGPUInfoHandler)
		v1.GET("/users_info", api.UsersInfoHandler)
		v1.GET("/user_info_byname", api.UserinfobynameHandler)
		v1.GET("/base_info", api.BaseinfoHandler)
	}
	r.Run()
}

func (r Router) notFound(q *gin.Context) {
	q.JSON(http.StatusNotFound, gin.H{
		"code": 1,
		"msg":  "抱歉，没有相对应的路径",
	})
}

func (r Router) noMethod(q *gin.Context) {
	q.JSON(http.StatusMethodNotAllowed, gin.H{
		"code": 1,
		"msg":  "抱歉，没有相对应的方法",
	})
}

func (r Router) corsMiddleware(q *gin.Context) {
	q.Writer.Header().Set("Access-Control-Allow-Origin", q.GetHeader("Origin"))
	q.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	q.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Token")
	q.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if q.Request.Method == "OPTIONS" {
		q.AbortWithStatus(200)
		return
	}
	q.Next()
}
