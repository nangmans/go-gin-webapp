package handler

import (
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/nangmans14/gin-web/model/gcs"
)

const projectID string = "prj-prod-datadev-8411"

// func SetBucket(c *gin.Context) {
// 	buckets, err := gcs.ListBuckets(ioutil.Discard, projectID)
// 	if err != nil {
// 		fmt.Printf("ListBuckets: %s", err)
// 	}
// 	c.Set("buckets", buckets)
// 	c.Next()
// }

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

	bucket, err := gcs.ListBuckets(ioutil.Discard, projectID, c.Param("bucket_id"))
	if err != nil {
		fmt.Printf("GetBucketByName: %s", err)
	}

	objects, err := gcs.ListObjects(ioutil.Discard, bucket[0])
	if err != nil {
		fmt.Printf("ListObjects: %s", err)
	}

	bucket[0].Objects = objects

	Render(c, gin.H{
		"bucket_id": bucket[0].Name,
		"payload":   objects,
	}, "article.html")

}

func ShowObjectPage(c *gin.Context) {

	bucket, err := gcs.ListBuckets(ioutil.Discard, projectID, c.Param("bucket_id"))
	if err != nil {
		fmt.Printf("GetBucketByName: %s", err)
	}

	objects, err := gcs.ListObjects(ioutil.Discard, bucket[0], c.Param("object_id"))
	if err != nil {
		fmt.Printf("ListObjects: %s", err)
	}

	bucket[0].Objects = objects

	Render(c, gin.H{
		"payload": objects,
	}, "article.html")

}
