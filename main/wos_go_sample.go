package main

import (
	"../examples"
	"fmt"
	"github.com/Wangsu-Cloud-Storage/wcs-go-sdk-v2/wos"
	"io/ioutil"
	"os"
	"strings"
)

const (
	endpoint   = "https://your-endpoint"
	ak         = "*** Provide your Access Key ***"
	sk         = "*** Provide your Secret Key ***"
	bucketName = "bucket-test"
	objectKey  = "object-test"
)

var wosClient *wos.WosClient

func getWosClient() *wos.WosClient {
	var err error
	if wosClient == nil {
		wosClient, err = wos.New(ak, sk, endpoint)
		if err != nil {
			panic(err)
		}
	}
	return wosClient
}

func listBuckets() {
	input := &wos.ListBucketsInput{}
	input.QueryLocation = true
	output, err := getWosClient().ListBuckets(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("Owner.DisplayName:%s, Owner.ID:%s\n", output.Owner.DisplayName, output.Owner.ID)
		for index, val := range output.Buckets {
			fmt.Printf("Bucket[%d]-Name:%s,CreationDate:%s,EndPoint:%s,Region:%s\n", index, val.Name, val.CreationDate, val.Endpoint, val.Region)
		}
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func listObjects() {
	input := &wos.ListObjectsInput{}
	input.Bucket = bucketName
	//	input.Prefix = "src/"
	output, err := getWosClient().ListObjects(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s,OwnerId:%s,OwnerName:%s,\n", output.StatusCode, output.RequestId, output.Owner.ID, output.Owner.DisplayName)
		for index, val := range output.Contents {
			fmt.Printf("Content[%d]-ETag:%s, Key:%s, LastModified:%s, Size:%d, StorageClass:%s\n",
				index, val.ETag, val.Key, val.LastModified, val.Size, val.StorageClass)
		}
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func listObjectsV2() {
	input := &wos.ListObjectsV2Input{}
	input.Bucket = bucketName
	output, err := getWosClient().ListObjectV2(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s,OwnerId:%s,OwnerName:%s,\n", output.StatusCode, output.RequestId, output.Owner.ID, output.Owner.DisplayName)
		for index, val := range output.Contents {
			fmt.Printf("Content[%d]-ETag:%s, Key:%s, LastModified:%s, Size:%d, StorageClass:%s\n",
				index, val.ETag, val.Key, val.LastModified, val.Size, val.StorageClass)
		}
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func headBucket() {
	output, err := getWosClient().HeadBucket(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
		} else {
			fmt.Println(err)
		}
	}
}

func setBucketLifecycleConfiguration() {
	input := &wos.SetBucketLifecycleConfigurationInput{}
	input.Bucket = bucketName

	var lifecycleRules [1]wos.LifecycleRule
	lifecycleRule0 := wos.LifecycleRule{}
	lifecycleRule0.ID = "rule0"
	lifecycleRule0.Status = wos.RuleStatusEnabled

	var filter wos.Filter
	filter.Prefix = "prefix0"

	var transitions [2]wos.Transition
	transitions[0] = wos.Transition{}
	transitions[0].Days = 30
	transitions[0].StorageClass = wos.StorageClassIA

	transitions[1] = wos.Transition{}
	transitions[1].Days = 90
	transitions[1].StorageClass = wos.StorageClassArchive
	lifecycleRule0.Transitions = transitions[:]

	lifecycleRule0.Expiration.Days = 100
	lifecycleRule0.Filter = filter

	lifecycleRules[0] = lifecycleRule0

	input.LifecycleRules = lifecycleRules[:]

	output, err := getWosClient().SetBucketLifecycleConfiguration(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getBucketLifecycleConfiguration() {
	output, err := getWosClient().GetBucketLifecycleConfiguration(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		for index, lifecycleRule := range output.LifecycleRules {
			fmt.Printf("LifecycleRule[%d]:\n", index)
			fmt.Printf("ID:%s, Prefix:%s, Status:%s\n", lifecycleRule.ID, lifecycleRule.Filter.Prefix, lifecycleRule.Status)

			for _, transition := range lifecycleRule.Transitions {
				fmt.Printf("transition.StorageClass:%s, Transition.Days:%d\n", transition.StorageClass, transition.Days)
			}

			fmt.Printf("Expiration.Days:%d\n", lifecycleRule.Expiration.Days)
		}
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func deleteBucketLifecycleConfiguration() {
	output, err := getWosClient().DeleteBucketLifecycleConfiguration(bucketName)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func listMultipartUploads() {
	input := &wos.ListMultipartUploadsInput{}
	input.Bucket = bucketName
	input.MaxUploads = 10
	output, err := getWosClient().ListMultipartUploads(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		for index, upload := range output.Uploads {
			fmt.Printf("Upload[%d]-OwnerId:%s, OwnerName:%s, UploadId:%s, Key:%s, Initiated:%s,StorageClass:%s\n",
				index, upload.Owner.ID, upload.Owner.DisplayName, upload.UploadId, upload.Key, upload.Initiated, upload.StorageClass)
		}
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func deleteObject() {
	input := &wos.DeleteObjectInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	output, err := getWosClient().DeleteObject(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func deleteObjects() {
	input := &wos.DeleteObjectsInput{}
	input.Bucket = bucketName
	var objects [3]wos.ObjectToDelete
	objects[0] = wos.ObjectToDelete{Key: "key1"}
	objects[1] = wos.ObjectToDelete{Key: "key2"}
	objects[2] = wos.ObjectToDelete{Key: "key3"}

	input.Objects = objects[:]
	output, err := getWosClient().DeleteObjects(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		for index, deleted := range output.Deleteds {
			fmt.Printf("Deleted[%d]-Key:%s\n", index, deleted.Key)
		}
		for index, err := range output.Errors {
			fmt.Printf("Error[%d]-Key:%s, Code:%s\n", index, err.Key, err.Code)
		}
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func restoreObject() {
	input := &wos.RestoreObjectInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	input.Days = 1
	output, err := getWosClient().RestoreObject(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getObjectMetadata() {
	input := &wos.GetObjectMetadataInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	output, err := getWosClient().GetObjectMetadata(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("StorageClass:%s, ETag:%s, ContentType:%s, ContentLength:%d, LastModified:%s\n",
			output.StorageClass, output.ETag, output.ContentType, output.ContentLength, output.LastModified)
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Printf("StatusCode:%d\n", wosError.StatusCode)
		} else {
			fmt.Println(err)
		}
	}
}

func copyObject() {
	input := &wos.CopyObjectInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	input.CopySourceBucket = bucketName
	input.CopySourceKey = objectKey + "-back"
	input.Metadata = map[string]string{"meta": "value"}
	input.MetadataDirective = wos.ReplaceMetadata

	output, err := getWosClient().CopyObject(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("ETag:%s, LastModified:%s\n",
			output.ETag, output.LastModified)
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func initiateMultipartUpload() {
	input := &wos.InitiateMultipartUploadInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	input.Metadata = map[string]string{"meta": "value"}
	output, err := getWosClient().InitiateMultipartUpload(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("Bucket:%s, Key:%s, UploadId:%s\n", output.Bucket, output.Key, output.UploadId)
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func abortMultipartUpload() {
	input := &wos.ListMultipartUploadsInput{}
	input.Bucket = bucketName
	output, err := getWosClient().ListMultipartUploads(input)
	if err == nil {
		for _, upload := range output.Uploads {
			input := &wos.AbortMultipartUploadInput{Bucket: bucketName}
			input.UploadId = upload.UploadId
			input.Key = upload.Key
			output, err := getWosClient().AbortMultipartUpload(input)
			if err == nil {
				fmt.Printf("Abort uploadId[%s] successfully\n", input.UploadId)
				fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
			}
		}
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func putObject() {
	input := &wos.PutObjectInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	input.Metadata = map[string]string{"meta": "value"}
	input.Body = strings.NewReader("Hello WOS")
	output, err := getWosClient().PutObject(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("ETag:%s", output.ETag)
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func putFile() {
	input := &wos.PutFileInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	input.SourceFile = "localfile"
	output, err := getWosClient().PutFile(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("ETag:%s\n", output.ETag)
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func headObject() {
	input := &wos.HeadObjectInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	output, err := getWosClient().HeadObject(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func uploadPart() {
	sourceFile := "localfile"
	var partSize int64 = 1024 * 1024 * 5
	fileInfo, statErr := os.Stat(sourceFile)
	if statErr != nil {
		panic(statErr)
	}
	partCount := fileInfo.Size() / partSize
	if fileInfo.Size()%partSize > 0 {
		partCount++
	}
	var i int64
	for i = 0; i < partCount; i++ {
		input := &wos.UploadPartInput{}
		input.Bucket = bucketName
		input.Key = objectKey
		input.UploadId = "uploadid"
		input.PartNumber = int(i + 1)
		input.Offset = i * partSize
		if i == partCount-1 {
			input.PartSize = fileInfo.Size()
		} else {
			input.PartSize = partSize
		}
		input.SourceFile = sourceFile
		output, err := getWosClient().UploadPart(input)
		if err == nil {
			fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
			fmt.Printf("ETag:%s\n", output.ETag)
		} else {
			if wosError, ok := err.(wos.WosError); ok {
				fmt.Println(wosError.StatusCode)
				fmt.Println(wosError.Code)
				fmt.Println(wosError.Message)
			} else {
				fmt.Println(err)
			}
		}
	}
}

func listParts() {
	input := &wos.ListPartsInput{}
	input.Bucket = bucketName
	input.Key = "100m.txt"
	input.UploadId = "27f44debc56948f49c0e3d357a186c89"
	output, err := getWosClient().ListParts(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		for index, part := range output.Parts {
			fmt.Printf("Part[%d]-ETag:%s, PartNumber:%d\n", index, part.ETag,
				part.PartNumber)
		}
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func completeMultipartUpload() {
	input := &wos.CompleteMultipartUploadInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	input.UploadId = "uploadid"
	input.Parts = []wos.Part{
		wos.Part{PartNumber: 1, ETag: "etag1"},
		wos.Part{PartNumber: 2, ETag: "etag2"},
		wos.Part{PartNumber: 3, ETag: "etag3"},
	}
	output, err := getWosClient().CompleteMultipartUpload(input)
	if err == nil {
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("Location:%s, Bucket:%s, Key:%s, ETag:%s\n", output.Location, output.Bucket, output.Key, output.ETag)
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func copyPart() {

	sourceBucket := "source-bucket"
	sourceKey := "source-key"
	input := &wos.GetObjectMetadataInput{}
	input.Bucket = sourceBucket
	input.Key = sourceKey
	output, err := getWosClient().GetObjectMetadata(input)
	if err == nil {
		objectSize := output.ContentLength
		var partSize int64 = 5 * 1024 * 1024
		partCount := objectSize / partSize
		if objectSize%partSize > 0 {
			partCount++
		}
		var i int64
		for i = 0; i < partCount; i++ {
			input := &wos.CopyPartInput{}
			input.Bucket = bucketName
			input.Key = objectKey
			input.UploadId = "uploadid"
			input.PartNumber = int(i + 1)
			input.CopySourceBucket = sourceBucket
			input.CopySourceKey = sourceKey
			input.CopySourceRangeStart = i * partSize
			if i == partCount-1 {
				input.CopySourceRangeEnd = objectSize - 1
			} else {
				input.CopySourceRangeEnd = (i+1)*partSize - 1
			}
			output, err := getWosClient().CopyPart(input)
			if err == nil {
				fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
				fmt.Printf("ETag:%s, PartNumber:%d\n", output.ETag, output.PartNumber)
			} else {
				if wosError, ok := err.(wos.WosError); ok {
					fmt.Println(wosError.StatusCode)
					fmt.Println(wosError.Code)
					fmt.Println(wosError.Message)
				} else {
					fmt.Println(err)
				}
			}
		}
	}
}

func getObject() {
	input := &wos.GetObjectInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	output, err := getWosClient().GetObject(input)
	if err == nil {
		defer output.Body.Close()
		fmt.Printf("StatusCode:%d, RequestId:%s\n", output.StatusCode, output.RequestId)
		fmt.Printf("StorageClass:%s, ETag:%s, ContentType:%s, ContentLength:%d, LastModified:%s\n",
			output.StorageClass, output.ETag, output.ContentType, output.ContentLength, output.LastModified)
		p := make([]byte, 1024)
		var readErr error
		var readCount int
		for {
			readCount, readErr = output.Body.Read(p)
			if readCount > 0 {
				fmt.Printf("%s", p[:readCount])
			}
			if readErr != nil {
				break
			}
		}
	} else {
		if wosError, ok := err.(wos.WosError); ok {
			fmt.Println(wosError.StatusCode)
			fmt.Println(wosError.Code)
			fmt.Println(wosError.Message)
		} else {
			fmt.Println(err)
		}
	}
}

func getAvinfo() {
	input := &wos.GetAvinfoInput{}
	input.Bucket = bucketName
	input.Key = objectKey
	getAvinfoOutput, err := getWosClient().GetAvinfo(input)
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

func runExamples() {
	examples.RunObjectOperationsSample()
	examples.RunDownloadSample()
	examples.RunCreateFolderSample()
	examples.RunDeleteObjectsSample()
	examples.RunListObjectsSample()
	examples.RunListObjectsInFolderSample()
	examples.RunConcurrentCopyPartSample()
	examples.RunConcurrentDownloadObjectSample()
	examples.RunConcurrentUploadPartSample()
	examples.RunSimpleMultipartUploadSample()
	examples.RunObjectMetaSample()
	examples.RunTemporarySignatureSample()
}

func main() {
	//---- init log ----
	defer wos.CloseLog()
	wos.InitLog("/temp/WOS-SDK.log", 1024*1024*100, 5, wos.LEVEL_WARN, false)

	//---- run examples----
	//	runExamples()

	//---- bucket related APIs ----

	// listBuckets()
	// wos.FlushLog()
	// listObjects()
	// listObjectsV2()
	// listMultipartUploads()
	// headBucket()
	// setBucketLifecycleConfiguration()
	// getBucketLifecycleConfiguration()
	// deleteBucketLifecycleConfiguration()

	//---- object related APIs ----
	// deleteObject()
	// deleteObjects()
	// restoreObject()
	// copyObject()
	// initiateMultipartUpload()
	// uploadPart()
	// copyPart()
	// listParts()
	// completeMultipartUpload()
	// abortMultipartUpload()
	// putObject()
	// putFile()
	// headObject()
	// getObjectMetadata()
	// getObject()
	// getAvinfo()
}
