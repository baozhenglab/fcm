package fcm

import (
	"flag"
	"fmt"

	goservice "github.com/baozhenglab/go-sdk"
)

// FcmMsg represents fcm response message - (tokens and topics)
type FcmResponseStatus struct {
	Ok            bool
	StatusCode    int
	MulticastId   int64               `json:"multicast_id"`
	Success       int                 `json:"success"`
	Fail          int                 `json:"failure"`
	Canonical_ids int                 `json:"canonical_ids"`
	Results       []map[string]string `json:"results,omitempty"`
	MsgId         int64               `json:"message_id,omitempty"`
	Err           string              `json:"error,omitempty"`
	RetryAfter    string
}

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

func (fcm *FcmClient) CreateMessage() *FcmMsg {
	return &FcmMsg{
		ApiKey: fcm.ApiKey,
	}
}

func (fcm *FcmClient) Get() interface{} {
	return fcm
}
