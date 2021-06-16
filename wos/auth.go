package wos

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

func (wosClient WosClient) doAuthTemporary(method, bucketName, objectKey string, params map[string]string,
	headers map[string][]string, expires int64) (requestURL string, err error) {
	sh := wosClient.getSecurity()
	isAkSkEmpty := sh.ak == "" || sh.sk == ""
	requestURL, canonicalizedURL := wosClient.conf.formatUrls(bucketName, objectKey, params, true)
	parsedRequestURL, err := url.Parse(requestURL)
	if err != nil {
		return "", err
	}
	encodeHeaders(headers)
	hostName := parsedRequestURL.Host

	isV4 := wosClient.conf.signature == SignatureV4
	isWos := wosClient.conf.signature == SignatureWos
	isV2 := wosClient.conf.signature == SignatureV2
	prepareHostAndDate(headers, hostName, isV4, isWos)

	if isAkSkEmpty {
		doLog(LEVEL_WARN, "No ak/sk provided, skip to construct authorization")
	} else {

		if isV2 {
			originDate := headers[HEADER_DATE_CAMEL][0]
			date, parseDateErr := time.Parse(RFC1123_FORMAT, originDate)
			if parseDateErr != nil {
				doLog(LEVEL_WARN, "Failed to parse date with reason: %v", parseDateErr)
				return "", parseDateErr
			}
			expires += date.Unix()
			headers[HEADER_DATE_CAMEL] = []string{Int64ToString(expires)}

			stringToSign := getV2StringToSign(method, canonicalizedURL, headers)
			signature := UrlEncode(Base64Encode(HmacSha1([]byte(sh.sk), []byte(stringToSign))), false)
			if strings.Index(requestURL, "?") < 0 {
				requestURL += "?"
			} else {
				requestURL += "&"
			}
			delete(headers, HEADER_DATE_CAMEL)

			if wosClient.conf.signature != SignatureWos {
				requestURL += "AWS"
			}
			requestURL += fmt.Sprintf("AccessKeyId=%s&Expires=%d&Signature=%s", UrlEncode(sh.ak, false), expires, signature)
		} else {
			date, parseDateErr := time.Parse(LONG_DATE_FORMAT, headers[HEADER_DATE_CAMEL][0])
			if parseDateErr != nil {
				doLog(LEVEL_WARN, "Failed to parse date with reason: %v", parseDateErr)
				return "", parseDateErr
			}
			delete(headers, HEADER_DATE_CAMEL)
			shortDate := date.Format(SHORT_DATE_FORMAT)
			longDate := date.Format(LONG_DATE_FORMAT)
			if len(headers[HEADER_HOST_CAMEL]) != 0 {
				index := strings.LastIndex(headers[HEADER_HOST_CAMEL][0], ":")
				if index != -1 {
					port := headers[HEADER_HOST_CAMEL][0][index+1:]
					if port == "80" || port == "443" {
						headers[HEADER_HOST_CAMEL] = []string{headers[HEADER_HOST_CAMEL][0][:index]}
					}
				}

			}

			signedHeaders, _headers := getSignedHeaders(headers)

			credential, scope := getCredential(sh.ak, wosClient.conf.region, shortDate, isWos)

			if isWos {
				params[PARAM_ALGORITHM_WOS_CAMEL] = V4_WOS_HASH_PREFIX
				params[PARAM_CREDENTIAL_WOS_CAMEL] = credential
				params[PARAM_DATE_WOS_CAMEL] = longDate
				params[PARAM_EXPIRES_WOS_CAMEL] = Int64ToString(expires)
				params[PARAM_SIGNEDHEADERS_WOS_CAMEL] = strings.Join(signedHeaders, ";")
			} else {
				params[PARAM_ALGORITHM_AMZ_CAMEL] = V4_HASH_PREFIX
				params[PARAM_CREDENTIAL_AMZ_CAMEL] = credential
				params[PARAM_DATE_AMZ_CAMEL] = longDate
				params[PARAM_EXPIRES_AMZ_CAMEL] = Int64ToString(expires)
				params[PARAM_SIGNEDHEADERS_AMZ_CAMEL] = strings.Join(signedHeaders, ";")
			}

			requestURL, canonicalizedURL = wosClient.conf.formatUrls(bucketName, objectKey, params, true)
			parsedRequestURL, _err := url.Parse(requestURL)
			if _err != nil {
				return "", _err
			}

			stringToSign := getV4StringToSign(method, canonicalizedURL, parsedRequestURL.RawQuery, scope, longDate, UNSIGNED_PAYLOAD, signedHeaders, _headers, isWos)
			signature := getSignature(stringToSign, sh.sk, wosClient.conf.region, shortDate, isWos)

			if isWos {
				requestURL += fmt.Sprintf("&%s=%s", PARAM_SIGNATURE_WOS_CAMEL, UrlEncode(signature, false))
			} else {
				requestURL += fmt.Sprintf("&%s=%s", PARAM_SIGNATURE_AMZ_CAMEL, UrlEncode(signature, false))
			}
		}
	}

	return
}

