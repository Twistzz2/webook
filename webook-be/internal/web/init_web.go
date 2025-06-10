package web

import (
	"fmt"
	"strings"
	"time"

	"github.com/Twistzz2/webook/webook-be/internal/repository"
	"github.com/Twistzz2/webook/webook-be/internal/repository/dao"
	"github.com/Twistzz2/webook/webook-be/internal/repository/service"
	"github.com/Twistzz2/webook/webook-be/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitRouter() *gin.Engine {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		// 只在初始化过程中使用 panic
		// 因为 panic 会使当前协程终止
		// 初始化过程有问题，这个协程就不应该启动
		panic("无法连接到数据库: " + err.Error())
	}

	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)

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

	u.RegisterRoutes(router)

	return router
}
