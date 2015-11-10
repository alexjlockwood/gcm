package gcm

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// Message contains all fields that can be sent to the GCM API.
// Not all fields should be set, it's recommended to read the reference before
// sending anything
// at https://developers.google.com/cloud-messaging/http-server-ref
type Message struct {
	To                    string                 `json:"to"`
	RegistrationIDs       []string               `json:"registration_ids,omitempty"`
	CollapseKey           string                 `json:"collapse_key,omitempty"`
	Priority              string                 `json:"priority,omitempty"`
	ContentAvailable      bool                   `json:"content_available,omitempty"`
	DelayWhileIdle        bool                   `json:"delay_while_idle,omitempty"`
	TimeToLive            int                    `json:"time_to_live,omitempty"`
	RestrictedPackageName string                 `json:"restricted_package_name,omitempty"`
	DryRun                bool                   `json:"dry_run,omitempty"`
	Data                  map[string]interface{} `json:"data,omitempty"`
	Notification          *Notification          `json:"notification,omitempty"`
}

// Notification containts all notification fields as defined in the GCM API reference
type Notification struct {
	Title        string `json:"title,omitempty"`
	Body         string `json:"body,omitempty"`
	Icon         string `json:"icon,omitempty"`
	Sound        string `json:"sound,omitempty"`
	Badge        string `json:"badge,omitempty"`
	Tag          string `json:"tag,omitempty"`
	Color        string `json:"color,omitempty"`
	ClickAction  string `json:"click_action,omitempty"`
	BodyLocKey   string `json:"body_loc_key,omitempty"`
	BodyLocArgs  string `json:"body_loc_args,omitempty"`
	TitleLocKey  string `json:"title_loc_key,omitempty"`
	TitleLocArgs string `json:"title_loc_args,omitempty"`
}

// Request creates a http.Request from a Message
func (m *Message) Request() (*http.Request, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(apiMethod, gcmSendEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

// Validate checks some basic Message errors so invalid data isn't sent
// to GCM
func (m *Message) Validate() error {
	if len(m.RegistrationIDs) > 1000 {
		return errors.New("the message may specify at most 1000 registration IDs")
	}
	if m.TimeToLive < 0 || 2419200 < m.TimeToLive {
		return errors.New("the message's TimeToLive field must be an integer " +
			"between 0 and 2419200 (4 weeks)")
	}
	return nil
}

// update checks for unavailable registration ids and modifies the Message
// so they can be retried
func (m *Message) update(r *Response) {
	var regIDs []string
	for i, result := range r.Results {
		if result.Error == errorUnavailable {
			regIDs = append(regIDs, m.RegistrationIDs[i])
		}
	}
	m.RegistrationIDs = regIDs
}
