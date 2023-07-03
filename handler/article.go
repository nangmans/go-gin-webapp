package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nangmans14/gin-web/model/gcs"
)

const projectID string = "prj-cc-sandbox-devops-0010"

func ShowIndexPage(c *gin.Context) {
	buckets, err := gcs.ListBuckets(ioutil.Discard, projectID)
	if err != nil {
		fmt.Printf("ListBuckets: %s", err)
	}

	Render(c, gin.H{
		"title":   "Home Page",
		"payload": buckets,
	}, "index.html")
}

func GetBucketObjects(c *gin.Context) {
	if objects, err := gcs.ListObjects(ioutil.Discard, c.Param("bucket_id")); err == nil {
		c.HTML(
			http.StatusOK,
			"article.html",
			gin.H{
				"name": objects,
			},
		)
	} else {
		c.AbortWithError(http.StatusNotFound, err)
	}
}
