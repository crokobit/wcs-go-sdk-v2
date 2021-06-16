package wos

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

// BaseModel defines base model response from WOS
type BaseModel struct {
	StatusCode      int                 `xml:"-"`
	RequestId       string              `xml:"RequestId" json:"request_id"`
	ResponseHeaders map[string][]string `xml:"-"`
}

// Bucket defines bucket properties
type Bucket struct {
	XMLName      xml.Name  `xml:"Bucket"`
	Name         string    `xml:"Name"`
	CreationDate time.Time `xml:"CreationDate"`
	Endpoint     string    `xml:"Endpoint,omitempty"`
	Region       string    `xml:Region,omitempty`
}

// Owner defines owner properties
type Owner struct {
	XMLName     xml.Name `xml:"Owner"`
	ID          string   `xml:"ID"`
	DisplayName string   `xml:"DisplayName,omitempty"`
}

// Initiator defines initiator properties
type Initiator struct {
	XMLName     xml.Name `xml:"Initiator"`
	ID          string   `xml:"ID"`
	DisplayName string   `xml:"DisplayName,omitempty"`
}

// ListBucketsInput is the input parameter of ListBuckets function
type ListBucketsInput struct {
	QueryLocation bool
}

// ListBucketsOutput is the result of ListBuckets function
type ListBucketsOutput struct {
	BaseModel
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Owner   Owner    `xml:"Owner"`
	Buckets []Bucket `xml:"Buckets>Bucket"`
}

type bucketLocationWos struct {
	XMLName  xml.Name `xml:"Location"`
	Location string   `xml:",chardata"`
}

// BucketLocation defines bucket location configuration
type BucketLocation struct {
	XMLName  xml.Name `xml:"CreateBucketConfiguration"`
	Location string   `xml:"LocationConstraint,omitempty"`
}

type bucketStoragePolicyWos struct {
	XMLName      xml.Name `xml:"StorageClass"`
	StorageClass string   `xml:",chardata"`
}
type getBucketStoragePolicyOutputWos struct {
	BaseModel
	bucketStoragePolicyWos
}

// ListObjsInput defines parameters for listing objects
type ListObjsInput struct {
	Prefix       string
	MaxKeys      int
	Delimiter    string
	EncodingType string
}

// ListObjectsInput is the input parameter of ListObjects function
type ListObjectsInput struct {
	ListObjsInput
	Bucket string
	Marker string
}

// ListObjectsInput is the input parameter of ListObjectsV2 function
type ListObjectsV2Input struct {
	ListObjsInput
	Bucket            string
	ContinuationToken string
	FetchOwner        bool
	StartAfter        string
}

// Content defines the object content properties
type Content struct {
	XMLName      xml.Name         `xml:"Contents"`
	ETag         string           `xml:"ETag"`
	Key          string           `xml:"Key"`
	LastModified time.Time        `xml:"LastModified"`
	Size         int64            `xml:"Size"`
	StorageClass StorageClassType `xml:"StorageClass"`
}

// ListObjectsOutput is the result of ListObjects function
type ListObjectsOutput struct {
	BaseModel
	XMLName        xml.Name  `xml:"ListBucketResult"`
	Delimiter      string    `xml:"Delimiter"`
	IsTruncated    bool      `xml:"IsTruncated"`
	Marker         string    `xml:"Marker"`
	NextMarker     string    `xml:"NextMarker"`
	MaxKeys        int       `xml:"MaxKeys"`
	Name           string    `xml:"Name"`
	Prefix         string    `xml:"Prefix"`
	Contents       []Content `xml:"Contents"`
	CommonPrefixes []string  `xml:"CommonPrefixes>Prefix"`
	Owner          Owner     `xml:"Owner"`
}

