/**
 * This sample demonstrates how to list objects under a specified folder of a bucket
 * from WOS using the WOS SDK for Go.
 */
package examples

import (
	"fmt"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk-v2/wos"
	"strconv"
	"strings"
)

type ListObjectsInFolderSample struct {
	bucketName string
	wosClient  *wos.WosClient
}

func newListObjectsInFolderSample(ak, sk, endpoint, bucketName string) *ListObjectsInFolderSample {
	wosClient, err := wos.New(ak, sk, endpoint)
	if err != nil {
		panic(err)
	}
	return &ListObjectsInFolderSample{wosClient: wosClient, bucketName: bucketName}
}

func (sample ListObjectsInFolderSample) prepareObjects(input *wos.PutObjectInput) {
	_, err := sample.wosClient.PutObject(input)
	if err != nil {
		panic(err)
	}
}

func (sample ListObjectsInFolderSample) PrepareFoldersAndObjects() {

	keyPrefix := "MyObjectKeyFolders"
	folderPrefix := "src"
	subFolderPrefix := "test"

	input := &wos.PutObjectInput{}
	input.Bucket = sample.bucketName

	// First prepare folders and sub folders
	for i := 0; i < 5; i++ {
		key := folderPrefix + strconv.Itoa(i) + "/"
		input.Key = key
		sample.prepareObjects(input)
		for j := 0; j < 3; j++ {
			subKey := key + subFolderPrefix + strconv.Itoa(j) + "/"
			input.Key = subKey
			sample.prepareObjects(input)
		}
	}

	// Insert 2 objects in each folder
	input.Body = strings.NewReader("Hello WOS")
	listObjectsInput := &wos.ListObjectsInput{}
	listObjectsInput.Bucket = sample.bucketName
	output, err := sample.wosClient.ListObjects(listObjectsInput)
	if err != nil {
		panic(err)
	}
	for _, content := range output.Contents {
		for i := 0; i < 2; i++ {
			objectKey := content.Key + keyPrefix + strconv.Itoa(i)
			input.Key = objectKey
			sample.prepareObjects(input)
		}
	}

	// Insert 2 objects in root path
	input.Key = keyPrefix + strconv.Itoa(0)
	sample.prepareObjects(input)
	input.Key = keyPrefix + strconv.Itoa(1)
	sample.prepareObjects(input)

	fmt.Println("Prepare folders and objects finished")
	fmt.Println()
}

func (sample ListObjectsInFolderSample) ListObjectsInFolders() {
	fmt.Println("List objects in folder src0/")
	input := &wos.ListObjectsInput{}
	input.Bucket = sample.bucketName
	input.Prefix = "src0/"
	output, err := sample.wosClient.ListObjects(input)
	if err != nil {
		panic(err)
	}
	for index, val := range output.Contents {
		fmt.Printf("Content[%d]-ETag:%s, Key:%s, Size:%d\n",
			index, val.ETag, val.Key, val.Size)
	}

	fmt.Println()

	fmt.Println("List objects in sub folder src0/test0/")

	input.Prefix = "src0/test0/"
	output, err = sample.wosClient.ListObjects(input)
	if err != nil {
		panic(err)
	}
	for index, val := range output.Contents {
		fmt.Printf("Content[%d]-ETag:%s, Key:%s, Size:%d\n",
			index, val.ETag, val.Key, val.Size)
	}

	fmt.Println()
}

func (sample ListObjectsInFolderSample) listObjectsByPrefixes(commonPrefixes []string) {
	input := &wos.ListObjectsInput{}
	input.Bucket = sample.bucketName
	input.Delimiter = "/"
	for _, prefix := range commonPrefixes {
		input.Prefix = prefix
		output, err := sample.wosClient.ListObjects(input)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Folder %s:\n", prefix)
		for index, val := range output.Contents {
			fmt.Printf("Content[%d]-ETag:%s, Key:%s, Size:%d\n",
				index, val.ETag, val.Key, val.Size)
		}
		fmt.Println()
		sample.listObjectsByPrefixes(output.CommonPrefixes)
	}
}

func (sample ListObjectsInFolderSample) ListObjectsGroupByFolder() {
	fmt.Println("List objects group by folder")
	input := &wos.ListObjectsInput{}
	input.Bucket = sample.bucketName
	input.Delimiter = "/"
	output, err := sample.wosClient.ListObjects(input)
	if err != nil {
		panic(err)
	}
	fmt.Println("Root path:")
	for index, val := range output.Contents {
		fmt.Printf("Content[%d]-ETag:%s, Key:%s, Size:%d\n",
			index, val.ETag, val.Key, val.Size)
	}
	fmt.Println()
	sample.listObjectsByPrefixes(output.CommonPrefixes)
}

func (sample ListObjectsInFolderSample) BatchDeleteObjects() {
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
	_, err = sample.wosClient.DeleteObjects(deleteObjectsInput)
	if err != nil {
		panic(err)
	}
	fmt.Println("Delete objects successfully!")
	fmt.Println()
}

func RunListObjectsInFolderSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
	)

	sample := newListObjectsInFolderSample(ak, sk, endpoint, bucketName)

	// First prepare folders and objects
	sample.PrepareFoldersAndObjects()

	// List objects in folders
	sample.ListObjectsInFolders()

	// List all objects group by folder
	sample.ListObjectsGroupByFolder()

	sample.BatchDeleteObjects()
}
