package gcm

// Message is used by the application server to send a message to
// the GCM server. See the documentation for GCM Architectural
// Overview for more information:
// http://developer.android.com/google/gcm/gcm.html#send-msg
type Message struct {
	RegistrationIDs       []string               `json:"registration_ids"`
	CollapseKey           string                 `json:"collapse_key,omitempty"`
	Data                  map[string]interface{} `json:"data,omitempty"`
	Notification          Notification           `json:"notification,omitempty"`
	DelayWhileIdle        bool                   `json:"delay_while_idle,omitempty"`
	TimeToLive            int                    `json:"time_to_live,omitempty"`
	RestrictedPackageName string                 `json:"restricted_package_name,omitempty"`
	DryRun                bool                   `json:"dry_run,omitempty"`
}

type Notification struct {
	Title			string	`json:"title"`
	Body			string	`json:"body,omitempty"`
	Icon			string	`json:"icon,omitempty"`
	Sound			string	`json:"sound,omitempty"`
	Badge			string	`json:"badge,omitempty"`
	Tag				string	`json:"tag,omitempty"`
	Color			string	`json:"color,omitempty"`
	ClickAction		string	`json:"click_action,omitempty"`
	BodyLocKey 		string	`json:"body_loc_key,omitempty"`
	BodyLocArgs		string	`json:"body_loc_args,omitempty"`
	TitleLocKey 	string	`json:"title_loc_key,omitempty"`
	TitleLocArgs	string	`json:"title_loc_args,omitempty"`
}

// NewMessage returns a new Message with the specified payload
// and registration IDs.
func NewMessage(data map[string]interface{}, regIDs ...string) *Message {
	return &Message{RegistrationIDs: regIDs, Data: data}
}

func NewMessageWithNotification(data map[string]interface{}, notification Notification, regIDs ...string) *Message {
	return &Message{RegistrationIDs: regIDs, Data: data, Notification: notification}
}