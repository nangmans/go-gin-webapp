package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Render(c *gin.Context, data gin.H, templateName string) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data["object"])
	case "application/xml":
		c.XML(http.StatusOK, data["object"])
	default:
		c.HTML(http.StatusOK, templateName, data)
	}
}
