/**
 * This sample demonstrates how to download an object
 * from WOS in different ways using the WOS SDK for Go.
 */
package examples

import (
	"fmt"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk-v2/wos"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type DownloadSample struct {
	bucketName string
	objectKey  string
	wosClient  *wos.WosClient
}

func newDownloadSample(ak, sk, endpoint, bucketName, objectKey string) *DownloadSample {
	wosClient, err := wos.New(ak, sk, endpoint, wos.WithSignature(wos.SignatureWos), wos.WithMaxRetryCount(0))
	if err != nil {
		panic(err)
	}
	return &DownloadSample{wosClient: wosClient, bucketName: bucketName, objectKey: objectKey}
}

func (sample DownloadSample) PutObject() {
	input := &wos.PutObjectInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.Body = strings.NewReader("Hello WOS")

	_, err := sample.wosClient.PutObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Put object:%s successfully!\n", sample.objectKey)
	fmt.Println()
}

func (sample DownloadSample) GetObject() {
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

func (sample DownloadSample) PutFile(sampleFilePath string) {
	input := &wos.PutFileInput{}
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.SourceFile = sampleFilePath

	_, err := sample.wosClient.PutFile(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Put object:%s with file:%s successfully!\n", sample.objectKey, sampleFilePath)
	fmt.Println()
}

func (sample DownloadSample) DeleteObject() {
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

func (DownloadSample) createSampleFile(sampleFilePath string) {
	if err := os.MkdirAll(filepath.Dir(sampleFilePath), os.ModePerm); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(sampleFilePath, []byte("Hello WOS from file"), os.ModePerm); err != nil {
		panic(err)
	}
}

func RunDownloadSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
		objectKey  = "object-test"
		location   = "yourbucketlocation"
	)
	sample := newDownloadSample(ak, sk, endpoint, bucketName, objectKey)

	fmt.Println("Uploading a new object to WOS from string")
	sample.PutObject()

	fmt.Println("Download object to string")
	sample.GetObject()

	fmt.Println("Uploading a new object to WOS from file")
	sampleFilePath := "/temp/text.txt"
	sample.createSampleFile(sampleFilePath)
	defer func() {
		errMsg := os.Remove(sampleFilePath)
		if errMsg != nil {
			panic(errMsg)
		}
	}()
	sample.PutFile(sampleFilePath)

	fmt.Println("Download file to string")
	sample.GetObject()

	sample.DeleteObject()
}
