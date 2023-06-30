package handler

import (
	"fmt"
	"io/ioutil"

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

// func GetBucket(c *gin.Context) {
// 	if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
// 		if article, err := article.GetArticleByID(articleID); err == nil {
// 			c.HTML(
// 				http.StatusOK,
// 				"article.html",
// 				gin.H{
// 					"title":   article.Title,
// 				},
// 			)
// 		} else {
// 			c.AbortWithError(http.StatusNotFound, err)
// 		}
// 	} else {
// 		c.AbortWithStatus(http.StatusNotFound)
// 	}
// }