// ListObjectsV2Output is the result of ListObjectsV2 function
type ListObjectsV2Output struct {
	BaseModel
	XMLName               xml.Name  `xml:"ListBucketResult"`
	Delimiter             string    `xml:"Delimiter"`
	IsTruncated           bool      `xml:"IsTruncated"`
	MaxKeys               int       `xml:"MaxKeys"`
	Name                  string    `xml:"Name"`
	Prefix                string    `xml:"Prefix"`
	Contents              []Content `xml:"Contents"`
	CommonPrefixes        []string  `xml:"CommonPrefixes>Prefix"`
	Owner                 Owner     `xml:"Owner"`
	ContinuationToken     string    `xml:"ContinuationToken"`
	NextContinuationToken string    `xml:"NextContinuationToken"`
}

// ListMultipartUploadsInput is the input parameter of ListMultipartUploads function
type ListMultipartUploadsInput struct {
	Bucket         string
	Prefix         string
	MaxUploads     int
	Delimiter      string
	KeyMarker      string
	UploadIdMarker string
	EncodingType   string
}

// Upload defines multipart upload properties
type Upload struct {
	XMLName      xml.Name         `xml:"Upload"`
	Key          string           `xml:"Key"`
	UploadId     string           `xml:"UploadId"`
	Initiated    time.Time        `xml:"Initiated"`
	StorageClass StorageClassType `xml:"StorageClass"`
	Owner        Owner            `xml:"Owner"`
	Initiator    Initiator        `xml:"Initiator"`
}

// ListMultipartUploadsOutput is the result of ListMultipartUploads function
type ListMultipartUploadsOutput struct {
	BaseModel
	XMLName            xml.Name `xml:"ListMultipartUploadsResult"`
	Bucket             string   `xml:"Bucket"`
	KeyMarker          string   `xml:"KeyMarker"`
	NextKeyMarker      string   `xml:"NextKeyMarker"`
	UploadIdMarker     string   `xml:"UploadIdMarker"`
	NextUploadIdMarker string   `xml:"NextUploadIdMarker"`
	Delimiter          string   `xml:"Delimiter"`
	IsTruncated        bool     `xml:"IsTruncated"`
	MaxUploads         int      `xml:"MaxUploads"`
	Prefix             string   `xml:"Prefix"`
	Uploads            []Upload `xml:"Upload"`
	CommonPrefixes     []string `xml:"CommonPrefixes>Prefix"`
	EncodingType       string   `xml:"EncodingType,omitempty"`
}

type getBucketLocationOutputS3 struct {
	BaseModel
	BucketLocation
}
type getBucketLocationOutputWos struct {
	BaseModel
	bucketLocationWos
}

// Grantee defines grantee properties
type Grantee struct {
	XMLName     xml.Name     `xml:"Grantee"`
	Type        GranteeType  `xml:"type,attr"`
	ID          string       `xml:"ID,omitempty"`
	DisplayName string       `xml:"DisplayName,omitempty"`
	URI         GroupUriType `xml:"URI,omitempty"`
}

type granteeWos struct {
	XMLName     xml.Name    `xml:"Grantee"`
	Type        GranteeType `xml:"type,attr"`
	ID          string      `xml:"ID,omitempty"`
	DisplayName string      `xml:"DisplayName,omitempty"`
	Canned      string      `xml:"Canned,omitempty"`
}

// Grant defines grant properties
type Grant struct {
	XMLName    xml.Name       `xml:"Grant"`
	Grantee    Grantee        `xml:"Grantee"`
	Permission PermissionType `xml:"Permission"`
	Delivered  bool           `xml:"Delivered"`
}
type grantWos struct {
	XMLName    xml.Name       `xml:"Grant"`
	Grantee    granteeWos     `xml:"Grantee"`
	Permission PermissionType `xml:"Permission"`
	Delivered  bool           `xml:"Delivered"`
}

