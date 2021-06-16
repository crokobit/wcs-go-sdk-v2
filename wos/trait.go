package wos

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
)

// IReadCloser defines interface with function: setReadCloser
type IReadCloser interface {
	setReadCloser(body io.ReadCloser)
}

func (output *GetObjectOutput) setReadCloser(body io.ReadCloser) {
	output.Body = body
}

func (output *GetAvinfoOutput) setReadCloser(body io.ReadCloser) {
	output.Body = body
}

func setHeaders(headers map[string][]string, header string, headerValue []string, isWos bool) {
	if isWos {
		header = HEADER_PREFIX_WOS + header
		headers[header] = headerValue
	} else {
		header = HEADER_PREFIX + header
		headers[header] = headerValue
	}
}

func setHeadersNext(headers map[string][]string, header string, headerNext string, headerValue []string, isWos bool) {
	if isWos {
		headers[header] = headerValue
	} else {
		headers[headerNext] = headerValue
	}
}

// IBaseModel defines interface for base response model
type IBaseModel interface {
	setStatusCode(statusCode int)

	setRequestID(requestID string)

	setResponseHeaders(responseHeaders map[string][]string)
}

// ISerializable defines interface with function: trans
type ISerializable interface {
	trans(isWos bool) (map[string]string, map[string][]string, interface{}, error)
}

// DefaultSerializable defines default serializable struct
type DefaultSerializable struct {
	params  map[string]string
	headers map[string][]string
	data    interface{}
}

func (input DefaultSerializable) trans(isWos bool) (map[string]string, map[string][]string, interface{}, error) {
	return input.params, input.headers, input.data, nil
}

var defaultSerializable = &DefaultSerializable{}

func newSubResourceSerial(subResource SubResourceType) *DefaultSerializable {
	return &DefaultSerializable{map[string]string{string(subResource): ""}, nil, nil}
}

func trans(subResource SubResourceType, input interface{}) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(subResource): ""}
	data, err = ConvertRequestToIoReader(input)
	return
}

func (baseModel *BaseModel) setStatusCode(statusCode int) {
	baseModel.StatusCode = statusCode
}

func (baseModel *BaseModel) setRequestID(requestID string) {
	baseModel.RequestId = requestID
}

func (baseModel *BaseModel) setResponseHeaders(responseHeaders map[string][]string) {
	baseModel.ResponseHeaders = responseHeaders
}

func (input ListBucketsInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	headers = make(map[string][]string)
	return
}

func (input ListObjsInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = make(map[string]string)
	if input.Prefix != "" {
		params["prefix"] = input.Prefix
	}
	if input.Delimiter != "" {
		params["delimiter"] = input.Delimiter
	}
	if input.MaxKeys > 0 {
		params["max-keys"] = IntToString(input.MaxKeys)
	}
	if input.EncodingType != "" {
		params["encoding-type"] = input.EncodingType
	}
	return
}

func (input ListObjectsInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params, headers, data, err = input.ListObjsInput.trans(isWos)
	if err != nil {
		return
	}
	if input.Marker != "" {
		params["marker"] = input.Marker
	}
	return
}

func (input ListObjectsV2Input) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params, headers, data, err = input.ListObjsInput.trans(isWos)
	if err != nil {
		return
	}
	if input.StartAfter != "" {
		params["start-after"] = input.StartAfter
	}
	if input.ContinuationToken != "" {
		params["continuation-token"] = input.ContinuationToken
	}
	if input.FetchOwner {
		params["fetch-owner"] = "true"
	}
	params["list-type"] = "2"
	return
}

func (input ListMultipartUploadsInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceUploads): ""}
	if input.Prefix != "" {
		params["prefix"] = input.Prefix
	}
	if input.Delimiter != "" {
		params["delimiter"] = input.Delimiter
	}
	if input.MaxUploads > 0 {
		params["max-uploads"] = IntToString(input.MaxUploads)
	}
	if input.KeyMarker != "" {
		params["key-marker"] = input.KeyMarker
	}
	if input.UploadIdMarker != "" {
		params["upload-id-marker"] = input.UploadIdMarker
	}
	if input.EncodingType != "" {
		params["encoding-type"] = input.EncodingType
	}
	return
}

func (input SetBucketLifecycleConfigurationInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceLifecycle): ""}
	data, md5 := ConvertLifecyleConfigurationToXml(input.BucketLifecyleConfiguration, true, isWos)
	headers = map[string][]string{HEADER_MD5_CAMEL: {md5}}
	return
}

func (input DeleteObjectInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = make(map[string]string)
	return
}

func (input DeleteObjectsInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceDelete): ""}
	for index, object := range input.Objects {
		urlEncodeStr := url.QueryEscape(object.Key)
		urlEncodeStr = strings.ReplaceAll(urlEncodeStr, "%2F", "/")
		urlEncodeStr = strings.ReplaceAll(urlEncodeStr, "+", "%20")
		input.Objects[index].Key = urlEncodeStr
	}
	data, md5 := convertDeleteObjectsToXML(input)
	headers = map[string][]string{HEADER_MD5_CAMEL: {md5}}
	return
}

func (input RestoreObjectInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{string(SubResourceRestore): ""}
	if !isWos {
		data, err = ConvertRequestToIoReader(input)
	} else {
		data = ConverntWosRestoreToXml(input)
	}
	return
}

// GetEncryption gets the Encryption field value from SseKmsHeader
func (header SseKmsHeader) GetEncryption() string {
	if header.Encryption != "" {
		return header.Encryption
	}
	if !header.isWos {
		return DEFAULT_SSE_KMS_ENCRYPTION
	}
	return DEFAULT_SSE_KMS_ENCRYPTION_WOS
}

// GetKey gets the Key field value from SseKmsHeader
func (header SseKmsHeader) GetKey() string {
	return header.Key
}

// GetEncryption gets the Encryption field value from SseCHeader
func (header SseCHeader) GetEncryption() string {
	if header.Encryption != "" {
		return header.Encryption
	}
	return DEFAULT_SSE_C_ENCRYPTION
}

// GetKey gets the Key field value from SseCHeader
func (header SseCHeader) GetKey() string {
	return header.Key
}

// GetKeyMD5 gets the KeyMD5 field value from SseCHeader
func (header SseCHeader) GetKeyMD5() string {
	if header.KeyMD5 != "" {
		return header.KeyMD5
	}

	if ret, err := Base64Decode(header.GetKey()); err == nil {
		return Base64Md5(ret)
	}
	return ""
}

func setSseHeader(headers map[string][]string, sseHeader ISseHeader, sseCOnly bool, isWos bool) {
	if sseHeader != nil {
		if sseCHeader, ok := sseHeader.(SseCHeader); ok {
			setHeaders(headers, HEADER_SSEC_ENCRYPTION, []string{sseCHeader.GetEncryption()}, isWos)
			setHeaders(headers, HEADER_SSEC_KEY, []string{sseCHeader.GetKey()}, isWos)
			setHeaders(headers, HEADER_SSEC_KEY_MD5, []string{sseCHeader.GetKeyMD5()}, isWos)
		} else if sseKmsHeader, ok := sseHeader.(SseKmsHeader); !sseCOnly && ok {
			sseKmsHeader.isWos = isWos
			setHeaders(headers, HEADER_SSEKMS_ENCRYPTION, []string{sseKmsHeader.GetEncryption()}, isWos)
			if sseKmsHeader.GetKey() != "" {
				setHeadersNext(headers, HEADER_SSEKMS_KEY_WOS, HEADER_SSEKMS_KEY_AMZ, []string{sseKmsHeader.GetKey()}, isWos)
			}
		}
	}
}