func (wosClient WosClient) doAuth(method, bucketName, objectKey string, params map[string]string,
	headers map[string][]string, hostName string) (requestURL string, err error) {
	sh := wosClient.getSecurity()
	isAkSkEmpty := sh.ak == "" || sh.sk == ""
	isWos := wosClient.conf.signature == SignatureWos
	requestURL, canonicalizedURL := wosClient.conf.formatUrls(bucketName, objectKey, params, true)
	parsedRequestURL, err := url.Parse(requestURL)
	if err != nil {
		return "", err
	}
	encodeHeaders(headers)

	if hostName == "" {
		hostName = parsedRequestURL.Host
	}

	isV4 := wosClient.conf.signature == SignatureV4
	isV2 := wosClient.conf.signature == SignatureV2
	prepareHostAndDate(headers, hostName, isV4, isWos)

	if isAkSkEmpty {
		doLog(LEVEL_WARN, "No ak/sk provided, skip to construct authorization")
	} else {
		ak := sh.ak
		sk := sh.sk
		var authorization string

		if isV2 {
			ret := v2Auth(ak, sk, method, canonicalizedURL, headers)
			hashPrefix := V2_HASH_PREFIX
			authorization = fmt.Sprintf("%s %s:%s", hashPrefix, ak, ret["Signature"])
		} else {
			if isWos {
				headers[HEADER_CONTENT_SHA256_WOS] = []string{UNSIGNED_PAYLOAD}
			} else {
				headers[HEADER_CONTENT_SHA256_AMZ] = []string{UNSIGNED_PAYLOAD}
			}
			ret := v4Auth(ak, sk, wosClient.conf.region, method, canonicalizedURL, parsedRequestURL.RawQuery, headers, isWos)
			if isWos {
				authorization = fmt.Sprintf("%s Credential=%s,SignedHeaders=%s,Signature=%s", V4_WOS_HASH_PREFIX, ret["Credential"], ret["SignedHeaders"], ret["Signature"])
			} else {
				authorization = fmt.Sprintf("%s Credential=%s,SignedHeaders=%s,Signature=%s", V4_HASH_PREFIX, ret["Credential"], ret["SignedHeaders"], ret["Signature"])
			}

		}
		headers[HEADER_AUTH_CAMEL] = []string{authorization}
	}
	return
}

func prepareHostAndDate(headers map[string][]string, hostName string, isV4 bool, isWos bool) {
	headers[HEADER_HOST_CAMEL] = []string{hostName}
	var (
		date []string
		ok   bool
	)

	if isWos {
		date, ok = headers[HEADER_DATE_WOS]
	} else {
		date, ok = headers[HEADER_DATE_AMZ]
	}
	if ok {
		flag := false
		if len(date) == 1 {
			if isV4 || isWos {
				if t, err := time.Parse(LONG_DATE_FORMAT, date[0]); err == nil {
					headers[HEADER_DATE_CAMEL] = []string{t.Format(LONG_DATE_FORMAT)}
					flag = true
				}
			} else {
				if strings.HasSuffix(date[0], "GMT") {
					headers[HEADER_DATE_CAMEL] = []string{date[0]}
					flag = true
				}
			}
		}
		if !flag {
			delete(headers, HEADER_DATE_AMZ)
		}
	}
	if _, ok := headers[HEADER_DATE_CAMEL]; !ok {
		if isV4 || isWos {
			headers[HEADER_DATE_CAMEL] = []string{time.Now().UTC().Format(LONG_DATE_FORMAT)}
		} else {
			headers[HEADER_DATE_CAMEL] = []string{FormatUtcToRfc1123(time.Now().UTC())}
		}
	}

}

