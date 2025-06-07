package web

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	u := &UserHandler{}
	u.RegisterRoutes(router)

	return router
}
