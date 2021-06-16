/**
 * This sample demonstrates how to multipart upload an object concurrently
 * from WOS using the WOS SDK for Go.
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

type ConcurrentUploadPartSample struct {
	bucketName string
	objectKey  string
	wosClient  *wos.WosClient
}

func newConcurrentUploadPartSample(ak, sk, endpoint, bucketName, objectKey string) *ConcurrentUploadPartSample {
	wosClient, err := wos.New(ak, sk, endpoint, wos.WithSignature(wos.SignatureWos), wos.WithMaxRetryCount(0))
	if err != nil {
		panic(err)
	}
	return &ConcurrentUploadPartSample{wosClient: wosClient, bucketName: bucketName, objectKey: objectKey}
}

func (sample ConcurrentUploadPartSample) checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func (sample ConcurrentUploadPartSample) createSampleFile(sampleFilePath string, byteCount int64) {
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

func (sample ConcurrentUploadPartSample) PutFile(sampleFilePath string) {
	input := &wos.PutFileInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.SourceFile = sampleFilePath
	_, err := sample.wosClient.PutFile(input)
	if err != nil {
		panic(err)
	}
}

func (sample ConcurrentUploadPartSample) DoConcurrentUploadPart(sampleFilePath string) {
	// Claim a upload id firstly
	input := &wos.InitiateMultipartUploadInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	output, err := sample.wosClient.InitiateMultipartUpload(input)
	if err != nil {
		panic(err)
	}
	uploadId := output.UploadId

	fmt.Printf("Claiming a new upload id %s\n", uploadId)
	fmt.Println()

	// Calculate how many blocks to be divided
	// 5MB
	var partSize int64 = 5 * 1024 * 1024

	stat, err := os.Stat(sampleFilePath)
	if err != nil {
		panic(err)
	}
	fileSize := stat.Size()

	partCount := int(fileSize / partSize)

	if fileSize%partSize != 0 {
		partCount++
	}
	fmt.Printf("Total parts count %d\n", partCount)
	fmt.Println()

	//  Upload parts
	fmt.Println("Begin to upload parts to WOS")

	partChan := make(chan wos.Part, 5)

	for i := 0; i < partCount; i++ {
		partNumber := i + 1
		offset := int64(i) * partSize
		currPartSize := partSize
		if i+1 == partCount {
			currPartSize = fileSize - offset
		}
		go func(index int, offset, partSize int64) {
			uploadPartInput := &wos.UploadPartInput{}
			uploadPartInput.Bucket = sample.bucketName
			uploadPartInput.Key = sample.objectKey
			uploadPartInput.UploadId = uploadId
			uploadPartInput.SourceFile = sampleFilePath
			uploadPartInput.PartNumber = index
			uploadPartInput.Offset = offset
			uploadPartInput.PartSize = partSize
			uploadPartInputOutput, errMsg := sample.wosClient.UploadPart(uploadPartInput)
			if errMsg == nil {
				fmt.Printf("%d finished\n", index)
				partChan <- wos.Part{ETag: uploadPartInputOutput.ETag, PartNumber: uploadPartInputOutput.PartNumber}
			} else {
				panic(errMsg)
			}
		}(partNumber, offset, currPartSize)
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
	completeMultipartUploadInput.Bucket = sample.bucketName
	completeMultipartUploadInput.Key = sample.objectKey
	completeMultipartUploadInput.UploadId = uploadId
	completeMultipartUploadInput.Parts = parts
	sample.doCompleteMultipartUpload(completeMultipartUploadInput)
}

func (sample ConcurrentUploadPartSample) doCompleteMultipartUpload(input *wos.CompleteMultipartUploadInput) {
	_, err := sample.wosClient.CompleteMultipartUpload(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Complete multiparts finished")
}

func RunConcurrentUploadPartSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
		objectKey  = "object-test"
		location   = "yourbucketlocation"
	)

	sample := newConcurrentUploadPartSample(ak, sk, endpoint, bucketName, objectKey)

	//60MB file
	sampleFilePath := "/temp/uploadText.txt"
	sample.createSampleFile(sampleFilePath, 1024*1024*60)

	sample.DoConcurrentUploadPart(sampleFilePath)
}