// AccessControlPolicy defines access control policy properties
type AccessControlPolicy struct {
	XMLName   xml.Name `xml:"AccessControlPolicy"`
	Owner     Owner    `xml:"Owner"`
	Grants    []Grant  `xml:"AccessControlList>Grant"`
	Delivered string   `xml:"Delivered,omitempty"`
}

type accessControlPolicyWos struct {
	XMLName xml.Name   `xml:"AccessControlPolicy"`
	Owner   Owner      `xml:"Owner"`
	Grants  []grantWos `xml:"AccessControlList>Grant"`
}

// CorsRule defines the CORS rules
type CorsRule struct {
	XMLName       xml.Name `xml:"CORSRule"`
	ID            string   `xml:"ID,omitempty"`
	AllowedOrigin []string `xml:"AllowedOrigin"`
	AllowedMethod []string `xml:"AllowedMethod"`
	AllowedHeader []string `xml:"AllowedHeader,omitempty"`
	MaxAgeSeconds int      `xml:"MaxAgeSeconds"`
	ExposeHeader  []string `xml:"ExposeHeader,omitempty"`
}

// IndexDocument defines the default page configuration
type IndexDocument struct {
	Suffix string `xml:"Suffix"`
}

// ErrorDocument defines the error page configuration
type ErrorDocument struct {
	Key string `xml:"Key,omitempty"`
}

// Condition defines condition in RoutingRule
type Condition struct {
	XMLName                     xml.Name `xml:"Condition"`
	KeyPrefixEquals             string   `xml:"KeyPrefixEquals,omitempty"`
	HttpErrorCodeReturnedEquals string   `xml:"HttpErrorCodeReturnedEquals,omitempty"`
}

// Redirect defines redirect in RoutingRule
type Redirect struct {
	XMLName              xml.Name     `xml:"Redirect"`
	Protocol             ProtocolType `xml:"Protocol,omitempty"`
	HostName             string       `xml:"HostName,omitempty"`
	ReplaceKeyPrefixWith string       `xml:"ReplaceKeyPrefixWith,omitempty"`
	ReplaceKeyWith       string       `xml:"ReplaceKeyWith,omitempty"`
	HttpRedirectCode     string       `xml:"HttpRedirectCode,omitempty"`
}

// RoutingRule defines routing rules
type RoutingRule struct {
	XMLName   xml.Name  `xml:"RoutingRule"`
	Condition Condition `xml:"Condition,omitempty"`
	Redirect  Redirect  `xml:"Redirect"`
}

// RedirectAllRequestsTo defines redirect in BucketWebsiteConfiguration
type RedirectAllRequestsTo struct {
	XMLName  xml.Name     `xml:"RedirectAllRequestsTo"`
	Protocol ProtocolType `xml:"Protocol,omitempty"`
	HostName string       `xml:"HostName"`
}

// BucketWebsiteConfiguration defines the bucket website configuration
type BucketWebsiteConfiguration struct {
	XMLName               xml.Name              `xml:"WebsiteConfiguration"`
	RedirectAllRequestsTo RedirectAllRequestsTo `xml:"RedirectAllRequestsTo,omitempty"`
	IndexDocument         IndexDocument         `xml:"IndexDocument,omitempty"`
	ErrorDocument         ErrorDocument         `xml:"ErrorDocument,omitempty"`
	RoutingRules          []RoutingRule         `xml:"RoutingRules>RoutingRule,omitempty"`
}

// SetObjectMetadataInput is the input parameter of SetObjectMetadata function
type SetObjectMetadataInput struct {
	Bucket                  string
	Key                     string
	MetadataDirective       MetadataDirectiveType
	CacheControl            string
	ContentDisposition      string
	ContentEncoding         string
	ContentLanguage         string
	ContentType             string
	Expires                 string
	WebsiteRedirectLocation string
	StorageClass            StorageClassType
	Metadata                map[string]string
}

