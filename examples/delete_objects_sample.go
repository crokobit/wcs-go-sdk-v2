/**
 * This sample demonstrates how to delete objects under specified bucket
 * from WOS using the WOS SDK for Go.
 */
package examples

import (
	"fmt"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk-v2/wos"
	"strconv"
	"strings"
)

const (
	MyObjectKey string = "MyObjectKey"
)

type DeleteObjectsSample struct {
	bucketName string
	wosClient  *wos.WosClient
}

func newDeleteObjectsSample(ak, sk, endpoint, bucketName string) *DeleteObjectsSample {
	wosClient, err := wos.New(ak, sk, endpoint)
	if err != nil {
		panic(err)
	}
	return &DeleteObjectsSample{wosClient: wosClient, bucketName: bucketName}
}

func (sample DeleteObjectsSample) BatchPutObjects() {
	content := "Thank you for using Object Storage Service"
	keyPrefix := MyObjectKey

	input := &wos.PutObjectInput{}
	input.Bucket = sample.bucketName
	input.Body = strings.NewReader(content)
	for i := 0; i < 100; i++ {
		input.Key = keyPrefix + strconv.Itoa(i)
		_, err := sample.wosClient.PutObject(input)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Succeed to put object %s\n", input.Key)
	}
}

func (sample DeleteObjectsSample) BatchDeleteObjects() {
	input := &wos.ListObjectsInput{}
	input.Bucket = sample.bucketName
	output, err := sample.wosClient.ListObjects(input)
	if err != nil {
		panic(err)
	}
	objects := make([]wos.ObjectToDelete, 0, len(output.Contents))
	for _, content := range output.Contents {
		objects = append(objects, wos.ObjectToDelete{Key: content.Key})
	}
	deleteObjectsInput := &wos.DeleteObjectsInput{}
	deleteObjectsInput.Bucket = sample.bucketName
	deleteObjectsInput.Objects = objects[:]
	deleteObjectsOutput, err := sample.wosClient.DeleteObjects(deleteObjectsInput)
	if err != nil {
		panic(err)
	}
	for _, deleted := range deleteObjectsOutput.Deleteds {
		fmt.Printf("Delete %s successfully\n", deleted.Key)
	}
	fmt.Println()
	for _, deleteError := range deleteObjectsOutput.Errors {
		fmt.Printf("Delete %s failed, code:%s, message:%s\n", deleteError.Key, deleteError.Code, deleteError.Message)
	}
	fmt.Println()
}

func RunDeleteObjectsSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
	)
	sample := newDeleteObjectsSample(ak, sk, endpoint, bucketName)

	// Batch put objects into the bucket
	sample.BatchPutObjects()

	// Delete all objects uploaded recently under the bucket
	sample.BatchDeleteObjects()
}
