/**
 * This sample demonstrates how to multipart upload an object concurrently by copy mode
 * to WOS using the WOS SDK for Go.
 */
package examples

import (
	"errors"
	"fmt"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk-v2/wos"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const (
	filePathSample string = "/temp/text.txt"
)

type ConcurrentCopyPartSample struct {
	bucketName string
	objectKey  string
	wosClient  *wos.WosClient
}

func newConcurrentCopyPartSample(ak, sk, endpoint, bucketName, objectKey string) *ConcurrentCopyPartSample {
	wosClient, err := wos.New(ak, sk, endpoint)
	if err != nil {
		panic(err)
	}
	return &ConcurrentCopyPartSample{wosClient: wosClient, bucketName: bucketName, objectKey: objectKey}
}

func (sample ConcurrentCopyPartSample) checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func (sample ConcurrentCopyPartSample) createSampleFile(sampleFilePath string, byteCount int64) {
	if err := os.MkdirAll(filepath.Dir(sampleFilePath), os.ModePerm); err != nil {
		panic(err)
	}

	fd, err := os.OpenFile(sampleFilePath, os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(errors.New("open file with error"))
	}

	const chunkSize = 1024
	b := [chunkSize]byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < chunkSize; i++ {
		b[i] = uint8(r.Intn(255))
	}

	var writedCount int64
	for {
		remainCount := byteCount - writedCount
		if remainCount <= 0 {
			break
		}
		if remainCount > chunkSize {
			_, errMsg := fd.Write(b[:])
			sample.checkError(errMsg)
			writedCount += chunkSize
		} else {
			_, errMsg := fd.Write(b[:remainCount])
			sample.checkError(errMsg)
			writedCount += remainCount
		}
	}

	defer func() {
		errMsg := fd.Close()
		sample.checkError(errMsg)
	}()

	err = fd.Sync()
	sample.checkError(err)
}

func (sample ConcurrentCopyPartSample) PutFile(sampleFilePath string) {
	input := &wos.PutFileInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.SourceFile = sampleFilePath
	_, err := sample.wosClient.PutFile(input)
	if err != nil {
		panic(err)
	}
}

func (sample ConcurrentCopyPartSample) DoConcurrentCopyPart() {
	destBucketName := sample.bucketName
	destObjectKey := sample.objectKey + "-back"
	sourceBucketName := sample.bucketName
	sourceObjectKey := sample.objectKey
	// Claim a upload id firstly
	input := &wos.InitiateMultipartUploadInput{}
	input.Bucket = destBucketName
	input.Key = destObjectKey
	output, err := sample.wosClient.InitiateMultipartUpload(input)
	if err != nil {
		panic(err)
	}
	uploadId := output.UploadId

	fmt.Printf("Claiming a new upload id %s\n", uploadId)
	fmt.Println()

	// Get size of the object
	getObjectMetadataInput := &wos.GetObjectMetadataInput{}
	getObjectMetadataInput.Bucket = sourceBucketName
	getObjectMetadataInput.Key = sourceObjectKey
	getObjectMetadataOutput, err := sample.wosClient.GetObjectMetadata(getObjectMetadataInput)
	if err != nil {
		panic(err)
	}

	objectSize := getObjectMetadataOutput.ContentLength

	// Calculate how many blocks to be divided
	// 5MB
	var partSize int64 = 5 * 1024 * 1024
	partCount := int(objectSize / partSize)

	if objectSize%partSize != 0 {
		partCount++
	}

	fmt.Printf("Total parts count %d\n", partCount)
	fmt.Println()

	//  Upload multiparts by copy mode
	fmt.Println("Begin to upload multiparts to WOS by copy mode")

	partChan := make(chan wos.Part, 5)

	for i := 0; i < partCount; i++ {
		partNumber := i + 1
		rangeStart := int64(i) * partSize
		rangeEnd := rangeStart + partSize - 1
		if i+1 == partCount {
			rangeEnd = objectSize - 1
		}
		go func(start, end int64, index int) {
			copyPartInput := &wos.CopyPartInput{}
			copyPartInput.Bucket = destBucketName
			copyPartInput.Key = destObjectKey
			copyPartInput.UploadId = uploadId
			copyPartInput.PartNumber = index
			copyPartInput.CopySourceBucket = sourceBucketName
			copyPartInput.CopySourceKey = sourceObjectKey
			copyPartInput.CopySourceRangeStart = start
			copyPartInput.CopySourceRangeEnd = end
			copyPartOutput, errMsg := sample.wosClient.CopyPart(copyPartInput)
			if errMsg == nil {
				fmt.Printf("%d finished\n", index)
				partChan <- wos.Part{ETag: copyPartOutput.ETag, PartNumber: copyPartOutput.PartNumber}
			} else {
				panic(errMsg)
			}
		}(rangeStart, rangeEnd, partNumber)
	}

	parts := make([]wos.Part, 0, partCount)

	for {
		part, ok := <-partChan
		if !ok {
			break
		}
		parts = append(parts, part)
		if len(parts) == partCount {
			close(partChan)
		}
	}

	fmt.Println()
	fmt.Println("Completing to upload multiparts")
	completeMultipartUploadInput := &wos.CompleteMultipartUploadInput{}
	completeMultipartUploadInput.Bucket = destBucketName
	completeMultipartUploadInput.Key = destObjectKey
	completeMultipartUploadInput.UploadId = uploadId
	completeMultipartUploadInput.Parts = parts
	sample.doCompleteMultipartUpload(completeMultipartUploadInput)
}

func (sample ConcurrentCopyPartSample) doCompleteMultipartUpload(input *wos.CompleteMultipartUploadInput) {
	_, err := sample.wosClient.CompleteMultipartUpload(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Complete multiparts finished")
}

func RunConcurrentCopyPartSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
		objectKey  = "object-test"
	)

	sample := newConcurrentCopyPartSample(ak, sk, endpoint, bucketName, objectKey)

	sampleFilePath := filePathSample
	//60MB file
	sample.createSampleFile(sampleFilePath, 1024*1024*60)
	//Upload an object to your source bucket
	sample.PutFile(sampleFilePath)

	sample.DoConcurrentCopyPart()
}