func (input GetObjectMetadataInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = make(map[string]string)
	headers = make(map[string][]string)

	if input.Origin != "" {
		headers[HEADER_ORIGIN_CAMEL] = []string{input.Origin}
	}

	if input.RequestHeader != "" {
		headers[HEADER_ACCESS_CONTROL_REQUEST_HEADER_CAMEL] = []string{input.RequestHeader}
	}
	setSseHeader(headers, input.SseHeader, true, isWos)
	return
}

func (input SetObjectMetadataInput) prepareContentHeaders(headers map[string][]string) {
	if input.ContentDisposition != "" {
		headers[HEADER_CONTENT_DISPOSITION_CAMEL] = []string{input.ContentDisposition}
	}
	if input.ContentEncoding != "" {
		headers[HEADER_CONTENT_ENCODING_CAMEL] = []string{input.ContentEncoding}
	}
	if input.ContentLanguage != "" {
		headers[HEADER_CONTENT_LANGUAGE_CAMEL] = []string{input.ContentLanguage}
	}

	if input.ContentType != "" {
		headers[HEADER_CONTENT_TYPE_CAML] = []string{input.ContentType}
	}
}

func (input SetObjectMetadataInput) prepareStorageClass(headers map[string][]string, isWos bool) {
	if storageClass := string(input.StorageClass); storageClass != "" {
		setHeaders(headers, HEADER_STORAGE_CLASS2, []string{storageClass}, isWos)
	}
}

func (input SetObjectMetadataInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = make(map[string]string)
	params = map[string]string{string(SubResourceMetadata): ""}

	headers = make(map[string][]string)

	if directive := string(input.MetadataDirective); directive != "" {
		setHeaders(headers, HEADER_METADATA_DIRECTIVE, []string{string(input.MetadataDirective)}, isWos)
	} else {
		setHeaders(headers, HEADER_METADATA_DIRECTIVE, []string{string(ReplaceNew)}, isWos)
	}
	if input.CacheControl != "" {
		headers[HEADER_CACHE_CONTROL_CAMEL] = []string{input.CacheControl}
	}
	input.prepareContentHeaders(headers)
	if input.Expires != "" {
		headers[HEADER_EXPIRES_CAMEL] = []string{input.Expires}
	}
	if input.WebsiteRedirectLocation != "" {
		setHeaders(headers, HEADER_WEBSITE_REDIRECT_LOCATION, []string{input.WebsiteRedirectLocation}, isWos)
	}
	input.prepareStorageClass(headers, isWos)
	if input.Metadata != nil {
		for key, value := range input.Metadata {
			key = strings.TrimSpace(key)
			setHeadersNext(headers, HEADER_PREFIX_META_WOS+key, HEADER_PREFIX_META+key, []string{value}, isWos)
		}
	}
	return
}

func (input GetObjectInput) prepareResponseParams(params map[string]string) {
	if input.ResponseCacheControl != "" {
		params[PARAM_RESPONSE_CACHE_CONTROL] = input.ResponseCacheControl
	}
	if input.ResponseContentDisposition != "" {
		params[PARAM_RESPONSE_CONTENT_DISPOSITION] = input.ResponseContentDisposition
	}
	if input.ResponseContentEncoding != "" {
		params[PARAM_RESPONSE_CONTENT_ENCODING] = input.ResponseContentEncoding
	}
	if input.ResponseContentLanguage != "" {
		params[PARAM_RESPONSE_CONTENT_LANGUAGE] = input.ResponseContentLanguage
	}
	if input.ResponseContentType != "" {
		params[PARAM_RESPONSE_CONTENT_TYPE] = input.ResponseContentType
	}
	if input.ResponseExpires != "" {
		params[PARAM_RESPONSE_EXPIRES] = input.ResponseExpires
	}
}

