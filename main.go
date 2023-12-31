package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nangmans14/gin-web/handler"
)

var router *gin.Engine

func main() {
	router = gin.Default()

	// router.UnescapePathValues = false
	// // router.UseRawPath = true
	// router.RemoveExtraSlash = false
	// router.RedirectTrailingSlash = true

	router.LoadHTMLGlob("templates/*")

	router.GET("/", handler.ShowIndexPage)

	router.GET("/storage/:bucket_id", handler.ShowBucketPage)

	router.GET("/storage/:bucket_id/object/*object_id", handler.ShowObjectPage)

	router.Run()
}
