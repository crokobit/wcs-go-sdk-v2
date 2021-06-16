/**
 * This sample demonstrates how to upload multiparts to WOS
 * using the WOS SDK for Go.
 */
package examples

import (
	"fmt"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk-v2/wos"
	"strings"
)

type SimpleMultipartUploadSample struct {
	bucketName string
	objectKey  string
	wosClient  *wos.WosClient
}

func newSimpleMultipartUploadSample(ak, sk, endpoint, bucketName, objectKey string) *SimpleMultipartUploadSample {
	wosClient, err := wos.New(ak, sk, endpoint)
	if err != nil {
		panic(err)
	}
	return &SimpleMultipartUploadSample{wosClient: wosClient, bucketName: bucketName, objectKey: objectKey}
}

func (sample SimpleMultipartUploadSample) InitiateMultipartUpload() string {
	input := &wos.InitiateMultipartUploadInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	output, err := sample.wosClient.InitiateMultipartUpload(input)
	if err != nil {
		panic(err)
	}
	return output.UploadId
}

func (sample SimpleMultipartUploadSample) UploadPart(uploadId string) (string, int) {
	input := &wos.UploadPartInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.UploadId = uploadId
	input.PartNumber = 1
	input.Body = strings.NewReader("Hello WOS")
	output, err := sample.wosClient.UploadPart(input)
	if err != nil {
		panic(err)
	}
	return output.ETag, output.PartNumber
}

func (sample SimpleMultipartUploadSample) CompleteMultipartUpload(uploadId, etag string, partNumber int) {
	input := &wos.CompleteMultipartUploadInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.UploadId = uploadId
	input.Parts = []wos.Part{
		wos.Part{PartNumber: partNumber, ETag: etag},
	}
	_, err := sample.wosClient.CompleteMultipartUpload(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Upload object %s successfully!\n", sample.objectKey)
}

func RunSimpleMultipartUploadSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
		objectKey  = "object-test"
	)
	sample := newSimpleMultipartUploadSample(ak, sk, endpoint, bucketName, objectKey)

	// Step 1: initiate multipart upload
	fmt.Println("Step 1: initiate multipart upload")
	uploadId := sample.InitiateMultipartUpload()

	// Step 2: upload a part
	fmt.Println("Step 2: upload a part")

	etag, partNumber := sample.UploadPart(uploadId)

	// Step 3: complete multipart upload
	fmt.Println("Step 3: complete multipart upload")
	sample.CompleteMultipartUpload(uploadId, etag, partNumber)

}
