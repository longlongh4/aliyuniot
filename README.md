# aliyuniot
阿里云物联网套件sdk go语言版

Aliyun IotHub sdk golang 

## Installing

install Go and run go get:

```
$ go get github.com/longlongh4/aliyuniot/...
```

## How to use

* for server package, run InitServer(newconfig ConfigStruct) first to init
* for client package, run InitClient(newClientInfo ClientInfoStruct) before use

## Notice

For security, we should never store product secret or device secret in client, so the client sdk use a different strategy from the official sdk,
we use the server package to generate all things we need and pass the data to the client

Here are the steps to init the client sdk

1. run RegisterDeviceRequest first(it doesn't matter even if the device is already registerd), then you will get RegisterDeviceStruct
2. run RegisterDeviceStruct.GetClientInfoStruct and get all you need to init a client
3. run InitClient to init the client sdk
4. run IotAuth() to get the server address
5. run ConnectToServer() to connect to the Iot Hub