func (input GetObjectInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params, headers, data, err = input.GetObjectMetadataInput.trans(isWos)
	if err != nil {
		return
	}
	input.prepareResponseParams(params)
	if input.RangeStart >= 0 && input.RangeEnd > input.RangeStart {
		headers[HEADER_RANGE] = []string{fmt.Sprintf("bytes=%d-%d", input.RangeStart, input.RangeEnd)}
	}

	if input.IfMatch != "" {
		headers[HEADER_IF_MATCH] = []string{input.IfMatch}
	}
	if input.IfNoneMatch != "" {
		headers[HEADER_IF_NONE_MATCH] = []string{input.IfNoneMatch}
	}
	if !input.IfModifiedSince.IsZero() {
		headers[HEADER_IF_MODIFIED_SINCE] = []string{FormatUtcToRfc1123(input.IfModifiedSince)}
	}
	if !input.IfUnmodifiedSince.IsZero() {
		headers[HEADER_IF_UNMODIFIED_SINCE] = []string{FormatUtcToRfc1123(input.IfUnmodifiedSince)}
	}
	return
}

func (input ObjectOperationInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	headers = make(map[string][]string)
	params = make(map[string]string)
	if storageClass := string(input.StorageClass); storageClass != "" {
		setHeaders(headers, HEADER_STORAGE_CLASS2, []string{storageClass}, isWos)
	}
	if input.Metadata != nil {
		for key, value := range input.Metadata {
			key = strings.TrimSpace(key)
			setHeadersNext(headers, HEADER_PREFIX_META_WOS+key, HEADER_PREFIX_META+key, []string{value}, isWos)
		}
	}
	return
}

func (input PutObjectBasicInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params, headers, data, err = input.ObjectOperationInput.trans(isWos)
	if err != nil {
		return
	}

	if input.ContentMD5 != "" {
		headers[HEADER_MD5_CAMEL] = []string{input.ContentMD5}
	}

	if input.ContentLength > 0 {
		headers[HEADER_CONTENT_LENGTH_CAMEL] = []string{Int64ToString(input.ContentLength)}
	}
	if input.ContentType != "" {
		headers[HEADER_CONTENT_TYPE_CAML] = []string{input.ContentType}
	}

	return
}

func (input PutObjectInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params, headers, data, err = input.PutObjectBasicInput.trans(isWos)
	if err != nil {
		return
	}
	if input.Body != nil {
		data = input.Body
	}
	return
}

func (input CopyObjectInput) prepareReplaceHeaders(headers map[string][]string) {
	if input.CacheControl != "" {
		headers[HEADER_CACHE_CONTROL] = []string{input.CacheControl}
	}
	if input.ContentDisposition != "" {
		headers[HEADER_CONTENT_DISPOSITION] = []string{input.ContentDisposition}
	}
	if input.ContentEncoding != "" {
		headers[HEADER_CONTENT_ENCODING] = []string{input.ContentEncoding}
	}
	if input.ContentLanguage != "" {
		headers[HEADER_CONTENT_LANGUAGE] = []string{input.ContentLanguage}
	}
	if input.ContentType != "" {
		headers[HEADER_CONTENT_TYPE] = []string{input.ContentType}
	}
	if input.Expires != "" {
		headers[HEADER_EXPIRES] = []string{input.Expires}
	}
}

func (input CopyObjectInput) prepareCopySourceHeaders(headers map[string][]string, isWos bool) {
	if input.CopySourceIfMatch != "" {
		setHeaders(headers, HEADER_COPY_SOURCE_IF_MATCH, []string{input.CopySourceIfMatch}, isWos)
	}
	if input.CopySourceIfNoneMatch != "" {
		setHeaders(headers, HEADER_COPY_SOURCE_IF_NONE_MATCH, []string{input.CopySourceIfNoneMatch}, isWos)
	}
	if !input.CopySourceIfModifiedSince.IsZero() {
		setHeaders(headers, HEADER_COPY_SOURCE_IF_MODIFIED_SINCE, []string{FormatUtcToRfc1123(input.CopySourceIfModifiedSince)}, isWos)
	}
	if !input.CopySourceIfUnmodifiedSince.IsZero() {
		setHeaders(headers, HEADER_COPY_SOURCE_IF_UNMODIFIED_SINCE, []string{FormatUtcToRfc1123(input.CopySourceIfUnmodifiedSince)}, isWos)
	}
}

