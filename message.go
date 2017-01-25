package gcm

// Priority represents the priority at which to send the message.
// Google only supports two values NormalPriority (the default, and
// equivalent to APNS level 5) and HighPriority (APNS level 10)
type Priority string

// HighPrioty is used to send messages that are urgent. For example
// when a new email arrives.
const HighPriority = Priority("high")

// NormalPriority is used to send messages in a manner that will try
// to optimize battery life over immediate delivery.
const NormalPriority = Priority("normal")

// Message is used by the application server to send a message to
// the GCM server. See the documentation for GCM Architectural
// Overview for more information:
// http://developer.android.com/google/gcm/gcm.html#send-msg
type Message struct {
	RegistrationIDs       []string               `json:"registration_ids"`
	CollapseKey           string                 `json:"collapse_key,omitempty"`
	ContentAvailable      bool                   `json:"content_available,omitempty"`
	Data                  map[string]interface{} `json:"data,omitempty"`
	DelayWhileIdle        bool                   `json:"delay_while_idle,omitempty"`
	TimeToLive            int                    `json:"time_to_live,omitempty"`
	RestrictedPackageName string                 `json:"restricted_package_name,omitempty"`
	DryRun                bool                   `json:"dry_run,omitempty"`
	Priority              Priority               `json:"priority,omitempty"`
	Notification          map[string]interface{} `json:"notification,omitempty"`
}

// NewMessage returns a new Message with the specified payload
// and registration IDs.
func NewMessage(data map[string]interface{}, regIDs ...string) *Message {
	return &Message{RegistrationIDs: regIDs, Data: data}
}
