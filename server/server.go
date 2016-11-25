package iotserver

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const serverAddress = "https://iot.aliyuncs.com"

type ConfigStruct struct {
	AccessKeyID     string
	AccessKeySecret string
}

var config ConfigStruct

func InitServer(newconfig ConfigStruct) {
	config = newconfig
}

func newBaseRequestParameters() *url.Values {
	return &url.Values{
		"Format":           {"json"},
		"Version":          {"2016-05-30"},
		"AccessKeyId":      {config.AccessKeyID},
		"SignatureMethod":  {"HMAC-SHA1"},
		"Timestamp":        {time.Now().UTC().Format("2006-01-02T15:04:05Z")},
		"SignatureVersion": {"1.0"},
		"SignatureNonce":   {strconv.FormatInt(time.Now().UnixNano(), 10)},
		"RegionId":         {"cn-hangzhou"},
	}
}

func getSignString(values *url.Values) string {
	canonicalizedQueryString := percentReplace(values.Encode())
	stringToSign := "POST&%2F&" + url.QueryEscape(canonicalizedQueryString)
	// Crypto by HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, []byte(config.AccessKeySecret+"&"))
	hmacSha1.Write([]byte(stringToSign))
	sign := hmacSha1.Sum(nil)

	// Encode to Base64
	return base64.StdEncoding.EncodeToString(sign)
}

func percentReplace(str string) string {
	str = strings.Replace(str, "+", "%20", -1)
	str = strings.Replace(str, "*", "%2A", -1)
	str = strings.Replace(str, "%7E", "~", -1)

	return str
}
