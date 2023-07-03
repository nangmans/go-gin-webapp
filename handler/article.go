package handler

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/nangmans14/gin-web/model/gcs"
)

const projectID string = "prj-prod-datadev-8411"

// This comment is required for go doc when the function is exported.
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

func ShowStoragePage(c *gin.Context) {
	objects, err := gcs.ListObjects(ioutil.Discard, *gcs.bucket)
	if err != nil {
		fmt.Printf("ListObjects: %s", err)
	}

	Render(c, gin.H{
		"bucket_id": c.Param("bucket_id"),
		"payload":   objects,
	}, "article.html")

}

func ShowObjectPage(c *gin.Context) {
	attrs, err := gcs.GetMetadata(ioutil.Discard, c.Param("bucket_id"), c.Param("object_name"))
	if err != nil {
		fmt.Printf("GetMetadata: %s", err)
	}

	Render(c, gin.H{
		"payload": attrs,
	}, "article.html")
}