func encodeHeaders(headers map[string][]string) {
	for key, values := range headers {
		for index, value := range values {
			values[index] = UrlEncode(value, true)
		}
		headers[key] = values
	}
}

func prepareDateHeader(dataHeader, dateCamelHeader string, headers, _headers map[string][]string) {
	if _, ok := _headers[HEADER_DATE_CAMEL]; ok {
		if _, ok := _headers[dataHeader]; ok {
			_headers[HEADER_DATE_CAMEL] = []string{""}
		} else if _, ok := headers[dateCamelHeader]; ok {
			_headers[HEADER_DATE_CAMEL] = []string{""}
		}
	} else if _, ok := _headers[strings.ToLower(HEADER_DATE_CAMEL)]; ok {
		if _, ok := _headers[dataHeader]; ok {
			_headers[HEADER_DATE_CAMEL] = []string{""}
		} else if _, ok := headers[dateCamelHeader]; ok {
			_headers[HEADER_DATE_CAMEL] = []string{""}
		}
	}
}

func getStringToSign(keys []string, _headers map[string][]string) []string {
	stringToSign := make([]string, 0, len(keys))
	for _, key := range keys {
		var value string
		prefixHeader := HEADER_PREFIX
		prefixMetaHeader := HEADER_PREFIX_META
		if strings.HasPrefix(key, prefixHeader) {
			if strings.HasPrefix(key, prefixMetaHeader) {
				for index, v := range _headers[key] {
					value += strings.TrimSpace(v)
					if index != len(_headers[key])-1 {
						value += ","
					}
				}
			} else {
				value = strings.Join(_headers[key], ",")
			}
			value = fmt.Sprintf("%s:%s", key, value)
		} else {
			value = strings.Join(_headers[key], ",")
		}
		stringToSign = append(stringToSign, value)
	}
	return stringToSign
}

func attachHeaders(headers map[string][]string) string {
	length := len(headers)
	_headers := make(map[string][]string, length)
	keys := make([]string, 0, length)

	for key, value := range headers {
		_key := strings.ToLower(strings.TrimSpace(key))
		if _key != "" {
			prefixheader := HEADER_PREFIX
			if _key == "content-md5" || _key == "content-type" || _key == "date" || strings.HasPrefix(_key, prefixheader) {
				keys = append(keys, _key)
				_headers[_key] = value
			}
		} else {
			delete(headers, key)
		}
	}

	for _, interestedHeader := range interestedHeaders {
		if _, ok := _headers[interestedHeader]; !ok {
			_headers[interestedHeader] = []string{""}
			keys = append(keys, interestedHeader)
		}
	}
	dateCamelHeader := PARAM_DATE_AMZ_CAMEL
	dataHeader := HEADER_DATE_AMZ
	prepareDateHeader(dataHeader, dateCamelHeader, headers, _headers)

	sort.Strings(keys)
	stringToSign := getStringToSign(keys, _headers)
	return strings.Join(stringToSign, "\n")
}

func getV2StringToSign(method, canonicalizedURL string, headers map[string][]string) string {
	stringToSign := strings.Join([]string{method, "\n", attachHeaders(headers), "\n", canonicalizedURL}, "")

	var isSecurityToken bool
	var securityToken []string
	securityToken, isSecurityToken = headers[HEADER_STS_TOKEN_AMZ]
	var query []string
	if !isSecurityToken {
		parmas := strings.Split(canonicalizedURL, "?")
		if len(parmas) > 1 {
			query = strings.Split(parmas[1], "&")
			for _, value := range query {
				if strings.HasPrefix(value, HEADER_STS_TOKEN_AMZ+"=") || strings.HasPrefix(value, HEADER_STS_TOKEN_WOS+"=") {
					if value[len(HEADER_STS_TOKEN_AMZ)+1:] != "" {
						securityToken = []string{value[len(HEADER_STS_TOKEN_AMZ)+1:]}
						isSecurityToken = true
					}
				}
			}
		}
	}
	logStringToSign := stringToSign
	if isSecurityToken && len(securityToken) > 0 {
		logStringToSign = strings.Replace(logStringToSign, securityToken[0], "******", -1)
	}
	doLog(LEVEL_DEBUG, "The v2 auth stringToSign:\n%s", logStringToSign)
	return stringToSign
}

