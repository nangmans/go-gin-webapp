package handler

import (
	"fmt"
	"io/ioutil"
	"strings"

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

	fmt.Print(c.Param("bucket_id"))
	bucket, err := gcs.ListBuckets(ioutil.Discard, projectID, c.Param("bucket_id"))
	if err != nil {
		fmt.Printf("GetBucketByName: %s", err)
	}

	objects, err := gcs.ListObjects(ioutil.Discard, bucket[0])
	if err != nil {
		fmt.Printf("ListObjects: %s", err)
	}

	bucket[0].Objects = objects
	fmt.Print(objects[3].Name)
	Render(c, gin.H{
		"object_id": "",
		"bucket_id": bucket[0].Name,
		"payload":   objects,
	}, "article.html")

}

func ShowObjectPage(c *gin.Context) {
	fmt.Print(c.Param("object_id"))

	name := strings.TrimPrefix(c.Param("object_id"), "/")
	bucket, err := gcs.ListBuckets(ioutil.Discard, projectID, c.Param("bucket_id"))
	if err != nil {
		fmt.Printf("GetBucketByName: %s", err)
	}

	objects, err := gcs.ListObjects(ioutil.Discard, bucket[0], name)
	if err != nil {
		fmt.Printf("ListObjects: %s", err)
	}

	bucket[0].Objects = objects
	fmt.Print(name)
	if name[len(name)-1] == '/' {
		Render(c, gin.H{
			"object_id": name,
			"bucket_id": c.Param("bucket_id"),
			"payload":   objects,
		}, "article.html")
	} else {
		Render(c, gin.H{
			"object_id": name,
			"bucket_id": c.Param("bucket_id"),
			"payload":   objects,
		}, "object.html")
	}

}