func (input CopyObjectInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params, headers, data, err = input.ObjectOperationInput.trans(isWos)
	if err != nil {
		return
	}

	copySource := fmt.Sprintf("%s/%s", input.CopySourceBucket, UrlEncode(input.CopySourceKey, false))
	setHeaders(headers, HEADER_COPY_SOURCE, []string{copySource}, isWos)

	if directive := string(input.MetadataDirective); directive != "" {
		setHeaders(headers, HEADER_METADATA_DIRECTIVE, []string{directive}, isWos)
	}

	if input.MetadataDirective == ReplaceMetadata {
		input.prepareReplaceHeaders(headers)
	}

	input.prepareCopySourceHeaders(headers, isWos)
	return
}

func (input AbortMultipartUploadInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{"uploadId": input.UploadId}
	return
}

func (input InitiateMultipartUploadInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params, headers, data, err = input.ObjectOperationInput.trans(isWos)
	if err != nil {
		return
	}
	if input.ContentType != "" {
		headers[HEADER_CONTENT_TYPE_CAML] = []string{input.ContentType}
	}
	params[string(SubResourceUploads)] = ""
	if input.EncodingType != "" {
		params["encoding-type"] = input.EncodingType
	}
	return
}

func (input UploadPartInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{"uploadId": input.UploadId, "partNumber": IntToString(input.PartNumber)}
	headers = make(map[string][]string)
	if input.ContentMD5 != "" {
		headers[HEADER_MD5_CAMEL] = []string{input.ContentMD5}
	}
	if input.Body != nil {
		data = input.Body
	}
	return
}

func (input CompleteMultipartUploadInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{"uploadId": input.UploadId}
	if input.EncodingType != "" {
		params["encoding-type"] = input.EncodingType
	}
	data, _ = ConvertCompleteMultipartUploadInputToXml(input, false)
	return
}

func (input ListPartsInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{"uploadId": input.UploadId}
	if input.MaxParts > 0 {
		params["max-parts"] = IntToString(input.MaxParts)
	}
	if input.PartNumberMarker > 0 {
		params["part-number-marker"] = IntToString(input.PartNumberMarker)
	}
	return
}

func (input CopyPartInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	params = map[string]string{"uploadId": input.UploadId, "partNumber": IntToString(input.PartNumber)}
	headers = make(map[string][]string, 1)
	var copySource string
	copySource = fmt.Sprintf("%s/%s", input.CopySourceBucket, UrlEncode(input.CopySourceKey, false))
	setHeaders(headers, HEADER_COPY_SOURCE, []string{copySource}, isWos)
	if input.CopySourceRangeStart >= 0 && input.CopySourceRangeEnd > input.CopySourceRangeStart {
		setHeaders(headers, HEADER_COPY_SOURCE_RANGE, []string{fmt.Sprintf("bytes=%d-%d", input.CopySourceRangeStart, input.CopySourceRangeEnd)}, isWos)
	}

	setSseHeader(headers, input.SseHeader, true, isWos)
	if input.SourceSseHeader != nil {
		if sseCHeader, ok := input.SourceSseHeader.(SseCHeader); ok {
			setHeaders(headers, HEADER_SSEC_COPY_SOURCE_ENCRYPTION, []string{sseCHeader.GetEncryption()}, isWos)
			setHeaders(headers, HEADER_SSEC_COPY_SOURCE_KEY, []string{sseCHeader.GetKey()}, isWos)
			setHeaders(headers, HEADER_SSEC_COPY_SOURCE_KEY_MD5, []string{sseCHeader.GetKeyMD5()}, isWos)
		}

	}
	return
}