func v2Auth(ak, sk, method, canonicalizedURL string, headers map[string][]string) map[string]string {
	stringToSign := getV2StringToSign(method, canonicalizedURL, headers)
	return map[string]string{"Signature": Base64Encode(HmacSha1([]byte(sk), []byte(stringToSign)))}
}

func getScope(region, shortDate string, isWos bool) string {
	var scope string
	if isWos {
		scope = fmt.Sprintf("%s/%s/%s/%s", shortDate, region, V4_WOS_SERVICE_NAME, V4_WOS_SERVICE_SUFFIX)
	} else {
		scope = fmt.Sprintf("%s/%s/%s/%s", shortDate, region, V4_SERVICE_NAME, V4_SERVICE_SUFFIX)
	}
	return scope
}

func getCredential(ak, region, shortDate string, isWos bool) (string, string) {
	scope := getScope(region, shortDate, isWos)
	return fmt.Sprintf("%s/%s", ak, scope), scope
}

func getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, payload string, signedHeaders []string, headers map[string][]string, isWos bool) string {
	canonicalRequest := make([]string, 0, 10+len(signedHeaders)*4)
	canonicalRequest = append(canonicalRequest, method)
	canonicalRequest = append(canonicalRequest, "\n")
	canonicalRequest = append(canonicalRequest, canonicalizedURL)
	canonicalRequest = append(canonicalRequest, "\n")
	canonicalRequest = append(canonicalRequest, queryURL)
	canonicalRequest = append(canonicalRequest, "\n")

	for _, signedHeader := range signedHeaders {
		values, _ := headers[signedHeader]
		for _, value := range values {
			canonicalRequest = append(canonicalRequest, signedHeader)
			canonicalRequest = append(canonicalRequest, ":")
			canonicalRequest = append(canonicalRequest, value)
			canonicalRequest = append(canonicalRequest, "\n")
		}
	}
	canonicalRequest = append(canonicalRequest, "\n")
	canonicalRequest = append(canonicalRequest, strings.Join(signedHeaders, ";"))
	canonicalRequest = append(canonicalRequest, "\n")
	canonicalRequest = append(canonicalRequest, payload)

	_canonicalRequest := strings.Join(canonicalRequest, "")

	var isSecurityToken bool
	var securityToken []string
	if securityToken, isSecurityToken = headers[HEADER_STS_TOKEN_WOS]; !isSecurityToken {
		securityToken, isSecurityToken = headers[HEADER_STS_TOKEN_AMZ]
	}
	var query []string
	if !isSecurityToken {
		query = strings.Split(queryURL, "&")
		for _, value := range query {
			if strings.HasPrefix(value, HEADER_STS_TOKEN_AMZ+"=") || strings.HasPrefix(value, HEADER_STS_TOKEN_WOS+"=") {
				if value[len(HEADER_STS_TOKEN_AMZ)+1:] != "" {
					securityToken = []string{value[len(HEADER_STS_TOKEN_AMZ)+1:]}
					isSecurityToken = true
				}
			}
		}
	}
	logCanonicalRequest := _canonicalRequest
	if isSecurityToken && len(securityToken) > 0 {
		logCanonicalRequest = strings.Replace(logCanonicalRequest, securityToken[0], "******", -1)
	}
	doLog(LEVEL_DEBUG, "The v4 auth canonicalRequest:\n%s", logCanonicalRequest)

	stringToSign := make([]string, 0, 7)

	if isWos {
		stringToSign = append(stringToSign, V4_WOS_HASH_PREFIX)
	} else {
		stringToSign = append(stringToSign, V4_HASH_PREFIX)
	}
	stringToSign = append(stringToSign, "\n")
	stringToSign = append(stringToSign, longDate)
	stringToSign = append(stringToSign, "\n")
	stringToSign = append(stringToSign, scope)
	stringToSign = append(stringToSign, "\n")
	stringToSign = append(stringToSign, HexSha256([]byte(_canonicalRequest)))

	_stringToSign := strings.Join(stringToSign, "")

	doLog(LEVEL_DEBUG, "The v4 auth stringToSign:\n%s", _stringToSign)
	return _stringToSign
}

