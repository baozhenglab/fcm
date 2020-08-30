package fcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// fcm_server_url fcm server url
	fcm_server_url = "https://fcm.googleapis.com/fcm/send"
	// MAX_TTL the default ttl for a notification
	MAX_TTL = 2419200
	// Priority_HIGH notification priority
	Priority_HIGH = "high"
	// Priority_NORMAL notification priority
	Priority_NORMAL = "normal"
	// retry_after_header header name
	retry_after_header = "Retry-After"
	// error_key readable error caching !
	error_key = "error"
)

var (
	// retreyableErrors whether the error is a retryable
	retreyableErrors = map[string]bool{
		"Unavailable":         true,
		"InternalServerError": true,
	}

	// fcmServerUrl for testing purposes
	fcmServerUrl = fcm_server_url
)

// FcmMsg represents fcm request message
type FcmMsg struct {
	Data                  interface{}         `json:"data,omitempty"`
	To                    string              `json:"to,omitempty"`
	RegistrationIds       []string            `json:"registration_ids,omitempty"`
	CollapseKey           string              `json:"collapse_key,omitempty"`
	Priority              string              `json:"priority,omitempty"`
	Notification          NotificationPayload `json:"notification,omitempty"`
	ContentAvailable      bool                `json:"content_available,omitempty"`
	DelayWhileIdle        bool                `json:"delay_while_idle,omitempty"`
	TimeToLive            int                 `json:"time_to_live,omitempty"`
	RestrictedPackageName string              `json:"restricted_package_name,omitempty"`
	DryRun                bool                `json:"dry_run,omitempty"`
	Condition             string              `json:"condition,omitempty"`
	MutableContent        bool                `json:"mutable_content,omitempty"`
	ApiKey                string              `json:"-"`
}

// NotificationPayload notification message payload
type NotificationPayload struct {
	Title            string `json:"title,omitempty"`
	Body             string `json:"body,omitempty"`
	Icon             string `json:"icon,omitempty"`
	Sound            string `json:"sound,omitempty"`
	Badge            string `json:"badge,omitempty"`
	Tag              string `json:"tag,omitempty"`
	Color            string `json:"color,omitempty"`
	ClickAction      string `json:"click_action,omitempty"`
	BodyLocKey       string `json:"body_loc_key,omitempty"`
	BodyLocArgs      string `json:"body_loc_args,omitempty"`
	TitleLocKey      string `json:"title_loc_key,omitempty"`
	TitleLocArgs     string `json:"title_loc_args,omitempty"`
	AndroidChannelID string `json:"android_channel_id,omitempty"`
}

// NewFcmTopicMsg sets the targeted token/topic and the data payload
func (this *FcmMsg) NewFcmTopicMsg(to string, body map[string]string) *FcmMsg {

	this.NewFcmMsgTo(to, body)

	return this
}

// NewFcmMsgTo sets the targeted token/topic and the data payload
func (this *FcmMsg) NewFcmMsgTo(to string, body interface{}) *FcmMsg {
	this.To = to
	this.Data = body

	return this
}

// SetMsgData sets data payload
func (this *FcmMsg) SetMsgData(body interface{}) *FcmMsg {
	this.Data = body
	return this
}

// NewFcmRegIdsMsg gets a list of devices with data payload
func (this *FcmMsg) NewFcmRegIdsMsg(list []string, body interface{}) *FcmMsg {
	this.newDevicesList(list)
	this.Data = body

	return this

}

// newDevicesList init the devices list
func (this *FcmMsg) newDevicesList(list []string) *FcmMsg {
	this.RegistrationIds = make([]string, len(list))
	copy(this.RegistrationIds, list)

	return this

}

// AppendDevices adds more devices/tokens to the Fcm request
func (this *FcmMsg) AppendDevices(list []string) *FcmMsg {

	this.RegistrationIds = append(this.RegistrationIds, list...)

	return this
}

// apiKeyHeader generates the value of the Authorization key
func (this *FcmMsg) apiKeyHeader() string {
	return fmt.Sprintf("key=%v", this.ApiKey)
}

// sendOnce send a single request to fcm
func (this *FcmMsg) sendOnce() (*FcmResponseStatus, error) {

	fcmRespStatus := new(FcmResponseStatus)

	jsonByte, err := this.toJsonByte()
	if err != nil {
		return fcmRespStatus, err
	}

	request, err := http.NewRequest("POST", fcmServerUrl, bytes.NewBuffer(jsonByte))
	request.Header.Set("Authorization", this.apiKeyHeader())
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return fcmRespStatus, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fcmRespStatus, err
	}

	fcmRespStatus.StatusCode = response.StatusCode

	fcmRespStatus.RetryAfter = response.Header.Get(retry_after_header)

	if response.StatusCode != 200 {
		return fcmRespStatus, nil
	}

	err = fcmRespStatus.parseStatusBody(body)
	if err != nil {
		return fcmRespStatus, err
	}
	fcmRespStatus.Ok = true

	return fcmRespStatus, nil
}

// Send to fcm
func (this *FcmMsg) Send() (*FcmResponseStatus, error) {
	return this.sendOnce()

}

