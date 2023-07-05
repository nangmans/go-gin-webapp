package gcs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type Bucket struct {
	Name    string
	Objects [](*Object)
}

func newBucket(n string) (*Bucket, error) {
	return &Bucket{
		Name: n,
	}, nil
}

// func (b *Bucket) Append(o string) (*Bucket, error) {
// 	obj, err := newObject(b, o)
// 	if err != nil {
// 		return nil, err
// 	}
// 	b.Objects = append(b.Objects, obj)

// 	return b, nil
// }

func GetBucketByName(b []*Bucket, n string) (*Bucket, error) {
	for _, a := range b {
		if a.Name == n {
			return a, nil
		}
	}
	return nil, errors.New("Bucket Not Found")
}

func ListBuckets(w io.Writer, projectID string) ([]*Bucket, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancle := context.WithTimeout(ctx, time.Second*30)
	defer cancle()

	var buckets []*Bucket
	it := client.Buckets(ctx, projectID)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		bucket, err := newBucket(battrs.Name)
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, bucket)
	}
	return buckets, nil
}

type Object struct {
	Name     string
	Parent   string
	IsRoot   bool
	Metadata *storage.ObjectAttrs
}

func newObject(b *Bucket, n string) (*Object, error) {
	lstidx := strings.LastIndex(n, "/")
	name := n[lstidx+1:]
	parent := n[:lstidx+1]
	return &Object{
		Name:   name,
		Parent: parent,
	}, nil
}

func GetObjectByName(o []*Object, n string) (*Object, error) {
	for _, a := range o {
		if a.Name == n {
			return a, nil
		}
	}
	return nil, errors.New("Object Not Found")
}

func ListObjects(w io.Writer, bucket *Bucket, parent ...string) ([]*Object, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var objects []*Object
	query := &storage.Query{
		Prefix:    "TM/",
		Delimiter: "/",
	}

	it := client.Bucket(bucket.Name).Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		object, err := newObject(bucket, attrs.Name)
		if err != nil {
			return nil, err
		}
		objects = append(objects, object)
	}
	return objects, nil
}

func GetMetadata(w io.Writer, bucket, object string) (*storage.ObjectAttrs, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(bucket).Object(object)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).Attrs: %w", object, err)
	}
	return attrs, nil
}
