package web

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	u := NewUserHandler()
	u.RegisterRoutes(router)

	return router
}
