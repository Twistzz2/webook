package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	u := &web.UserHandler{}
	u.RegisterRoutes(router)

	router.Run(":8080")
}