func (input HeadObjectInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	headers = make(map[string][]string)
	if input.IfMatch != "" {
		headers[HEADER_IF_MATCH] = []string{input.IfMatch}
	}
	if input.IfNoneMatch != "" {
		headers[HEADER_IF_NONE_MATCH] = []string{input.IfNoneMatch}
	}
	if !input.IfModifiedSince.IsZero() {
		headers[HEADER_IF_MODIFIED_SINCE] = []string{FormatUtcToRfc1123(input.IfModifiedSince)}
	}
	if !input.IfUnmodifiedSince.IsZero() {
		headers[HEADER_IF_UNMODIFIED_SINCE] = []string{FormatUtcToRfc1123(input.IfUnmodifiedSince)}
	}
	if input.RangeStart >= 0 && input.RangeEnd > input.RangeStart {
		headers[HEADER_RANGE] = []string{fmt.Sprintf("bytes=%d-%d", input.RangeStart, input.RangeEnd)}
	}
	return
}

type partSlice []Part

func (parts partSlice) Len() int {
	return len(parts)
}

func (parts partSlice) Less(i, j int) bool {
	return parts[i].PartNumber < parts[j].PartNumber
}

func (parts partSlice) Swap(i, j int) {
	parts[i], parts[j] = parts[j], parts[i]
}

type readerWrapper struct {
	reader      io.Reader
	mark        int64
	totalCount  int64
	readedCount int64
}

func (rw *readerWrapper) seek(offset int64, whence int) (int64, error) {
	if r, ok := rw.reader.(*strings.Reader); ok {
		return r.Seek(offset, whence)
	} else if r, ok := rw.reader.(*bytes.Reader); ok {
		return r.Seek(offset, whence)
	} else if r, ok := rw.reader.(*os.File); ok {
		return r.Seek(offset, whence)
	}
	return offset, nil
}

func (rw *readerWrapper) Read(p []byte) (n int, err error) {
	if rw.totalCount == 0 {
		return 0, io.EOF
	}
	if rw.totalCount > 0 {
		n, err = rw.reader.Read(p)
		readedOnce := int64(n)
		remainCount := rw.totalCount - rw.readedCount
		if remainCount > readedOnce {
			rw.readedCount += readedOnce
			return n, err
		}
		rw.readedCount += remainCount
		return int(remainCount), io.EOF
	}
	return rw.reader.Read(p)
}

type fileReaderWrapper struct {
	readerWrapper
	filePath string
}

func (input SetBucketFetchPolicyInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	contentType, _ := mimeTypes["json"]
	headers = make(map[string][]string, 2)
	headers[HEADER_CONTENT_TYPE] = []string{contentType}
	setHeaders(headers, headerOefMarker, []string{"yes"}, isWos)
	data, err = convertFetchPolicyToJSON(input)
	return
}

func (input GetBucketFetchPolicyInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	headers = make(map[string][]string, 1)
	setHeaders(headers, headerOefMarker, []string{"yes"}, isWos)
	return
}

func (input DeleteBucketFetchPolicyInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	headers = make(map[string][]string, 1)
	setHeaders(headers, headerOefMarker, []string{"yes"}, isWos)
	return
}

func (input SetBucketFetchJobInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	contentType, _ := mimeTypes["json"]
	headers = make(map[string][]string, 2)
	headers[HEADER_CONTENT_TYPE] = []string{contentType}
	setHeaders(headers, headerOefMarker, []string{"yes"}, isWos)
	data, err = convertFetchJobToJSON(input)
	return
}

func (input GetBucketFetchJobInput) trans(isWos bool) (params map[string]string, headers map[string][]string, data interface{}, err error) {
	headers = make(map[string][]string, 1)
	setHeaders(headers, headerOefMarker, []string{"yes"}, isWos)
	return
}
