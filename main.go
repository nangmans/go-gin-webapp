package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nangmans14/gin-web/handler"
)

var router *gin.Engine

func main() {
	router = gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.Use(handler.SetBucket)

	router.GET("/", handler.ShowIndexPage)

	router.GET("/storage/:bucket_id", handler.ShowStoragePage)

	router.GET("/storage/:bucket_id/object/:object_name", handler.ShowObjectPage)

	router.Run()
}
