package wos

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// CreateSignedUrl creates signed url with the specified CreateSignedUrlInput, and returns the CreateSignedUrlOutput and error
func (wosClient WosClient) CreateSignedUrl(input *CreateSignedUrlInput) (output *CreateSignedUrlOutput, err error) {
	if input == nil {
		return nil, errors.New("CreateSignedUrlInput is nil")
	}

	params := make(map[string]string, len(input.QueryParams))
	for key, value := range input.QueryParams {
		params[key] = value
	}

	if input.SubResource != "" {
		params[string(input.SubResource)] = ""
	}

	headers := make(map[string][]string, len(input.Headers))
	for key, value := range input.Headers {
		headers[key] = []string{value}
	}

	if input.Expires <= 0 {
		input.Expires = 300
	}

	requestURL, err := wosClient.doAuthTemporary(string(input.Method), input.Bucket, input.Key, params, headers, int64(input.Expires))
	if err != nil {
		return nil, err
	}

	output = &CreateSignedUrlOutput{
		SignedUrl:                  requestURL,
		ActualSignedRequestHeaders: headers,
	}
	return
}

// CreateBrowserBasedSignature gets the browser based signature with the specified CreateBrowserBasedSignatureInput,
// and returns the CreateBrowserBasedSignatureOutput and error
func (wosClient WosClient) CreateBrowserBasedSignature(input *CreateBrowserBasedSignatureInput) (output *CreateBrowserBasedSignatureOutput, err error) {
	if input == nil {
		return nil, errors.New("CreateBrowserBasedSignatureInput is nil")
	}

	params := make(map[string]string, len(input.FormParams))
	for key, value := range input.FormParams {
		params[key] = value
	}

	date := time.Now().UTC()
	shortDate := date.Format(SHORT_DATE_FORMAT)
	longDate := date.Format(LONG_DATE_FORMAT)
	sh := wosClient.getSecurity()

	credential, _ := getCredential(sh.ak, wosClient.conf.region, shortDate, wosClient.conf.signature == SignatureWos)

	if input.Expires <= 0 {
		input.Expires = 300
	}

	expiration := date.Add(time.Second * time.Duration(input.Expires)).Format(ISO8601_DATE_FORMAT)
	if wosClient.conf.signature == SignatureV4 {
		params[PARAM_ALGORITHM_AMZ_CAMEL] = V4_HASH_PREFIX
		params[PARAM_CREDENTIAL_AMZ_CAMEL] = credential
		params[PARAM_DATE_AMZ_CAMEL] = longDate
	}

	matchAnyBucket := true
	matchAnyKey := true
	count := 5
	if bucket := strings.TrimSpace(input.Bucket); bucket != "" {
		params["bucket"] = bucket
		matchAnyBucket = false
		count--
	}

	if key := strings.TrimSpace(input.Key); key != "" {
		params["key"] = key
		matchAnyKey = false
		count--
	}

	originPolicySlice := make([]string, 0, len(params)+count)
	originPolicySlice = append(originPolicySlice, fmt.Sprintf("{\"expiration\":\"%s\",", expiration))
	originPolicySlice = append(originPolicySlice, "\"conditions\":[")
	for key, value := range params {
		if _key := strings.TrimSpace(strings.ToLower(key)); _key != "" {
			originPolicySlice = append(originPolicySlice, fmt.Sprintf("{\"%s\":\"%s\"},", _key, value))
		}
	}

	if matchAnyBucket {
		originPolicySlice = append(originPolicySlice, "[\"starts-with\", \"$bucket\", \"\"],")
	}

	if matchAnyKey {
		originPolicySlice = append(originPolicySlice, "[\"starts-with\", \"$key\", \"\"],")
	}

	originPolicySlice = append(originPolicySlice, "]}")

	originPolicy := strings.Join(originPolicySlice, "")
	policy := Base64Encode([]byte(originPolicy))
	var signature string

	if wosClient.conf.signature == SignatureV2 {
		signature = Base64Encode(HmacSha1([]byte(sh.sk), []byte(policy)))
	} else {
		signature = getSignature(policy, sh.sk, wosClient.conf.region, shortDate, wosClient.conf.signature == SignatureWos)
	}

	output = &CreateBrowserBasedSignatureOutput{
		OriginPolicy: originPolicy,
		Policy:       policy,
		Algorithm:    params[PARAM_ALGORITHM_AMZ_CAMEL],
		Credential:   params[PARAM_CREDENTIAL_AMZ_CAMEL],
		Date:         params[PARAM_DATE_AMZ_CAMEL],
		Signature:    signature,
	}
	return
}

