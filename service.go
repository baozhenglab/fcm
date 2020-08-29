package fcm

import (
	"flag"
	"fmt"

	goservice "github.com/baozhenglab/go-sdk"
)

// FcmClient stores the key and the Message (FcmMsg)
type FcmClient struct {
	ApiKey  string
	Message FcmMsg
}

// NewFcmClient init and create fcm client
func NewFcmClient() goservice.PrefixConfigure {
	return new(FcmClient)
}

func (fcm *FcmClient) Name() string {
	return KeyService
}

func (fcm *FcmClient) GetPrefix() string {
	return KeyService
}

func (fcm *FcmClient) InitFlags() {
	prefix := fmt.Sprintf("%s-", fcm.Name())
	flag.StringVar(&fcm.ApiKey, prefix+"api-key", "", "API key server for Firebase cloud message")
}

func (fcm *FcmClient) Get() interface{} {
	return fcm
}