// toJsonByte converts FcmMsg to a json byte
func (this *FcmMsg) toJsonByte() ([]byte, error) {

	return json.Marshal(this)

}

// parseStatusBody parse FCM response body
func (this *FcmResponseStatus) parseStatusBody(body []byte) error {

	if err := json.Unmarshal([]byte(body), &this); err != nil {
		return err
	}
	return nil

}

// SetPriority Sets the priority of the message.
// Priority_HIGH or Priority_NORMAL
func (this *FcmMsg) SetPriority(p string) *FcmMsg {

	if p == Priority_HIGH {
		this.Priority = Priority_HIGH
	} else {
		this.Priority = Priority_NORMAL
	}

	return this
}

// SetCollapseKey This parameter identifies a group of messages
// (e.g., with collapse_key: "Updates Available") that can be collapsed,
// so that only the last message gets sent when delivery can be resumed.
// This is intended to avoid sending too many of the same messages when the
// device comes back online or becomes active (see delay_while_idle).
func (this *FcmMsg) SetCollapseKey(val string) *FcmMsg {

	this.CollapseKey = val

	return this
}

// SetNotificationPayload sets the notification payload based on the specs
// https://firebase.google.com/docs/cloud-messaging/http-server-ref
func (this *FcmMsg) SetNotificationPayload(payload *NotificationPayload) *FcmMsg {

	this.Notification = *payload

	return this
}

// SetContentAvailable On iOS, use this field to represent content-available
// in the APNS payload. When a notification or message is sent and this is set
// to true, an inactive client app is awoken. On Android, data messages wake
// the app by default. On Chrome, currently not supported.
func (this *FcmMsg) SetContentAvailable(isContentAvailable bool) *FcmMsg {

	this.ContentAvailable = isContentAvailable

	return this
}

// SetDelayWhileIdle When this parameter is set to true, it indicates that
// the message should not be sent until the device becomes active.
// The default value is false.
func (this *FcmMsg) SetDelayWhileIdle(isDelayWhileIdle bool) *FcmMsg {

	this.DelayWhileIdle = isDelayWhileIdle

	return this
}

// SetTimeToLive This parameter specifies how long (in seconds) the message
// should be kept in FCM storage if the device is offline. The maximum time
// to live supported is 4 weeks, and the default value is 4 weeks.
// For more information, see
// https://firebase.google.com/docs/cloud-messaging/concept-options#ttl
func (this *FcmMsg) SetTimeToLive(ttl int) *FcmMsg {

	if ttl > MAX_TTL {

		this.TimeToLive = MAX_TTL

	} else {

		this.TimeToLive = ttl

	}
	return this
}

// SetRestrictedPackageName This parameter specifies the package name of the
// application where the registration tokens must match in order to
// receive the message.
func (this *FcmMsg) SetRestrictedPackageName(pkg string) *FcmMsg {

	this.RestrictedPackageName = pkg

	return this
}

// SetDryRun This parameter, when set to true, allows developers to test
// a request without actually sending a message.
// The default value is false
func (this *FcmMsg) SetDryRun(drun bool) *FcmMsg {

	this.DryRun = drun

	return this
}

// SetMutableContent Currently for iOS 10+ devices only. On iOS,
// use this field to represent mutable-content in the APNs payload.
// When a notification is sent and this is set to true, the content
// of the notification can be modified before it is displayed,
// using a Notification Service app extension.
// This parameter will be ignored for Android and web.
func (this *FcmMsg) SetMutableContent(mc bool) *FcmMsg {

	this.MutableContent = mc

	return this
}

// PrintResults prints the FcmResponseStatus results for fast using and debugging
func (this *FcmResponseStatus) PrintResults() {
	fmt.Println("Status Code   :", this.StatusCode)
	fmt.Println("Success       :", this.Success)
	fmt.Println("Fail          :", this.Fail)
	fmt.Println("Canonical_ids :", this.Canonical_ids)
	fmt.Println("Topic MsgId   :", this.MsgId)
	fmt.Println("Topic Err     :", this.Err)
	for i, val := range this.Results {
		fmt.Printf("Result(%d)> \n", i)
		for k, v := range val {
			fmt.Println("\t", k, " : ", v)
		}
	}
}

// IsTimeout check whether the response timeout based on http response status
// code and if any error is retryable
func (this *FcmResponseStatus) IsTimeout() bool {
	if this.StatusCode >= 500 {
		return true
	} else if this.StatusCode == 200 {
		for _, val := range this.Results {
			for k, v := range val {
				if k == error_key && retreyableErrors[v] == true {
					return true
				}
			}
		}
	}

	return false
}

// GetRetryAfterTime converts the retrey after response header
// to a time.Duration
func (this *FcmResponseStatus) GetRetryAfterTime() (t time.Duration, e error) {
	t, e = time.ParseDuration(this.RetryAfter)
	return
}

// SetCondition to set a logical expression of conditions that determine the message target
func (this *FcmMsg) SetCondition(condition string) *FcmMsg {
	this.Condition = condition
	return this
}