// ListBucketsWithSignedUrl lists buckets with the specified signed url and signed request headers
func (wosClient WosClient) ListBucketsWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *ListBucketsOutput, err error) {
	output = &ListBucketsOutput{}
	err = wosClient.doHTTPWithSignedURL("ListBuckets", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// ListObjectsWithSignedUrl lists objects in a bucket with the specified signed url and signed request headers
func (wosClient WosClient) ListObjectsWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *ListObjectsOutput, err error) {
	output = &ListObjectsOutput{}
	err = wosClient.doHTTPWithSignedURL("ListObjects", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// ListMultipartUploadsWithSignedUrl lists the multipart uploads that are initialized but not combined or aborted in a
// specified bucket with the specified signed url and signed request headers
func (wosClient WosClient) ListMultipartUploadsWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *ListMultipartUploadsOutput, err error) {
	output = &ListMultipartUploadsOutput{}
	err = wosClient.doHTTPWithSignedURL("ListMultipartUploads", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	} else if output.EncodingType == "url" {
		err = decodeListMultipartUploadsOutput(output)
		if err != nil {
			doLog(LEVEL_ERROR, "Failed to get ListMultipartUploadsOutput with error: %v.", err)
			output = nil
		}
	}
	return
}

// HeadBucketWithSignedUrl checks whether a bucket exists with the specified signed url and signed request headers
func (wosClient WosClient) HeadBucketWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = wosClient.doHTTPWithSignedURL("HeadBucket", HTTP_HEAD, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// HeadObjectWithSignedUrl checks whether an object exists with the specified signed url and signed request headers
func (wosClient WosClient) HeadObjectWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = wosClient.doHTTPWithSignedURL("HeadObject", HTTP_HEAD, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// SetBucketLifecycleConfigurationWithSignedUrl sets lifecycle rules for a bucket with the specified signed url and signed request headers and data
func (wosClient WosClient) SetBucketLifecycleConfigurationWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = wosClient.doHTTPWithSignedURL("SetBucketLifecycleConfiguration", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetBucketLifecycleConfigurationWithSignedUrl gets lifecycle rules of a bucket with the specified signed url and signed request headers
func (wosClient WosClient) GetBucketLifecycleConfigurationWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetBucketLifecycleConfigurationOutput, err error) {
	output = &GetBucketLifecycleConfigurationOutput{}
	err = wosClient.doHTTPWithSignedURL("GetBucketLifecycleConfiguration", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// DeleteBucketLifecycleConfigurationWithSignedUrl deletes lifecycle rules of a bucket with the specified signed url and signed request headers
func (wosClient WosClient) DeleteBucketLifecycleConfigurationWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = wosClient.doHTTPWithSignedURL("DeleteBucketLifecycleConfiguration", HTTP_DELETE, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// DeleteObjectWithSignedUrl deletes an object with the specified signed url and signed request headers
func (wosClient WosClient) DeleteObjectWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *DeleteObjectOutput, err error) {
	output = &DeleteObjectOutput{}
	err = wosClient.doHTTPWithSignedURL("DeleteObject", HTTP_DELETE, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// DeleteObjectsWithSignedUrl deletes objects in a batch with the specified signed url and signed request headers and data
func (wosClient WosClient) DeleteObjectsWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *DeleteObjectsOutput, err error) {
	output = &DeleteObjectsOutput{}
	err = wosClient.doHTTPWithSignedURL("DeleteObjects", HTTP_POST, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// RestoreObjectWithSignedUrl restores an object with the specified signed url and signed request headers and data
func (wosClient WosClient) RestoreObjectWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = wosClient.doHTTPWithSignedURL("RestoreObject", HTTP_POST, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetObjectMetadataWithSignedUrl gets object metadata with the specified signed url and signed request headers
func (wosClient WosClient) GetObjectMetadataWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetObjectMetadataOutput, err error) {
	output = &GetObjectMetadataOutput{}
	err = wosClient.doHTTPWithSignedURL("GetObjectMetadata", HTTP_HEAD, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	} else {
		ParseGetObjectMetadataOutput(output)
	}
	return
}

// GetAvinfoWithSignedUrl get object avinfo with the specified signed url and signed request headers
func (wosClient WosClient) GetAvinfoWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetAvinfoOutput, err error) {
	output = &GetAvinfoOutput{}
	err = wosClient.doHTTPWithSignedURL("GetObject", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// GetObjectWithSignedUrl downloads object with the specified signed url and signed request headers
func (wosClient WosClient) GetObjectWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *GetObjectOutput, err error) {
	output = &GetObjectOutput{}
	err = wosClient.doHTTPWithSignedURL("GetObject", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	} else {
		ParseGetObjectOutput(output)
	}
	return
}

// PutObjectWithSignedUrl uploads an object to the specified bucket with the specified signed url and signed request headers and data
func (wosClient WosClient) PutObjectWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *PutObjectOutput, err error) {
	output = &PutObjectOutput{}
	err = wosClient.doHTTPWithSignedURL("PutObject", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	} else {
		ParsePutObjectOutput(output)
	}
	return
}

// PutFileWithSignedUrl uploads a file to the specified bucket with the specified signed url and signed request headers and sourceFile path
func (wosClient WosClient) PutFileWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, sourceFile string) (output *PutObjectOutput, err error) {
	var data io.Reader
	sourceFile = strings.TrimSpace(sourceFile)
	if sourceFile != "" {
		fd, _err := os.Open(sourceFile)
		if _err != nil {
			err = _err
			return nil, err
		}
		defer func() {
			errMsg := fd.Close()
			if errMsg != nil {
				doLog(LEVEL_WARN, "Failed to close file with reason: %v", errMsg)
			}
		}()

		stat, _err := fd.Stat()
		if _err != nil {
			err = _err
			return nil, err
		}
		fileReaderWrapper := &fileReaderWrapper{filePath: sourceFile}
		fileReaderWrapper.reader = fd

		var contentLength int64
		if value, ok := actualSignedRequestHeaders[HEADER_CONTENT_LENGTH_CAMEL]; ok {
			contentLength = StringToInt64(value[0], -1)
		} else if value, ok := actualSignedRequestHeaders[HEADER_CONTENT_LENGTH]; ok {
			contentLength = StringToInt64(value[0], -1)
		} else {
			contentLength = stat.Size()
		}
		if contentLength > stat.Size() {
			return nil, errors.New("ContentLength is larger than fileSize")
		}
		fileReaderWrapper.totalCount = contentLength
		data = fileReaderWrapper
	}

	output = &PutObjectOutput{}
	err = wosClient.doHTTPWithSignedURL("PutObject", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	} else {
		ParsePutObjectOutput(output)
	}
	return
}

// CopyObjectWithSignedUrl creates a copy for an existing object with the specified signed url and signed request headers
func (wosClient WosClient) CopyObjectWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *CopyObjectOutput, err error) {
	output = &CopyObjectOutput{}
	err = wosClient.doHTTPWithSignedURL("CopyObject", HTTP_PUT, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	} else {
		ParseCopyObjectOutput(output)
	}
	return
}

// AbortMultipartUploadWithSignedUrl aborts a multipart upload in a specified bucket by using the multipart upload ID with the specified signed url and signed request headers
func (wosClient WosClient) AbortMultipartUploadWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *BaseModel, err error) {
	output = &BaseModel{}
	err = wosClient.doHTTPWithSignedURL("AbortMultipartUpload", HTTP_DELETE, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	}
	return
}

// InitiateMultipartUploadWithSignedUrl initializes a multipart upload with the specified signed url and signed request headers
func (wosClient WosClient) InitiateMultipartUploadWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *InitiateMultipartUploadOutput, err error) {
	output = &InitiateMultipartUploadOutput{}
	err = wosClient.doHTTPWithSignedURL("InitiateMultipartUpload", HTTP_POST, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	} else {
		ParseInitiateMultipartUploadOutput(output)
		if output.EncodingType == "url" {
			err = decodeInitiateMultipartUploadOutput(output)
			if err != nil {
				doLog(LEVEL_ERROR, "Failed to get InitiateMultipartUploadOutput with error: %v.", err)
				output = nil
			}
		}
	}
	return
}

// UploadPartWithSignedUrl uploads a part to a specified bucket by using a specified multipart upload ID
// with the specified signed url and signed request headers and data
func (wosClient WosClient) UploadPartWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *UploadPartOutput, err error) {
	output = &UploadPartOutput{}
	err = wosClient.doHTTPWithSignedURL("UploadPart", HTTP_PUT, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	} else {
		ParseUploadPartOutput(output)
	}
	return
}

// CompleteMultipartUploadWithSignedUrl combines the uploaded parts in a specified bucket by using the multipart upload ID
// with the specified signed url and signed request headers and data
func (wosClient WosClient) CompleteMultipartUploadWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header, data io.Reader) (output *CompleteMultipartUploadOutput, err error) {
	output = &CompleteMultipartUploadOutput{}
	err = wosClient.doHTTPWithSignedURL("CompleteMultipartUpload", HTTP_POST, signedUrl, actualSignedRequestHeaders, data, output, true)
	if err != nil {
		output = nil
	} else {
		ParseCompleteMultipartUploadOutput(output)
		if output.EncodingType == "url" {
			err = decodeCompleteMultipartUploadOutput(output)
			if err != nil {
				doLog(LEVEL_ERROR, "Failed to get CompleteMultipartUploadOutput with error: %v.", err)
				output = nil
			}
		}
	}
	return
}

// ListPartsWithSignedUrl lists the uploaded parts in a bucket by using the multipart upload ID with the specified signed url and signed request headers
func (wosClient WosClient) ListPartsWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *ListPartsOutput, err error) {
	output = &ListPartsOutput{}
	err = wosClient.doHTTPWithSignedURL("ListParts", HTTP_GET, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	} else if output.EncodingType == "url" {
		err = decodeListPartsOutput(output)
		if err != nil {
			doLog(LEVEL_ERROR, "Failed to get ListPartsOutput with error: %v.", err)
			output = nil
		}
	}
	return
}

// CopyPartWithSignedUrl copy a part to a specified bucket by using a specified multipart upload ID with the specified signed url and signed request headers
func (wosClient WosClient) CopyPartWithSignedUrl(signedUrl string, actualSignedRequestHeaders http.Header) (output *CopyPartOutput, err error) {
	output = &CopyPartOutput{}
	err = wosClient.doHTTPWithSignedURL("CopyPart", HTTP_PUT, signedUrl, actualSignedRequestHeaders, nil, output, true)
	if err != nil {
		output = nil
	} else {
		ParseCopyPartOutput(output)
	}
	return
}