//SetObjectMetadataOutput is the result of SetObjectMetadata function
type SetObjectMetadataOutput struct {
	BaseModel
	MetadataDirective       MetadataDirectiveType
	CacheControl            string
	ContentDisposition      string
	ContentEncoding         string
	ContentLanguage         string
	ContentType             string
	Expires                 string
	WebsiteRedirectLocation string
	StorageClass            StorageClassType
	Metadata                map[string]string
}

// GetBucketMetadataOutput is the result of GetBucketMetadata function
type GetBucketMetadataOutput struct {
	BaseModel
	StorageClass  StorageClassType
	Location      string
	AllowOrigin   string
	AllowMethod   string
	AllowHeader   string
	MaxAgeSeconds int
	ExposeHeader  string
	Epid          string
	FSStatus      FSStatusType
}

// Filter defines filter property in LifecycleRule
type Filter struct {
	Prefix string `xml:"Prefix"`
}

// Transition defines transition property in LifecycleRule
type Transition struct {
	XMLName      xml.Name         `xml:"Transition"`
	Days         int              `xml:"Days,omitempty"`
	StorageClass StorageClassType `xml:"StorageClass"`
}

// Expiration defines expiration property in LifecycleRule
type Expiration struct {
	XMLName xml.Name `xml:"Expiration"`
	Days    int      `xml:"Days,omitempty"`
}

// LifecycleRule defines lifecycle rule
type LifecycleRule struct {
	ID          string         `xml:"ID,omitempty"`
	Filter      Filter         `xml:"Filter"`
	Status      RuleStatusType `xml:"Status"`
	Transitions []Transition   `xml:"Transition,omitempty"`
	Expiration  Expiration     `xml:"Expiration,omitempty"`
}

// BucketLifecyleConfiguration defines the bucket lifecycle configuration
type BucketLifecyleConfiguration struct {
	XMLName        xml.Name        `xml:"LifecycleConfiguration"`
	LifecycleRules []LifecycleRule `xml:"Rule"`
}

// SetBucketLifecycleConfigurationInput is the input parameter of SetBucketLifecycleConfiguration function
type SetBucketLifecycleConfigurationInput struct {
	Bucket string `xml:"-"`
	BucketLifecyleConfiguration
}

// GetBucketLifecycleConfigurationOutput is the result of GetBucketLifecycleConfiguration function
type GetBucketLifecycleConfigurationOutput struct {
	BaseModel
	BucketLifecyleConfiguration
}

// Tag defines tag property in BucketTagging
type Tag struct {
	XMLName xml.Name `xml:"Tag"`
	Key     string   `xml:"Key"`
	Value   string   `xml:"Value"`
}

// FilterRule defines filter rule in TopicConfiguration
type FilterRule struct {
	XMLName xml.Name `xml:"FilterRule"`
	Name    string   `xml:"Name,omitempty"`
	Value   string   `xml:"Value,omitempty"`
}

// TopicConfiguration defines the topic configuration
type TopicConfiguration struct {
	XMLName     xml.Name     `xml:"TopicConfiguration"`
	ID          string       `xml:"Id,omitempty"`
	Topic       string       `xml:"Topic"`
	Events      []EventType  `xml:"Event"`
	FilterRules []FilterRule `xml:"Filter>Object>FilterRule"`
}

type topicConfigurationS3 struct {
	XMLName     xml.Name     `xml:"TopicConfiguration"`
	ID          string       `xml:"Id,omitempty"`
	Topic       string       `xml:"Topic"`
	Events      []string     `xml:"Event"`
	FilterRules []FilterRule `xml:"Filter>S3Key>FilterRule"`
}

type bucketNotificationS3 struct {
	XMLName             xml.Name               `xml:"NotificationConfiguration"`
	TopicConfigurations []topicConfigurationS3 `xml:"TopicConfiguration"`
}

type getBucketNotificationOutputS3 struct {
	BaseModel
	bucketNotificationS3
}

// DeleteObjectInput is the input parameter of DeleteObject function
type DeleteObjectInput struct {
	Bucket string
	Key    string
}

