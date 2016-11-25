package iotserver

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/longlongh4/aliyuniot/client"
	"github.com/parnurzeal/gorequest"
)

type BaseResponseStruct struct {
	RequestID string `json:"RequestId"`
	Success   bool   `json:"Success"`
	Code      string `json:"Code"`
	HostID    string `json:"HostId"`
	Message   string `json:"Message"`
}

type DevicePermitsStruct struct {
	DevicePermissions struct {
		DevicePermission []struct {
			GrantType     string `json:"GrantType"`
			ID            int    `json:"Id"`
			TopicFullName string `json:"TopicFullName"`
			TopicUserID   int    `json:"TopicUserId"`
		} `json:"DevicePermission"`
	} `json:"DevicePermissions"`
	BaseResponseStruct
}

type DevicePermitStruct struct {
	ErrorMessage string `json:"ErrorMessage"`
	ID           int    `json:"id"`
	BaseResponseStruct
}

type RegisterDeviceStruct struct {
	DeviceID     string `json:"DeviceId"`
	DeviceName   string `json:"DeviceName"`
	DeviceSecret string `json:"DeviceSecret"`
	DeviceStatus string `json:"DeviceStatus"`
	ErrorMessage string `json:"ErrorMessage"`
	BaseResponseStruct
}

func (s *RegisterDeviceStruct) GetClientInfoStruct(productKey string, productSecret string) iotclient.ClientInfoStruct {
	clientInfo := iotclient.ClientInfoStruct{DeviceID: s.DeviceID, ProductKey: productKey, DeviceName: s.DeviceName}
	mac := hmac.New(md5.New, []byte(productSecret+s.DeviceSecret))
	mac.Write([]byte(fmt.Sprintf("deviceName%sproductKey%s", s.DeviceName, productKey)))
	clientInfo.Sign = strings.ToUpper(hex.EncodeToString(mac.Sum(nil)))

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(productKey + productSecret + s.DeviceID + s.DeviceSecret))
	clientInfo.UserNameMD5 = strings.ToUpper(hex.EncodeToString(md5Ctx.Sum(nil)))
	return clientInfo
}

func GetDevicePermitsRequest(productKey string, deviceName string) *DevicePermitsStruct {
	result := DevicePermitsStruct{}
	request := newBaseRequestParameters()
	request.Set("Action", "ListDevicePermits")
	request.Set("ProductKey", productKey)
	request.Set("DeviceName", deviceName)
	request.Set("Signature", getSignString(request))
	gorequest.New().Post(serverAddress).Type("form").Send(&request).EndStruct(&result)
	return &result
}

func PubMessageToTopicRequest(productKey string, message string, topicFullName string, Qos int) *BaseResponseStruct {
	result := BaseResponseStruct{}
	request := newBaseRequestParameters()
	request.Set("Action", "Pub")
	request.Set("ProductKey", productKey)
	request.Set("MessageContent", base64.StdEncoding.EncodeToString([]byte(message)))
	request.Set("TopicFullName", topicFullName)
	request.Set("Qos", strconv.Itoa(Qos))
	request.Set("Signature", getSignString(request))
	gorequest.New().Post(serverAddress).Type("form").Send(&request).EndStruct(&result)
	return &result
}

func SubTopicRequest(productKey string, subCallback string, topics ...string) *BaseResponseStruct {
	result := BaseResponseStruct{}
	request := newBaseRequestParameters()
	request.Set("Action", "Sub")
	request.Set("ProductKey", productKey)
	request.Set("SubCallback", subCallback)
	for k, v := range topics {
		request.Set("Topic."+strconv.Itoa(k+1), v)
	}
	request.Set("Signature", getSignString(request))
	gorequest.New().Post(serverAddress).Type("form").Send(&request).EndStruct(&result)
	return &result
}

func CancelSubTopicRequest(productKey string, topics ...string) *BaseResponseStruct {
	result := BaseResponseStruct{}
	request := newBaseRequestParameters()
	request.Set("Action", "UnSub")
	request.Set("ProductKey", productKey)
	for k, v := range topics {
		request.Set("Topic."+strconv.Itoa(k+1), v)
	}
	request.Set("Signature", getSignString(request))
	gorequest.New().Post(serverAddress).Type("form").Send(&request).EndStruct(&result)
	return &result
}

func GrantDeviceRequest(productKey string, deviceName string, grantType string, topicFullName string) *DevicePermitStruct {
	result := DevicePermitStruct{}
	request := newBaseRequestParameters()
	request.Set("Action", "DeviceGrant")
	request.Set("ProductKey", productKey)
	request.Set("DeviceName", deviceName)
	request.Set("GrantType", grantType)
	request.Set("TopicFullName", topicFullName)
	request.Set("Signature", getSignString(request))
	gorequest.New().Post(serverAddress).Type("form").Send(&request).EndStruct(&result)
	return &result
}

func ModifyDevicePermitRequest(productKey string, deviceName string, ruleID string, grantType string, topicFullName string) *DevicePermitStruct {
	result := DevicePermitStruct{}
	request := newBaseRequestParameters()
	request.Set("Action", "DevicePermitModify")
	request.Set("ProductKey", productKey)
	request.Set("DeviceName", deviceName)
	request.Set("RuleId", ruleID)
	request.Set("GrantType", grantType)
	request.Set("TopicFullName", topicFullName)
	request.Set("Signature", getSignString(request))
	gorequest.New().Post(serverAddress).Type("form").Send(&request).EndStruct(&result)
	return &result
}

func RemoveDevicePermitByIDRequest(productKey string, deviceName string, ruleID string) *DevicePermitStruct {
	result := DevicePermitStruct{}
	request := newBaseRequestParameters()
	request.Set("Action", "DeviceRevokeById")
	request.Set("ProductKey", productKey)
	request.Set("DeviceName", deviceName)
	request.Set("RuleId", ruleID)
	request.Set("Signature", getSignString(request))
	gorequest.New().Post(serverAddress).Type("form").Send(&request).EndStruct(&result)
	return &result
}

func ModifyDevicePermitByTopicRequest(productKey string, deviceName string, grantType string, topicFullName string) *DevicePermitStruct {
	result := DevicePermitStruct{}
	request := newBaseRequestParameters()
	request.Set("Action", "DeviceRevokeByTopic")
	request.Set("ProductKey", productKey)
	request.Set("DeviceName", deviceName)
	request.Set("GrantType", grantType)
	request.Set("TopicFullName", topicFullName)
	request.Set("Signature", getSignString(request))
	gorequest.New().Post(serverAddress).Type("form").Send(&request).EndStruct(&result)
	return &result
}

func RegisterDeviceRequest(productKey string, deviceName string) *RegisterDeviceStruct {
	result := RegisterDeviceStruct{}
	request := newBaseRequestParameters()
	request.Set("Action", "RegistDevice")
	request.Set("ProductKey", productKey)
	request.Set("DeviceName", deviceName)
	request.Set("Signature", getSignString(request))
	gorequest.New().Post(serverAddress).Type("form").Send(&request).EndStruct(&result)
	return &result
}

func ServerOnlineRequest(productKey string) *BaseResponseStruct {
	result := BaseResponseStruct{}
	request := newBaseRequestParameters()
	request.Set("Action", "ServerOnline")
	request.Set("ProductKey", productKey)
	request.Set("Signature", getSignString(request))
	gorequest.New().Post(serverAddress).Type("form").Send(&request).EndStruct(&result)
	return &result
}
