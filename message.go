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
	title			string	`json:"title"`
	body			string	`json:"body,omitempty"`
	icon			string	`json:"icon,omitempty"`
	sound			string	`json:"sound,omitempty"`
	badge			string	`json:"badge,omitempty"`
	tag				string	`json:"tag,omitempty"`
	color			string	`json:"color,omitempty"`
	clickAction		string	`json:"click_action,omitempty"`
	bodyLocKey 		string	`json:"body_loc_key,omitempty"`
	bodyLocArgs		string	`json:"body_loc_args,omitempty"`
	titleLocKey 	string	`json:"title_loc_key,omitempty"`
	titleLocArgs	string	`json:"title_loc_args,omitempty"`
}

// NewMessage returns a new Message with the specified payload
// and registration IDs.
func NewMessage(data map[string]interface{}, regIDs ...string) *Message {
	return &Message{RegistrationIDs: regIDs, Data: data}
}

func NewMessageWithNotification(data map[string]interface{}, notification Notification, regIDs ...string) *Message {
	return &Message{RegistrationIDs: regIDs, Data: data, Notification: notification}
}