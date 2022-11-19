package main

import (
	"OSS/apiServer/bucket"
	"OSS/apiServer/heartbeat"
	"OSS/apiServer/locate"
	"OSS/apiServer/objects"
	"OSS/apiServer/system"
	"OSS/apiServer/temp"
	"OSS/apiServer/versions"
	RedisMQ "OSS/lib/Redis"
	"OSS/tools"
	utils2 "OSS/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()
		c.Next()
	}
}
func InitRouter(r *gin.Engine) {

	r.Use(Cors())

	//objects
	{
		r.PUT("/objects/*id", objects.Put)
		r.GET("/objects/*id", objects.Get)
		r.DELETE("/objects/*id", objects.Del)
		r.POST("/objects/*id", objects.Post)
	}
	//locate
	{
		r.GET("/locate/*id", locate.Get)
	}
	//versions
	{
		r.GET("/versions/*id", versions.Get)
	}
	//temp
	{
		r.HEAD("/temp/*id", temp.Head)
		r.PUT("/temp/*id", temp.Put)
	}
	//headrtbeat
	{
		r.GET("/heartbeat", heartbeat.Get)
	}
	//system
	{
		r.GET("/nodeSystemInfo/*id", system.Get)
	}
	//bucket
	{
		r.GET("/bucket/*id", bucket.Get)
		r.PUT("/bucket/*id", bucket.Put)
		r.DELETE("/bucket/*id", bucket.Del)
		r.POST("/bucket/*id", bucket.Head)
	}
	//allVersions
	{
		r.GET("/allVersions/*id", versions.AllGet)
	}
	//utils
	{
		r.GET("/deleteOldMetadata/*id", tools.DelOldMetaDate)
		r.GET("/deleteOrphanServer/*id", tools.DelOrphan)
		r.GET("/objectScanner/*id", tools.ObjectScanner)
	}
}
func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)

	utils2.Rds = RedisMQ.NewRedis(os.Getenv("REDIS_SERVER"))
	defer utils2.Rds.Client.Close()

	//开始连接apiServers这个exchanges，将数据服务节点的地址保存起来
	go heartbeat.ListenHeartbeat()

	r := gin.Default()
	InitRouter(r)

	fmt.Println(os.Getenv("LISTEN_ADDRESS"), "===>apiServer Start running <===")

	//监听并启动 ip在tools中规划好了
	//目前是10.29.2.1和10.29.2.2
	log.Fatal(r.Run(os.Getenv("LISTEN_ADDRESS")))
}
