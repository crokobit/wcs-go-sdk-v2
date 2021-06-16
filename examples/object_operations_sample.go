/**
 * This sample demonstrates how to do object-related operations
 * (such as create/delete/get/copy object, do object ACL)
 * on WOS using the WOS SDK for Go.
 */
package examples

import (
	"fmt"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk-v2/wos"
	"io/ioutil"
	"strings"
)

type ObjectOperationsSample struct {
	bucketName string
	objectKey  string
	wosClient  *wos.WosClient
}

func newObjectOperationsSample(ak, sk, endpoint, bucketName, objectKey string) *ObjectOperationsSample {
	wosClient, err := wos.New(ak, sk, endpoint)
	if err != nil {
		panic(err)
	}
	return &ObjectOperationsSample{wosClient: wosClient, bucketName: bucketName, objectKey: objectKey}
}

func (sample ObjectOperationsSample) GetObjectMeta() {
	input := &wos.GetObjectMetadataInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	output, err := sample.wosClient.GetObjectMetadata(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Object content-type:%s\n", output.ContentType)
	fmt.Printf("Object content-length:%d\n", output.ContentLength)
	fmt.Println()
}

func (sample ObjectOperationsSample) CreateObject() {
	input := &wos.PutObjectInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.Body = strings.NewReader("Hello WOS")

	_, err := sample.wosClient.PutObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Create object:%s successfully!\n", sample.objectKey)
	fmt.Println()
}

func (sample ObjectOperationsSample) GetObject() {
	input := &wos.GetObjectInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey

	output, err := sample.wosClient.GetObject(input)
	if err != nil {
		panic(err)
	}
	defer func() {
		errMsg := output.Body.Close()
		if errMsg != nil {
			panic(errMsg)
		}
	}()
	fmt.Println("Object content:")
	body, err := ioutil.ReadAll(output.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	fmt.Println()
}

func (sample ObjectOperationsSample) CopyObject() {
	input := &wos.CopyObjectInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey + "-back"
	input.CopySourceBucket = sample.bucketName
	input.CopySourceKey = sample.objectKey

	_, err := sample.wosClient.CopyObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Copy object successfully!")
	fmt.Println()
}

func (sample ObjectOperationsSample) DeleteObject() {
	input := &wos.DeleteObjectInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey

	_, err := sample.wosClient.DeleteObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Delete object:%s successfully!\n", input.Key)
	fmt.Println()

	input.Key = sample.objectKey + "-back"

	_, err = sample.wosClient.DeleteObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Delete object:%s successfully!\n", input.Key)
	fmt.Println()
}

func RunObjectOperationsSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
		objectKey  = "object-test"
	)

	sample := newObjectOperationsSample(ak, sk, endpoint, bucketName, objectKey)

	sample.CreateObject()

	sample.GetObjectMeta()

	sample.GetObject()

	sample.CopyObject()

	sample.DeleteObject()
}
