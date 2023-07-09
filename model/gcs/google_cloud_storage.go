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

func ListBuckets(w io.Writer, projectID string, name ...string) ([]*Bucket, error) {
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

	if name != nil {
		it.Prefix = name[0]
	}

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
	Metadata *storage.ObjectAttrs
}

func newObject(b *Bucket, n string) (*Object, error) {
	var (
		lstIdx       int
		name, parent string
	)

	if n[len(n)-1] != '/' {
		lstIdx = strings.LastIndex(n, "/")
	} else {
		lstIdx = strings.LastIndex(n[:len(n)-1], "/")
	}

	if n[:lstIdx+1] == "" {
		name = n[lstIdx+1:]
		parent = b.Name
	} else {
		name = n[lstIdx+1:]
		parent = n[:lstIdx+1]
	}

	return &Object{
		Name:   name,
		Parent: parent,
	}, nil

}

func GetObjectByName(o []*Object, p string) (*Object, error) {
	for _, a := range o {
		if a.Parent == p {
			return a, nil
		}
	}
	return nil, errors.New("Object Not Found")
}

func ListObjects(w io.Writer, b *Bucket, q ...string) ([]*Object, error) {
	var (
		objects []*Object
		name    string
	)

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	if len(q) != 0 {
		name = q[0]
	}
	fmt.Printf("query is %s", name)
	query := &storage.Query{
		Delimiter:                "/",
		IncludeTrailingDelimiter: true,
		Prefix:                   name,
	}

	it := client.Bucket(b.Name).Objects(ctx, query)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		name := attrs.Name

		switch {
		case len(name) == 0:
			name = attrs.Prefix
		case name[len(name)-1] == '/':
			continue
		}

		object, err := newObject(b, name)
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
