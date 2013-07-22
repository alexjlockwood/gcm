// Copyright 2013 Alex Lockwood. All rights reserved.

package gcm

// Message is used by the application server to send a message to
// the GCM server. See the documentation for GCM Architectural
// Overview for more information:
// http://developer.android.com/google/gcm/gcm.html#send-msg
type Message struct {
	RegistrationIDs       []string          `json:"registration_ids"`
	CollapseKey           string            `json:"collapse_key,omitempty"`
	Data                  map[string]string `json:"data,omitempty"`
	DelayWhileIdle        bool              `json:"delay_while_idle,omitempty"`
	TimeToLive            int               `json:"time_to_live,omitempty"`
	RestrictedPackageName string            `json:"restricted_package_name,omitempty"`
	DryRun                bool              `json:"dry_run,omitempty"`
}

// NewMessage returns a new Message with the specified payload
// and registration ids.
func NewMessage(data map[string]string, regIds ...string) *Message {
	return &Message{RegistrationIDs: regIds, Data: data}
}
