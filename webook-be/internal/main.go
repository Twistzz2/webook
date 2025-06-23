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
	"github.com/gin-contrib/sessions/redis"
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
	dsn := "root:root@tcp(localhost:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// 只在初始化过程中使用 panic
		// 因为 panic 会使当前协程终止
		// 初始化过程有问题，这个协程就不应该启动
		fmt.Printf("数据库连接失败: %v\n", err)
		fmt.Printf("请确保MySQL服务在端口13316上运行\n")
		panic("无法连接到数据库: " + err.Error())
	}

	err = dao.InitTable(db)
	if err != nil {
		fmt.Printf("数据库表初始化失败: %v\n", err)
		panic("无法初始化数据库表: " + err.Error())
	}

	fmt.Println("数据库连接成功")
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

	// 跨域问题的解决
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

	// 使用 cookie 存储 session
	// store := cookie.NewStore([]byte("secret"))

	// 基于内存的实现，第一个参数是 32 位或 64 位的 authenticaiton key，第二个参数是 encryption key
	// store := memstore.NewStore([]byte("authenticaitonKey"), []byte("encryptionKey"))
	// store := memstore.NewStore([]byte("7F3C9A1E5D6B2F4A9831B6D472C0E9AD"), []byte("A4B29D5E7C183F62E5DA07B1C3498FE0"))	// func NewStore(size int, network, address, username, password string, keyPairs ...[]byte) (Store, error) {
	// 恢复使用两个密钥，确保加密解密正确
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "", "",
		[]byte("7F3C9A1E5D6B2F4A9831B6D472C0E9AD"),
		[]byte("A4B29D5E7C183F62E5DA07B1C3498FE0"))
	if err != nil {
		fmt.Printf("Redis连接失败: %v\n", err)
		fmt.Printf("请确保Redis服务在端口6379上运行\n")
		panic("无法连接到 Redis: " + err.Error())
	}

	fmt.Println("Redis连接成功")

	// 设置全局会话选项
	store.Options(sessions.Options{
		Path:     "/",   // 整个网站都可以访问
		MaxAge:   86400, // 1天
		HttpOnly: true,  // 防止JavaScript访问
		// Secure: true,  // 生产环境应设置为true
	})

	router.Use(sessions.Sessions("mysession", store))

	// 使用自定义的 LoginMiddlewareBuilder 中间件
	router.Use(middleware.NewLoginMiddlewareBuilder().Build())

	// 注册用户相关的路由
	u.RegisterRoutes(router)

	return router
}