// DeleteObjectOutput is the result of DeleteObject function
type DeleteObjectOutput struct {
	BaseModel
}

// ObjectToDelete defines the object property in DeleteObjectsInput
type ObjectToDelete struct {
	XMLName xml.Name `xml:"Object"`
	Key     string   `xml:"Key"`
}

// DeleteObjectsInput is the input parameter of DeleteObjects function
type DeleteObjectsInput struct {
	Bucket  string           `xml:"-"`
	XMLName xml.Name         `xml:"Delete"`
	Quiet   bool             `xml:"Quiet,omitempty"`
	Objects []ObjectToDelete `xml:"Object"`
}

// Deleted defines the deleted property in DeleteObjectsOutput
type Deleted struct {
	XMLName xml.Name `xml:"Deleted"`
	Key     string   `xml:"Key"`
}

// Error defines the error property in DeleteObjectsOutput
type Error struct {
	XMLName xml.Name `xml:"Error"`
	Key     string   `xml:"Key"`
	Code    string   `xml:"Code"`
	Message string   `xml:"Message"`
}

// DeleteObjectsOutput is the result of DeleteObjects function
type DeleteObjectsOutput struct {
	BaseModel
	XMLName  xml.Name  `xml:"DeleteResult"`
	Deleteds []Deleted `xml:"Deleted"`
	Errors   []Error   `xml:"Error"`
}

// RestoreObjectInput is the input parameter of RestoreObject function
type RestoreObjectInput struct {
	Bucket  string   `xml:"-"`
	Key     string   `xml:"-"`
	XMLName xml.Name `xml:"RestoreRequest"`
	Days    int      `xml:"Days"`
}

// ISseHeader defines the sse encryption header
type ISseHeader interface {
	GetEncryption() string
	GetKey() string
}

// SseKmsHeader defines the SseKms header
type SseKmsHeader struct {
	Encryption string
	Key        string
	isWos      bool
}

// SseCHeader defines the SseC header
type SseCHeader struct {
	Encryption string
	Key        string
	KeyMD5     string
}

// GetObjectMetadataInput is the input parameter of GetObjectMetadata function
type GetObjectMetadataInput struct {
	Bucket        string
	Key           string
	Origin        string
	RequestHeader string
	SseHeader     ISseHeader
}

// GetObjectMetadataOutput is the result of GetObjectMetadata function
type GetObjectMetadataOutput struct {
	BaseModel
	Expiration         string
	Restore            string
	ObjectType         string
	NextAppendPosition string
	StorageClass       StorageClassType
	ContentLength      int64
	ContentType        string
	ETag               string
	AllowOrigin        string
	AllowHeader        string
	AllowMethod        string
	ExposeHeader       string
	MaxAgeSeconds      int
	LastModified       time.Time
	Metadata           map[string]string
}

type GetAvinfoInput struct {
	Bucket        string
	Key           string
	RequestHeader string
	SseHeader     ISseHeader
}

// GetAvinfoOutput is the result of GetAvinfo function
type GetAvinfoOutput struct {
	BaseModel
	Body io.ReadCloser
}

// GetObjectInput is the input parameter of GetObject function
type GetObjectInput struct {
	GetObjectMetadataInput
	IfMatch                    string
	IfNoneMatch                string
	IfUnmodifiedSince          time.Time
	IfModifiedSince            time.Time
	RangeStart                 int64
	RangeEnd                   int64
	ResponseCacheControl       string
	ResponseContentDisposition string
	ResponseContentEncoding    string
	ResponseContentLanguage    string
	ResponseContentType        string
	ResponseExpires            string
}

// GetObjectOutput is the result of GetObject function
type GetObjectOutput struct {
	GetObjectMetadataOutput
	DeleteMarker       bool
	CacheControl       string
	ContentDisposition string
	ContentEncoding    string
	ContentLanguage    string
	Expires            string
	Body               io.ReadCloser
}

