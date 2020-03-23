package fcm

// Message is used by the application server to send a message to
// the FCM server. See the documentation for FCM Architectural
// Overview for more information:
// https://firebase.google.com/docs/cloud-messaging/http-server-ref
type Message struct {
	RegistrationIDs       []string               `json:"registration_ids"`
	CollapseKey           string                 `json:"collapse_key,omitempty"`
	Data                  map[string]interface{} `json:"data,omitempty"`
	DelayWhileIdle        bool                   `json:"delay_while_idle,omitempty"`
	TimeToLive            int                    `json:"time_to_live,omitempty"`
	RestrictedPackageName string                 `json:"restricted_package_name,omitempty"`
	DryRun                bool                   `json:"dry_run,omitempty"`
	Notification          *Notification           `json:"notification"`
}

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// NewMessage returns a new Message with the specified payload
// and registration IDs.
// @DEPRECATED as no validation here and client should create itself freely
func NewMessage(data map[string]interface{}, regIDs ...string) *Message {
	return &Message{
		RegistrationIDs: regIDs,
		Data: data,
		Notification: &Notification{
			Title: data["title"].(string),
			Body: data["message"].(string),
		},
	}
}
