package main

import (
	_ "Praiseson6065/ocrolus-be/config"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	err := Server()

	if err != nil {
		log.Fatal(err)
	}
}
