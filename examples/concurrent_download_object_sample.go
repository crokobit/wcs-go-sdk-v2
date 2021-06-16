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
	"sync"
	"time"
)

type ConcurrentDownloadObjectSample struct {
	bucketName string
	objectKey  string
	wosClient  *wos.WosClient
}

func newConcurrentDownloadObjectSample(ak, sk, endpoint, bucketName, objectKey string) *ConcurrentDownloadObjectSample {
	wosClient, err := wos.New(ak, sk, endpoint, wos.WithPathStyle(true))
	if err != nil {
		panic(err)
	}
	return &ConcurrentDownloadObjectSample{wosClient: wosClient, bucketName: bucketName, objectKey: objectKey}
}

func (sample ConcurrentDownloadObjectSample) createSampleFile(sampleFilePath string, byteCount int64) {
	if err := os.MkdirAll(filepath.Dir(sampleFilePath), os.ModePerm); err != nil {
		panic(err)
	}

	fd, err := os.OpenFile(sampleFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
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

func (sample ConcurrentDownloadObjectSample) PutFile(sampleFilePath string) {
	input := &wos.PutFileInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.SourceFile = sampleFilePath
	_, err := sample.wosClient.PutFile(input)
	if err != nil {
		panic(err)
	}
}

func (sample ConcurrentDownloadObjectSample) checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func (sample ConcurrentDownloadObjectSample) DoConcurrentDownload(sampleFilePath string) {

	// Get size of the object
	getObjectMetadataInput := &wos.GetObjectMetadataInput{}
	getObjectMetadataInput.Bucket = sample.bucketName
	getObjectMetadataInput.Key = sample.objectKey
	getObjectMetadataOutput, err := sample.wosClient.GetObjectMetadata(getObjectMetadataInput)
	sample.checkError(err)

	objectSize := getObjectMetadataOutput.ContentLength

	// Calculate how many blocks to be divided
	// 5MB
	var partSize int64 = 1024 * 1024 * 5
	partCount := int(objectSize / partSize)

	if objectSize%partSize != 0 {
		partCount++
	}

	fmt.Printf("Total parts count %d\n", partCount)
	fmt.Println()

	downloadFilePath := filepath.Dir(sampleFilePath) + "/" + sample.objectKey

	var wg sync.WaitGroup
	wg.Add(partCount)

	fd, err := os.OpenFile(downloadFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		panic(errors.New("open file with error"))
	}

	err = fd.Close()
	sample.checkError(err)

	//Download the object concurrently
	fmt.Printf("Start to download %s \n", sample.objectKey)

	for i := 0; i < partCount; i++ {
		index := i + 1
		rangeStart := int64(i) * partSize
		rangeEnd := rangeStart + partSize - 1
		if index == partCount {
			rangeEnd = objectSize - 1
		}
		go func(start, end int64, num int) {
			defer wg.Done()
			getObjectInput := &wos.GetObjectInput{}
			getObjectInput.Bucket = sample.bucketName
			getObjectInput.Key = sample.objectKey
			getObjectInput.RangeStart = start
			getObjectInput.RangeEnd = end
			getObjectOutput, err := sample.wosClient.GetObject(getObjectInput)
			if err == nil {
				defer func() {
					errMsg := getObjectOutput.Body.Close()
					sample.checkError(errMsg)
				}()
				wfd, err := os.OpenFile(downloadFilePath, os.O_WRONLY, 0600)
				sample.checkError(err)
				b := make([]byte, 1024)
				for {
					n, err := getObjectOutput.Body.Read(b)
					if n > 0 {
						wcnt, err := wfd.WriteAt(b[0:n], start)
						sample.checkError(err)
						if n != wcnt {
							panic(fmt.Sprintf("wcnt %d, n %d", wcnt, n))
						}
						start += int64(n)
					}

					if err != nil {
						break
					}
				}
				errMsg := wfd.Sync()
				sample.checkError(errMsg)
				errMsg = wfd.Close()
				sample.checkError(errMsg)
				fmt.Printf("%d finished\n", num)
			} else {
				panic(err)
			}
		}(rangeStart, rangeEnd, index)
	}
	wg.Wait()

	fmt.Printf("Download object finished, downloadPath:%s\n", downloadFilePath)
}

func RunConcurrentDownloadObjectSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
		objectKey  = "object-test"
	)

	sample := newConcurrentDownloadObjectSample(ak, sk, endpoint, bucketName, objectKey)

	//60MB file
	sampleFilePath := "/temp/uploadText.txt"
	sample.createSampleFile(sampleFilePath, 1024*1024*60)
	//Upload an object to your source bucket
	sample.PutFile(sampleFilePath)

	sample.DoConcurrentDownload(sampleFilePath)

}
