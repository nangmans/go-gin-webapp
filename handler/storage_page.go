package handler

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nangmans14/gin-web/model/gcs"
)

const projectID string = "sandbox-393317"

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
		fmt.Printf("%+v\n", err)
	}
	switch {
	// Error handling for string format
	case strings.Contains(fmt.Sprint(err), "could not find default credentials"):
		Render(c, gin.H{}, "credential_not_found.html")
	case buckets == nil:
		Render(c, gin.H{}, "bucket_not_found.html")
	default:
		Render(c, gin.H{
			"title":   "Home Page",
			"payload": buckets,
		}, "index.html")
	}
}

func ShowStoragePage(c *gin.Context) {

	bucket, err := gcs.ListBuckets(ioutil.Discard, projectID, c.Param("bucket_id"))
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	objects, err := gcs.ListObjects(ioutil.Discard, bucket[0])
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	bucket[0].Objects = objects

	// If any objects are not found
	if objects == nil {
		Render(c, gin.H{}, "object_not_found.html")
	} else {
		Render(c, gin.H{
			"object_id": "",
			"bucket_id": bucket[0].Name,
			"payload":   objects,
		}, "article.html")
	}
}

func ShowObjectPage(c *gin.Context) {

	name := strings.TrimPrefix(c.Param("object_id"), "/")
	bucket, err := gcs.ListBuckets(ioutil.Discard, projectID, c.Param("bucket_id"))
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	objects, err := gcs.ListObjects(ioutil.Discard, bucket[0], name)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	//bucket[0].Objects = objects
	switch {
	case objects == nil:
		Render(c, gin.H{}, "object_not_found.html")
	case name[len(name)-1] == '/':
		Render(c, gin.H{
			"object_id": name,
			"bucket_id": c.Param("bucket_id"),
			"payload":   objects,
		}, "article.html")
	default:
		Render(c, gin.H{
			"object_id": name,
			"bucket_id": c.Param("bucket_id"),
			"payload":   objects[0].Metadata,
		}, "object.html")
	}
}
