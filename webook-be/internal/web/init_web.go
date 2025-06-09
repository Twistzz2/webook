package web

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.Use(func(ctx *gin.Context) {
		fmt.Println("这是第一个 middleware")
	})

	router.Use(func(ctx *gin.Context) {
		fmt.Println("这是第二个 middleware")
	})

	router.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"http://localhost:3000"},
		// AllowMethods 不传时，默认允许所有方法
		// AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		// 业务请求中允许的请求头
		AllowHeaders:  []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
		// 是否允许携带 Cookie
		AllowCredentials: true,
		// 哪些来源是允许的
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourdomain.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	u := NewUserHandler()
	u.RegisterRoutes(router)

	return router
}