// ObjectOperationInput defines the object operation properties
type ObjectOperationInput struct {
	Bucket       string
	Key          string
	StorageClass StorageClassType
	Metadata     map[string]string
}

// PutObjectBasicInput defines the basic object operation properties
type PutObjectBasicInput struct {
	ObjectOperationInput
	ContentType   string
	ContentMD5    string
	ContentLength int64
}

// PutObjectInput is the input parameter of PutObject function
type PutObjectInput struct {
	PutObjectBasicInput
	Body io.Reader
}

// PutFileInput is the input parameter of PutFile function
type PutFileInput struct {
	PutObjectBasicInput
	SourceFile string
}

// PutObjectOutput is the result of PutObject function
type PutObjectOutput struct {
	BaseModel
	SseHeader ISseHeader
	ETag      string
}

// CopyObjectInput is the input parameter of CopyObject function
type CopyObjectInput struct {
	ObjectOperationInput
	CopySourceBucket            string
	CopySourceKey               string
	CopySourceIfMatch           string
	CopySourceIfNoneMatch       string
	CopySourceIfUnmodifiedSince time.Time
	CopySourceIfModifiedSince   time.Time
	CacheControl                string
	ContentDisposition          string
	ContentEncoding             string
	ContentLanguage             string
	ContentType                 string
	Expires                     string
	MetadataDirective           MetadataDirectiveType
}

// CopyObjectOutput is the result of CopyObject function
type CopyObjectOutput struct {
	BaseModel
	XMLName      xml.Name  `xml:"CopyObjectResult"`
	LastModified time.Time `xml:"LastModified"`
	ETag         string    `xml:"ETag"`
}

// AbortMultipartUploadInput is the input parameter of AbortMultipartUpload function
type AbortMultipartUploadInput struct {
	Bucket   string
	Key      string
	UploadId string
}

// InitiateMultipartUploadInput is the input parameter of InitiateMultipartUpload function
type InitiateMultipartUploadInput struct {
	ObjectOperationInput
	ContentType  string
	EncodingType string
}

// InitiateMultipartUploadOutput is the result of InitiateMultipartUpload function
type InitiateMultipartUploadOutput struct {
	BaseModel
	XMLName      xml.Name `xml:"InitiateMultipartUploadResult"`
	Bucket       string   `xml:"Bucket"`
	Key          string   `xml:"Key"`
	UploadId     string   `xml:"UploadId"`
	SseHeader    ISseHeader
	EncodingType string `xml:"EncodingType,omitempty"`
}

// UploadPartInput is the input parameter of UploadPart function
type UploadPartInput struct {
	Bucket     string
	Key        string
	PartNumber int
	UploadId   string
	ContentMD5 string
	Body       io.Reader
	SourceFile string
	Offset     int64
	PartSize   int64
}

// UploadPartOutput is the result of UploadPart function
type UploadPartOutput struct {
	BaseModel
	PartNumber int
	ETag       string
}

// Part defines the part properties
type Part struct {
	XMLName      xml.Name  `xml:"Part"`
	PartNumber   int       `xml:"PartNumber"`
	ETag         string    `xml:"ETag"`
	LastModified time.Time `xml:"LastModified,omitempty"`
	Size         int64     `xml:"Size,omitempty"`
}

// CompleteMultipartUploadInput is the input parameter of CompleteMultipartUpload function
type CompleteMultipartUploadInput struct {
	Bucket       string   `xml:"-"`
	Key          string   `xml:"-"`
	UploadId     string   `xml:"-"`
	XMLName      xml.Name `xml:"CompleteMultipartUpload"`
	Parts        []Part   `xml:"Part"`
	EncodingType string   `xml:"-"`
}

