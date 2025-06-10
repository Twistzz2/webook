package main

import (
	"github.com/Twistzz2/webook/webook-be/internal/handler"
)

func main() {
	router := handler.InitRouter()

	router.Run(":8080")
}
