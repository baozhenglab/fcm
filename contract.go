package fcm

const (
	KeyService = "fcm"
)

type FcmService interface {
	CreateMessage() *FcmMsg
	GetInfo(withDetails bool, instanceIdToken string) (*InstanceIdInfoResponse, error)
	SubscribeToTopic(instanceIdToken string, topic string) (*SubscribeResponse, error)
	BatchSubscribeToTopic(tokens []string, topic string) (*BatchResponse, error)
	BatchUnsubscribeFromTopic(tokens []string, topic string) (*BatchResponse, error)
	ApnsBatchImportRequest(apnsReq *ApnsBatchRequest) (*ApnsBatchResponse, error)
}
