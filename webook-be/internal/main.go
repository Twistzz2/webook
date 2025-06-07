package main

import (
	"github.com/Twistzz2/webook/webook-be/internal/web"
)

func main() {
	router := web.InitRouter()

	router.Run(":8080")
}
