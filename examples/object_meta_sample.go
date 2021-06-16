/**
 * This sample demonstrates how to set/get self-defined metadata for object
 * on WOS using the WOS SDK for Go.
 */
package examples

import (
	"fmt"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk-v2/wos"
	"strings"
)

type ObjectMetaSample struct {
	bucketName string
	objectKey  string
	wosClient  *wos.WosClient
}

func newObjectMetaSample(ak, sk, endpoint, bucketName, objectKey string) *ObjectMetaSample {
	wosClient, err := wos.New(ak, sk, endpoint)
	if err != nil {
		panic(err)
	}
	return &ObjectMetaSample{wosClient: wosClient, bucketName: bucketName, objectKey: objectKey}
}

func (sample ObjectMetaSample) SetObjectMeta() {
	input := &wos.PutObjectInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.Body = strings.NewReader("Hello WOS")
	// Setting object mime type
	input.ContentType = "text/plain"
	// Setting self-defined metadata
	input.Metadata = map[string]string{"meta1": "value1", "meta2": "value2"}
	_, err := sample.wosClient.PutObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Set object meatdata successfully!")
	fmt.Println()
}

func (sample ObjectMetaSample) GetObjectMeta() {
	input := &wos.GetObjectMetadataInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	output, err := sample.wosClient.GetObjectMetadata(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Object content-type:%s\n", output.ContentType)
	for key, val := range output.Metadata {
		fmt.Printf("%s:%s\n", key, val)
	}
	fmt.Println()
}
func (sample ObjectMetaSample) DeleteObject() {
	input := &wos.DeleteObjectInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey

	_, err := sample.wosClient.DeleteObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Delete object:%s successfully!\n", sample.objectKey)
	fmt.Println()
}

func RunObjectMetaSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
		objectKey  = "object-test"
	)
	sample := newObjectMetaSample(ak, sk, endpoint, bucketName, objectKey)

	sample.SetObjectMeta()

	sample.GetObjectMeta()

	sample.DeleteObject()
}