// CompleteMultipartUploadOutput is the result of CompleteMultipartUpload function
type CompleteMultipartUploadOutput struct {
	BaseModel
	SseHeader    ISseHeader `xml:"-"`
	XMLName      xml.Name   `xml:"CompleteMultipartUploadResult"`
	Location     string     `xml:"Location"`
	Bucket       string     `xml:"Bucket"`
	Key          string     `xml:"Key"`
	ETag         string     `xml:"ETag"`
	EncodingType string     `xml:"EncodingType,omitempty"`
}

// ListPartsInput is the input parameter of ListParts function
type ListPartsInput struct {
	Bucket           string
	Key              string
	UploadId         string
	MaxParts         int
	PartNumberMarker int
}

// ListPartsOutput is the result of ListParts function
type ListPartsOutput struct {
	BaseModel
	XMLName              xml.Name         `xml:"ListPartsResult"`
	Bucket               string           `xml:"Bucket"`
	Key                  string           `xml:"Key"`
	UploadId             string           `xml:"UploadId"`
	PartNumberMarker     int              `xml:"PartNumberMarker"`
	NextPartNumberMarker int              `xml:"NextPartNumberMarker"`
	MaxParts             int              `xml:"MaxParts"`
	IsTruncated          bool             `xml:"IsTruncated"`
	StorageClass         StorageClassType `xml:"StorageClass"`
	Initiator            Initiator        `xml:"Initiator"`
	Owner                Owner            `xml:"Owner"`
	Parts                []Part           `xml:"Part"`
	EncodingType         string           `xml:"EncodingType,omitempty"`
}

// CopyPartInput is the input parameter of CopyPart function
type CopyPartInput struct {
	Bucket               string
	Key                  string
	UploadId             string
	PartNumber           int
	CopySourceBucket     string
	CopySourceKey        string
	CopySourceRangeStart int64
	CopySourceRangeEnd   int64
	SseHeader            ISseHeader
	SourceSseHeader      ISseHeader
}

// CopyPartOutput is the result of CopyPart function
type CopyPartOutput struct {
	BaseModel
	XMLName      xml.Name   `xml:"CopyPartResult"`
	PartNumber   int        `xml:"-"`
	ETag         string     `xml:"ETag"`
	LastModified time.Time  `xml:"LastModified"`
	SseHeader    ISseHeader `xml:"-"`
}

// CreateSignedUrlInput is the input parameter of CreateSignedUrl function
type CreateSignedUrlInput struct {
	Method      HttpMethodType
	Bucket      string
	Key         string
	SubResource SubResourceType
	Expires     int
	Headers     map[string]string
	QueryParams map[string]string
}

// CreateSignedUrlOutput is the result of CreateSignedUrl function
type CreateSignedUrlOutput struct {
	SignedUrl                  string
	ActualSignedRequestHeaders http.Header
}

// CreateBrowserBasedSignatureInput is the input parameter of CreateBrowserBasedSignature function.
type CreateBrowserBasedSignatureInput struct {
	Bucket     string
	Key        string
	Expires    int
	FormParams map[string]string
}

// CreateBrowserBasedSignatureOutput is the result of CreateBrowserBasedSignature function.
type CreateBrowserBasedSignatureOutput struct {
	OriginPolicy string
	Policy       string
	Algorithm    string
	Credential   string
	Date         string
	Signature    string
}

// HeadObjectInput is the input parameter of HeadObject function
type HeadObjectInput struct {
	Bucket            string
	Key               string
	IfMatch           string
	IfNoneMatch       string
	IfModifiedSince   time.Time
	IfUnmodifiedSince time.Time
	RangeStart        int64
	RangeEnd          int64
}

// UploadFileInput is the input parameter of UploadFile function
type UploadFileInput struct {
	ObjectOperationInput
	ContentType      string
	UploadFile       string
	PartSize         int64
	TaskNum          int
	EnableCheckpoint bool
	CheckpointFile   string
	EncodingType     string
}

