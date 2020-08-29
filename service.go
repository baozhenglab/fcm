package fcm

// FcmClient stores the key and the Message (FcmMsg)
type FcmClient struct {
	ApiKey  string
	Message FcmMsg
}

// NewFcmClient init and create fcm client
func NewFcmClient(apiKey string) *FcmClient {
	fcmc := new(FcmClient)
	fcmc.ApiKey = apiKey

	return fcmc
}
