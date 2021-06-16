/**
 * This sample demonstrates how to list objects under specified bucket
 * from WOS using the WOS SDK for Go.
 */
package examples

import (
	"fmt"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk-v2/wos"
	"strconv"
	"strings"
)

type ListObjectsSample struct {
	bucketName string
	wosClient  *wos.WosClient
}

func newListObjectsSample(ak, sk, endpoint, bucketName string) *ListObjectsSample {
	wosClient, err := wos.New(ak, sk, endpoint)
	if err != nil {
		panic(err)
	}
	return &ListObjectsSample{wosClient: wosClient, bucketName: bucketName}
}

func (sample ListObjectsSample) DoInsertObjects() []string {

	keyPrefix := "MyObjectKey"

	input := &wos.PutObjectInput{}
	input.Bucket = sample.bucketName
	input.Body = strings.NewReader("Hello WOS")
	keys := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		input.Key = keyPrefix + strconv.Itoa(i)
		_, err := sample.wosClient.PutObject(input)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Succeed to put object %s\n", input.Key)
		keys = append(keys, input.Key)
	}
	fmt.Println()
	return keys
}

func (sample ListObjectsSample) ListObjects() {
	input := &wos.ListObjectsInput{}
	input.Bucket = sample.bucketName
	output, err := sample.wosClient.ListObjects(input)
	if err != nil {
		panic(err)
	}
	for index, val := range output.Contents {
		fmt.Printf("Content[%d]-ETag:%s, Key:%s, Size:%d\n",
			index, val.ETag, val.Key, val.Size)
	}
	fmt.Println()
}

func (sample ListObjectsSample) ListObjectsByMarker() {
	input := &wos.ListObjectsInput{}
	input.Bucket = sample.bucketName
	input.MaxKeys = 10
	output, err := sample.wosClient.ListObjects(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("List the first 10 objects :")
	for index, val := range output.Contents {
		fmt.Printf("Content[%d]-ETag:%s, Key:%s, Size:%d\n",
			index, val.ETag, val.Key, val.Size)
	}
	fmt.Println()

	input.Marker = output.NextMarker
	output, err = sample.wosClient.ListObjects(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("List the second 10 objects using marker:")
	for index, val := range output.Contents {
		fmt.Printf("Content[%d]-ETag:%s, Key:%s, Size:%d\n",
			index, val.ETag, val.Key, val.Size)
	}
	fmt.Println()
}

func (sample ListObjectsSample) ListObjectsByPage() {

	pageSize := 10
	pageNum := 1
	input := &wos.ListObjectsInput{}
	input.Bucket = sample.bucketName
	input.MaxKeys = pageSize

	for {
		output, err := sample.wosClient.ListObjects(input)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Page:%d\n", pageNum)
		for index, val := range output.Contents {
			fmt.Printf("Content[%d]-ETag:%s, Key:%s, Size:%d\n",
				index, val.ETag, val.Key, val.Size)
		}
		if output.IsTruncated {
			input.Marker = output.NextMarker
			pageNum++
		} else {
			break
		}
	}

	fmt.Println()
}

func (sample ListObjectsSample) DeleteObjects(keys []string) {
	input := &wos.DeleteObjectsInput{}
	input.Bucket = sample.bucketName

	objects := make([]wos.ObjectToDelete, 0, len(keys))
	for _, key := range keys {
		objects = append(objects, wos.ObjectToDelete{Key: key})
	}
	input.Objects = objects
	_, err := sample.wosClient.DeleteObjects(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Delete objects successfully!")
}

func RunListObjectsSample() {

	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
	)

	sample := newListObjectsSample(ak, sk, endpoint, bucketName)

	// First insert 100 objects for demo
	keys := sample.DoInsertObjects()

	// List objects using default parameters, will return up to 1000 objects
	sample.ListObjects()

	// List the first 10 and second 10 objects
	sample.ListObjectsByMarker()

	// List objects in way of pagination
	sample.ListObjectsByPage()

	// Delete all the objects created
	sample.DeleteObjects(keys)
}
