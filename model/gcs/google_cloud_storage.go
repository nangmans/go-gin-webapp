package gcs

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type bucket struct {
	Name    string
	Objects [](*object)
}

func newBucket(n string) (*bucket, error) {
	return &bucket{
		Name: n,
	}, nil
}

func (b *bucket) Append(o string) (*bucket, error) {
	obj, err := newObject(b, o)
	if err != nil {
		return nil, err
	}
	b.Objects = append(b.Objects, obj)

	return b, nil
}

type object struct {
	name     string
	metadata *storage.ObjectAttrs
}

func newObject(b *bucket, n string) (*object, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	attrs, err := client.Bucket(b.Name).Object(n).Attrs(ctx)
	if err != nil {
		return nil, err
	}
	return &object{
		name:     attrs.Name,
		metadata: attrs,
	}, nil
}

func ListBuckets(w io.Writer, projectID string) ([]*bucket, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancle := context.WithTimeout(ctx, time.Second*30)
	defer cancle()

	var buckets []*bucket
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

func ListObjects(w io.Writer, bucket *bucket) ([]*object, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	var objects []*object
	it := client.Bucket(bucket.Name).Objects(ctx, nil)
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