// DownloadFileInput is the input parameter of DownloadFile function
type DownloadFileInput struct {
	GetObjectMetadataInput
	IfMatch           string
	IfNoneMatch       string
	IfModifiedSince   time.Time
	IfUnmodifiedSince time.Time
	DownloadFile      string
	PartSize          int64
	TaskNum           int
	EnableCheckpoint  bool
	CheckpointFile    string
}

// SetBucketFetchPolicyInput is the input parameter of SetBucketFetchPolicy function
type SetBucketFetchPolicyInput struct {
	Bucket string
	Status FetchPolicyStatusType `json:"status"`
	Agency string                `json:"agency"`
}

// GetBucketFetchPolicyInput is the input parameter of GetBucketFetchPolicy function
type GetBucketFetchPolicyInput struct {
	Bucket string
}

// GetBucketFetchPolicyOutput is the result of GetBucketFetchPolicy function
type GetBucketFetchPolicyOutput struct {
	BaseModel
	FetchResponse `json:"fetch"`
}

// FetchResponse defines the response fetch policy configuration
type FetchResponse struct {
	Status FetchPolicyStatusType `json:"status"`
	Agency string                `json:"agency"`
}

// DeleteBucketFetchPolicyInput is the input parameter of DeleteBucketFetchPolicy function
type DeleteBucketFetchPolicyInput struct {
	Bucket string
}

// SetBucketFetchJobInput is the input parameter of SetBucketFetchJob function
type SetBucketFetchJobInput struct {
	Bucket           string            `json:"bucket"`
	URL              string            `json:"url"`
	Host             string            `json:"host,omitempty"`
	Key              string            `json:"key,omitempty"`
	Md5              string            `json:"md5,omitempty"`
	CallBackURL      string            `json:"callbackurl,omitempty"`
	CallBackBody     string            `json:"callbackbody,omitempty"`
	CallBackBodyType string            `json:"callbackbodytype,omitempty"`
	CallBackHost     string            `json:"callbackhost,omitempty"`
	FileType         string            `json:"file_type,omitempty"`
	IgnoreSameKey    bool              `json:"ignore_same_key,omitempty"`
	ObjectHeaders    map[string]string `json:"objectheaders,omitempty"`
	Etag             string            `json:"etag,omitempty"`
	TrustName        string            `json:"trustname,omitempty"`
}

// SetBucketFetchJobOutput is the result of SetBucketFetchJob function
type SetBucketFetchJobOutput struct {
	BaseModel
	SetBucketFetchJobResponse
}

// SetBucketFetchJobResponse defines the response SetBucketFetchJob configuration
type SetBucketFetchJobResponse struct {
	ID   string `json:"id"`
	Wait int    `json:"Wait"`
}

// GetBucketFetchJobInput is the input parameter of GetBucketFetchJob function
type GetBucketFetchJobInput struct {
	Bucket string
	JobID  string
}

// GetBucketFetchJobOutput is the result of GetBucketFetchJob function
type GetBucketFetchJobOutput struct {
	BaseModel
	GetBucketFetchJobResponse
}

// GetBucketFetchJobResponse defines the response fetch job configuration
type GetBucketFetchJobResponse struct {
	Err    string      `json:"err"`
	Code   string      `json:"code"`
	Status string      `json:"status"`
	Job    JobResponse `json:"job"`
}

// JobResponse defines the response job configuration
type JobResponse struct {
	Bucket           string `json:"bucket"`
	URL              string `json:"url"`
	Host             string `json:"host"`
	Key              string `json:"key"`
	Md5              string `json:"md5"`
	CallBackURL      string `json:"callbackurl"`
	CallBackBody     string `json:"callbackbody"`
	CallBackBodyType string `json:"callbackbodytype"`
	CallBackHost     string `json:"callbackhost"`
	FileType         string `json:"file_type"`
	IgnoreSameKey    bool   `json:"ignore_same_key"`
}
