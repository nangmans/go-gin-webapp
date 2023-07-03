package gcs_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/nangmans14/gin-web/model/gcs"
)

const projectID string = "prj-cc-sandbox-devops-0010"

func TestListBuckets(t *testing.T) {

	buckets, err := gcs.ListBuckets(ioutil.Discard, projectID)
	if err != nil {
		t.Fatalf("listBuckets: %v", err)
	}
	fmt.Printf("buckets : %s", buckets)
}

func TestListObjects(t *testing.T) {

	objects, err := gcs.ListObjects(ioutil.Discard, "iso-collection")
	if err != nil {
		t.Fatalf("listObjects: %v", err)
	}
	fmt.Printf("Objects: %s", objects)
}