func getSignedHeaders(headers map[string][]string) ([]string, map[string][]string) {
	length := len(headers)
	_headers := make(map[string][]string, length)
	signedHeaders := make([]string, 0, length)
	for key, value := range headers {
		_key := strings.ToLower(strings.TrimSpace(key))
		if _key != "" {
			signedHeaders = append(signedHeaders, _key)
			_headers[_key] = value
		} else {
			delete(headers, key)
		}
	}
	sort.Strings(signedHeaders)
	return signedHeaders, _headers
}

func getSignature(stringToSign, sk, region, shortDate string, isWos bool) string {
	var signature string
	if isWos {
		key := HmacSha256([]byte(V4_WOS_HASH_PRE+sk), []byte(shortDate))
		key = HmacSha256(key, []byte(region))
		key = HmacSha256(key, []byte(V4_WOS_SERVICE_NAME))
		key = HmacSha256(key, []byte(V4_WOS_SERVICE_SUFFIX))
		signature = Hex(HmacSha256(key, []byte(stringToSign)))
	} else {
		key := HmacSha256([]byte(V4_HASH_PRE+sk), []byte(shortDate))
		key = HmacSha256(key, []byte(region))
		key = HmacSha256(key, []byte(V4_SERVICE_NAME))
		key = HmacSha256(key, []byte(V4_SERVICE_SUFFIX))
		signature = Hex(HmacSha256(key, []byte(stringToSign)))
	}

	return signature
}

// V4Auth is a wrapper for v4Auth
func V4Auth(ak, sk, region, method, canonicalizedURL, queryURL string, headers map[string][]string) map[string]string {
	return v4Auth(ak, sk, region, method, canonicalizedURL, queryURL, headers, false)
}

func v4Auth(ak, sk, region, method, canonicalizedURL, queryURL string, headers map[string][]string, isWos bool) map[string]string {
	var t time.Time
	var headDate string
	if isWos {
		headDate = HEADER_DATE_WOS
	} else {
		headDate = HEADER_DATE_AMZ
	}

	if val, ok := headers[headDate]; ok {
		var err error
		t, err = time.Parse(LONG_DATE_FORMAT, val[0])
		if err != nil {
			t = time.Now().UTC()
		}
	} else if val, ok := headers[PARAM_DATE_AMZ_CAMEL]; ok {
		var err error
		t, err = time.Parse(LONG_DATE_FORMAT, val[0])
		if err != nil {
			t = time.Now().UTC()
		}
	} else if val, ok := headers[HEADER_DATE_CAMEL]; ok {
		var err error
		t, err = time.Parse(LONG_DATE_FORMAT, val[0])
		if err != nil {
			t = time.Now().UTC()
		}
	} else if val, ok := headers[strings.ToLower(HEADER_DATE_CAMEL)]; ok {
		var err error
		t, err = time.Parse(RFC1123_FORMAT, val[0])
		if err != nil {
			t = time.Now().UTC()
		}
	} else {
		t = time.Now().UTC()
	}

	shortDate := t.Format(SHORT_DATE_FORMAT)
	longDate := t.Format(LONG_DATE_FORMAT)

	signedHeaders, _headers := getSignedHeaders(headers)

	credential, scope := getCredential(ak, region, shortDate, isWos)

	payload := UNSIGNED_PAYLOAD
	if val, ok := headers[HEADER_CONTENT_SHA256_AMZ]; ok {
		payload = val[0]
	}
	stringToSign := getV4StringToSign(method, canonicalizedURL, queryURL, scope, longDate, payload, signedHeaders, _headers, isWos)

	signature := getSignature(stringToSign, sk, region, shortDate, isWos)

	ret := make(map[string]string, 3)
	ret["Credential"] = credential
	ret["SignedHeaders"] = strings.Join(signedHeaders, ";")
	ret["Signature"] = signature
	return ret
}
