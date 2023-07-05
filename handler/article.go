package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nangmans14/gin-web/model/gcs"
)

const projectID string = "prj-prod-datadev-8411"

func SetBucket(c *gin.Context) {
	buckets, err := gcs.ListBuckets(ioutil.Discard, projectID)
	if err != nil {
		fmt.Printf("ListBuckets: %s", err)
	}
	c.Set("buckets", buckets)
	c.Next()
}

// This comment is required for go doc when the function is exported.
func ShowIndexPage(c *gin.Context) {
	buckets := c.MustGet("buckets").([]*gcs.Bucket)

	Render(c, gin.H{
		"title":   "Home Page",
		"payload": buckets,
	}, "index.html")

}

func ShowStoragePage(c *gin.Context) {
	buckets := c.MustGet("buckets").([]*gcs.Bucket)

	bucket, err := gcs.GetBucketByName(buckets, c.Param("bucket_id"))
	if err != nil {
		fmt.Printf("GetBucketByName: %s", err)
	}

	objects, err := gcs.ListObjects(ioutil.Discard, bucket)
	if err != nil {
		fmt.Printf("ListObjects: %s", err)
	}
	bucket.Objects = objects

	// Must be located before Render()
	c.Set("objects", objects)

	Render(c, gin.H{
		"bucket_id": bucket.Name,
		"payload":   bucket.Objects,
	}, "article.html")

}

func ShowObjectPage(c *gin.Context) {
	objects := c.MustGet("objects").([]*gcs.Object)

	object, err := gcs.GetObjectByName(objects, c.Param("object_name"))
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
	} else {
		Render(c, gin.H{
			"payload": object,
		}, "object.html")
	}
}
