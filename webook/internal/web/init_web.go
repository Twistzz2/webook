package web

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes() *gin.Engine {
	router := gin.Default()

	registerUsersRoutes(router)

	return router
}

func registerUsersRoutes(router *gin.Engine) {
	userHandler := &UserHandler{}
	userHandler.RegisterRoutes(router)

	// 这里可以注册其他模块的路由
	// registerOtherModuleRoutes(router)
}
