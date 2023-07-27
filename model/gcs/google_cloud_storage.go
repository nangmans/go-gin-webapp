package gcs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

func GetBucketByName(b []*Bucket, n string) (*Bucket, error) {
	for _, a := range b {
		if a.Name == n {
			return a, nil
		}
	}
	return nil, errors.New("Bucket Not Found")
}

func ListBuckets(w io.Writer, projectID string, name ...string) ([]*Bucket, error) {

	var buckets []*Bucket

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancle := context.WithTimeout(ctx, time.Second*30)
	defer cancle()

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
	Metadata *Metadata
}

type Metadata struct {
	Name    string
	Bucket  string
	Size    int64
	Created time.Time
	Updated time.Time
}

func newObjectWithMetadata(b *Bucket, n string) (*Object, error) {
	var (
		lstIdx       int
		objAttr      *Metadata
		err          error
		name, parent string
	)

	// Check whether n is folder or file
	if n[len(n)-1] != '/' {
		// If n is file
		lstIdx = strings.LastIndex(n, "/")
		objAttr, err = GetMetadata(ioutil.Discard, b, n)
		if err != nil {
			return nil, fmt.Errorf("GetMetadata: %w", err)
		}
	} else {
		// If n is folder
		lstIdx = strings.LastIndex(n[:len(n)-1], "/")
	}

	// Check whether n is at root or not
	if n[:lstIdx+1] == "" {
		name = n[lstIdx+1:]
		parent = b.Name
	} else {
		name = n[lstIdx+1:]
		parent = n[:lstIdx+1]
	}

	return &Object{
		Name:     name,
		Parent:   parent,
		Metadata: objAttr,
	}, nil

}

func newObjectWithoutMetadata(b *Bucket, n string) (*Object, error) {
	var (
		lstIdx       int
		name, parent string
	)

	// Check whether n is folder or file
	if n[len(n)-1] != '/' {
		// If n is file
		lstIdx = strings.LastIndex(n, "/")
	} else {
		// If n is folder
		lstIdx = strings.LastIndex(n[:len(n)-1], "/")
	}

	// Check whether n is at root or not
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

func GetObject(b *Bucket, n string) (*Object, error) {
	o, err := newObjectWithMetadata(b, n)
	if err != nil {
		return nil, fmt.Errorf("newObject: %w", err)
	}

	return o, nil
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

		// name 변수를 아래 if문에서 할당하면 if문 빠져나오는 순간 값이 사라지므로 따로 선언한다.
		name := attrs.Name

		if len(name) == 0 {
			name = attrs.Prefix
		}

		object, err := newObjectWithoutMetadata(b, name)
		if err != nil {
			return nil, err
		}
		objects = append(objects, object)
	}
	return objects, nil
}

func GetMetadata(w io.Writer, b *Bucket, object string) (*Metadata, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(b.Name).Object(object)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).Attrs: %w", object, err)
	}
	return &Metadata{
		Name:    attrs.Name,
		Bucket:  attrs.Bucket,
		Size:    attrs.Size,
		Created: attrs.Created,
		Updated: attrs.Updated,
	}, nil
}
