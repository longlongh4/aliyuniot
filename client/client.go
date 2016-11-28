package iotclient

import (
	"fmt"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/parnurzeal/gorequest"
)

type ClientInfoStruct struct {
	ProductKey  string `json:"productKey"`
	DeviceID    string `json:"deviceId"`
	DeviceName  string `json:"deviceName"`
	Sign        string `json:"sign"`
	UserNameMD5 string
}

type AuthResponseStruct struct {
	Servers   string `json:"servers"`
	Pubkey    string `json:"pubkey"`
	PkVersion string `json:"pkVersion"`
	DeviceID  string `json:"deviceId"`
	Success   bool   `json:"success"`
	Sign      string `json："sign"`
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}

func (s *AuthResponseStruct) IsSuccess() bool {
	return s.Success
}

func (s *AuthResponseStruct) GetServerAddress() string {
	array := strings.Split(s.Servers, "|")
	if len(array) >= 1 {
		return array[0]
	}
	return ""
}

var clientInfo ClientInfoStruct

func InitClient(newClientInfo ClientInfoStruct) {
	clientInfo = newClientInfo
}

func IotAuth() *AuthResponseStruct {
	response := AuthResponseStruct{}
	fmt.Println(clientInfo)
	_, body, _ := gorequest.New().Post("http://iot.channel.aliyun.com/iot/auth").Type("form").Send(
		map[string]string{
			"deviceName": clientInfo.DeviceName,
			"productKey": clientInfo.ProductKey,
			"sign":       clientInfo.Sign,
		}).EndStruct(&response)
	fmt.Println(string(body))
	return &response
}

func GetMqttClientOpts(serverAddress string) *MQTT.ClientOptions {
	opts := MQTT.NewClientOptions().AddBroker("tcp://" + serverAddress)
	opts.SetClientID(clientInfo.ProductKey + ":" + clientInfo.DeviceID)
	opts.SetUsername(clientInfo.UserNameMD5)
	opts.SetAutoReconnect(true)
	opts.SetKeepAlive(80 * time.Second)
	opts.SetMaxReconnectInterval(5 * time.Second)

	return opts
}
