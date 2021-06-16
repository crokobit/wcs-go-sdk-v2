/**
 * This sample demonstrates how to create an empty folder under
 * specified bucket to WOS using the WOS SDK for Go.
 */
package examples

import (
	"fmt"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk-v2/wos"
)

type CreateFolderSample struct {
	bucketName string
	wosClient  *wos.WosClient
}

func newCreateFolderSample(ak, sk, endpoint, bucketName string) *CreateFolderSample {
	wosClient, err := wos.New(ak, sk, endpoint, wos.WithSignature(wos.SignatureWos), wos.WithMaxRetryCount(0))
	if err != nil {
		panic(err)
	}
	return &CreateFolderSample{wosClient: wosClient, bucketName: bucketName}
}

func RunCreateFolderSample() {
	const (
		endpoint   = "https://your-endpoint"
		ak         = "*** Provide your Access Key ***"
		sk         = "*** Provide your Secret Key ***"
		bucketName = "bucket-test"
	)
	sample := newCreateFolderSample(ak, sk, endpoint, bucketName)

	keySuffixWithSlash1 := "MyObjectKey1/"
	keySuffixWithSlash2 := "MyObjectKey2/"

	// Create two empty folder without request body, note that the key must be suffixed with a slash
	var input = &wos.PutObjectInput{}
	input.Bucket = bucketName
	input.Key = keySuffixWithSlash1

	_, err := sample.wosClient.PutObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Create empty folder:%s successfully!\n", keySuffixWithSlash1)
	fmt.Println()

	input.Key = keySuffixWithSlash2
	_, err = sample.wosClient.PutObject(input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Create empty folder:%s successfully!\n", keySuffixWithSlash2)
	fmt.Println()

	// Verify whether the size of the empty folder is zero
	var input2 = &wos.GetObjectMetadataInput{}
	input2.Bucket = bucketName
	input2.Key = keySuffixWithSlash1
	output, err := sample.wosClient.GetObjectMetadata(input2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Size of the empty folder %s is %d \n", keySuffixWithSlash1, output.ContentLength)
	fmt.Println()

	input2.Key = keySuffixWithSlash2
	output, err = sample.wosClient.GetObjectMetadata(input2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Size of the empty folder %s is %d \n", keySuffixWithSlash2, output.ContentLength)
	fmt.Println()

}
