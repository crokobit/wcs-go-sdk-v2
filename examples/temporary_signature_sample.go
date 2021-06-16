/**
 * This sample demonstrates how to do common operations in temporary signature way
 * on WOS using the WOS SDK for Go.
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

type TemporarySignatureSample struct {
	bucketName string
	objectKey  string
	wosClient  *wos.WosClient
}

func newTemporarySignatureSample(ak, sk, endpoint, bucketName, objectKey string) *TemporarySignatureSample {
	wosClient, err := wos.New(ak, sk, endpoint)
	if err != nil {
		panic(err)
	}
	return &TemporarySignatureSample{wosClient: wosClient, bucketName: bucketName, objectKey: objectKey}
}

func (sample TemporarySignatureSample) ListBuckets() {
	input := &wos.CreateSignedUrlInput{}
	input.Method = wos.HttpMethodGet
	input.Expires = 3600
	output, err := sample.wosClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "ListBuckets")
	fmt.Println(output.SignedUrl)

	listBucketsOutput, err := sample.wosClient.ListBucketsWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Owner.DisplayName:%s, Owner.ID:%s\n", listBucketsOutput.Owner.DisplayName, listBucketsOutput.Owner.ID)
	for index, val := range listBucketsOutput.Buckets {
		fmt.Printf("Bucket[%d]-Name:%s,CreationDate:%s\n", index, val.Name, val.CreationDate)
	}
	fmt.Println()
}

func (sample TemporarySignatureSample) PutObject() {
	input := &wos.CreateSignedUrlInput{}
	input.Method = wos.HttpMethodPut
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.Expires = 3600
	output, err := sample.wosClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "PutObject")
	fmt.Println(output.SignedUrl)

	data := strings.NewReader("Hello WOS")
	_, err = sample.wosClient.PutObjectWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders, data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Put object:%s successfully!\n", sample.objectKey)
	fmt.Println()
}

func (TemporarySignatureSample) createSampleFile(sampleFilePath string) {
	if err := os.MkdirAll(filepath.Dir(sampleFilePath), os.ModePerm); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(sampleFilePath, []byte("Hello WOS from file"), os.ModePerm); err != nil {
		panic(err)
	}
}

func (sample TemporarySignatureSample) PutFile(sampleFilePath string) {
	input := &wos.CreateSignedUrlInput{}
	input.Method = wos.HttpMethodPut
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.Expires = 3600
	output, err := sample.wosClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "PutFile")
	fmt.Println(output.SignedUrl)

	_, err = sample.wosClient.PutFileWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders, sampleFilePath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Put file:%s successfully!\n", sample.objectKey)
	fmt.Println()
}

func (sample TemporarySignatureSample) GetObject() {
	input := &wos.CreateSignedUrlInput{}
	input.Method = wos.HttpMethodGet
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.Expires = 3600
	output, err := sample.wosClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "GetObject")
	fmt.Println(output.SignedUrl)

	getObjectOutput, err := sample.wosClient.GetObjectWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders)
	if err != nil {
		panic(err)
	}
	defer func() {
		errMsg := getObjectOutput.Body.Close()
		if errMsg != nil {
			panic(errMsg)
		}
	}()
	fmt.Println("Object content:")
	body, err := ioutil.ReadAll(getObjectOutput.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	fmt.Println()
}

func (sample TemporarySignatureSample) GetAvinfo() {
	input := &wos.CreateSignedUrlInput{}
	input.Method = wos.HttpMethodGet
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.SubResource = wos.SubResourceAvinfo
	output, err := sample.wosClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "GetAvinfo")
	fmt.Println(output.SignedUrl)

	getAvinfoOutput, err := sample.wosClient.GetAvinfoWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders)
	if err != nil {
		panic(err)
	}
	defer func() {
		errMsg := getAvinfoOutput.Body.Close()
		if errMsg != nil {
			panic(errMsg)
		}
	}()
	fmt.Println("Object content:")
	body, err := ioutil.ReadAll(getAvinfoOutput.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	fmt.Println()
}

func (sample TemporarySignatureSample) DeleteObject() {
	input := &wos.CreateSignedUrlInput{}
	input.Method = wos.HttpMethodDelete
	input.Bucket = sample.bucketName
	input.Key = sample.objectKey
	input.Expires = 3600
	output, err := sample.wosClient.CreateSignedUrl(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s using temporary signature url:\n", "DeleteObject")
	fmt.Println(output.SignedUrl)

	_, err = sample.wosClient.DeleteObjectWithSignedUrl(output.SignedUrl, output.ActualSignedRequestHeaders)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Delete object:%s successfully!\n", sample.objectKey)
	fmt.Println()
}

func RunTemporarySignatureSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
		objectKey  = "object-test"
	)

	sample := newTemporarySignatureSample(ak, sk, endpoint, bucketName, objectKey)

	// List buckets
	sample.ListBuckets()

	// Put object
	sample.PutObject()

	// Get object
	sample.GetObject()

	// Put file
	sampleFilePath := "/temp/sampleText.txt"
	sample.createSampleFile(sampleFilePath)

	sample.PutFile(sampleFilePath)
	// Get object
	sample.GetObject()

	// Get avinfo
	sample.GetAvinfo()

	// Delete object
	sample.DeleteObject()
}
