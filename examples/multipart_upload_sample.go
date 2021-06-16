/**
 * This sample demonstrates how to download an object concurrently
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

type MultipartUploadSample struct {
	bucketName string
	objectKey  string
	wosClient  *wos.WosClient
}

type CurrentPartUploadResult struct {
	part    wos.Part
	success bool
}

func newMultipartUploadSample(ak, sk, endpoint, bucketName, objectKey string) *MultipartUploadSample {
	wosClient, err := wos.New(ak, sk, endpoint, wos.WithSignature(wos.SignatureWos), wos.WithMaxRetryCount(0))
	if err != nil {
		panic(err)
	}
	return &MultipartUploadSample{wosClient: wosClient, bucketName: bucketName, objectKey: objectKey}
}

func (sample MultipartUploadSample) checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func (sample MultipartUploadSample) createSampleFile(sampleFilePath string, byteCount int64) {
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

func (sample MultipartUploadSample) DoMultipartUploadFile(sampleFilePath string) {
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

	jobs := make(chan *wos.UploadPartInput, 100)
	partChan := make(chan *CurrentPartUploadResult, 10000)

	stopFlag := false
	poolSize := 5
	for i := 0; i < poolSize; i++ {
		go sample.worker(jobs, partChan, &stopFlag)
	}

	go func() {
		for i := 0; i < partCount; i++ {
			partNumber := i + 1
			offset := int64(i) * partSize
			currPartSize := partSize
			if i+1 == partCount {
				currPartSize = fileSize - offset
			}

			uploadPartInput := &wos.UploadPartInput{}
			uploadPartInput.Bucket = sample.bucketName
			uploadPartInput.Key = sample.objectKey
			uploadPartInput.UploadId = uploadId
			uploadPartInput.SourceFile = sampleFilePath
			uploadPartInput.PartNumber = partNumber
			uploadPartInput.Offset = offset
			uploadPartInput.PartSize = currPartSize

			jobs <- uploadPartInput
		}
		close(jobs)
	}()

	parts := make([]wos.Part, 0, partCount)

	for i := 0; i < partCount; i++ {
		result, ok := <-partChan
		if !ok {
			break
		}
		if result.success {
			parts = append(parts, result.part)
		}
	}
	close(partChan)

	if len(parts) == partCount {
		fmt.Println()
		fmt.Println("Completing to upload multiparts")
		completeMultipartUploadInput := &wos.CompleteMultipartUploadInput{}
		completeMultipartUploadInput.Bucket = sample.bucketName
		completeMultipartUploadInput.Key = sample.objectKey
		completeMultipartUploadInput.UploadId = uploadId
		completeMultipartUploadInput.Parts = parts
		sample.doCompleteMultipartUpload(completeMultipartUploadInput)
	} else {
		fmt.Println()
		fmt.Println("upload parts fail.")
	}

}

func (sample MultipartUploadSample) doCompleteMultipartUpload(input *wos.CompleteMultipartUploadInput) {
	_, err := sample.wosClient.CompleteMultipartUpload(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Complete multiparts finished")
}

func (sample MultipartUploadSample) worker(jobs <-chan *wos.UploadPartInput, partChan chan<- *CurrentPartUploadResult, stopFlag *bool) {
	for uploadPartInput := range jobs {
		if *stopFlag {
			partChan <- &CurrentPartUploadResult{part: wos.Part{ETag: "", PartNumber: uploadPartInput.PartNumber}, success: false}
			continue
		}
		uploadPartInputOutput, errMsg := sample.wosClient.UploadPart(uploadPartInput)
		if errMsg == nil {
			fmt.Printf("uploadPart %d success\n", uploadPartInput.PartNumber)
			partChan <- &CurrentPartUploadResult{part: wos.Part{ETag: uploadPartInputOutput.ETag, PartNumber: uploadPartInput.PartNumber}, success: true}
		} else {
			fmt.Printf("uploadPart %d fail\n", uploadPartInput.PartNumber)
			partChan <- &CurrentPartUploadResult{part: wos.Part{ETag: "", PartNumber: uploadPartInput.PartNumber}, success: false}
			*stopFlag = true
			continue
		}
	}
}

func RunMultipartUploadSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
		objectKey  = "object-test"
	)

	sample := newMultipartUploadSample(ak, sk, endpoint, bucketName, objectKey)

	//60MB file
	sampleFilePath := "/temp/uploadText.txt"
	sample.createSampleFile(sampleFilePath, 1024*1024*60)

	sample.DoMultipartUploadFile(sampleFilePath)
}
