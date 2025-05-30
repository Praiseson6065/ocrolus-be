package main

import (
	"Praiseson6065/ocrolus-be/config"
	"Praiseson6065/ocrolus-be/middleware"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Server() error {
	// Get port from config
	port := config.Config.Server.Port

	r := gin.New()
	r.Use(middleware.CORS())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Hello": "World",
		})
	})
	AuthRouter(r)
	ApiRouter(r)
	fmt.Println("Server is starting on 8000")
	err := r.Run(port)
	if err != nil {
		return err
	}
	return nil
}
