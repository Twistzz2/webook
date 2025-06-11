package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Twistzz2/webook/webook-be/internal/handler"
	"github.com/Twistzz2/webook/webook-be/internal/handler/middleware"
	"github.com/Twistzz2/webook/webook-be/internal/repository"
	"github.com/Twistzz2/webook/webook-be/internal/repository/dao"
	"github.com/Twistzz2/webook/webook-be/internal/repository/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := InitDB()
	u := InitUser(db)

	router := InitRouter(u)

	router.Run(":8080")
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		// 只在初始化过程中使用 panic
		// 因为 panic 会使当前协程终止
		// 初始化过程有问题，这个协程就不应该启动
		panic("无法连接到数据库: " + err.Error())
	}

	err = dao.InitTable(db)
	if err != nil {
		panic("无法初始化数据库表: " + err.Error())
	}

	return db
}

func InitUser(db *gorm.DB) *handler.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := handler.NewUserHandler(svc)
	return u
}

func InitRouter(u *handler.UserHandler) *gin.Engine {
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

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.Use(middleware.NewSignInMiddlewareBuilder().Build())

	u.RegisterRoutes(router)

	return router
}
