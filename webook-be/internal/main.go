package main

import (
	"time"

	"github.com/Twistzz2/webook/webook-be/internal/web"
	"github.com/gin-contrib/cors"
)

func main() {
	router := web.InitRouter()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://foo.com"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.Run(":8080")
}
