package gcm

import "errors"

type errorString string

func (e errorString) Error() string {
	return string(e)
}

var (
	errorMap = map[int]error{
		400: errors.New("Invalid JSON"),
		401: errors.New("Authentication Error"),
	}
)

const (
	errorUnavailable = "Unavailable"
)

// Response represents the GCM server's response to the application
// server's sent message. See the documentation for more information:
// https://developers.google.com/cloud-messaging/http-server-ref#interpret-downstream
type Response struct {
	MessageID    int         `json:"message_id,omitempty"`
	Error        errorString `json:"error,omitempty"`
	MulticastID  int         `json:"multicast_id,omitempty"`
	Success      int         `json:"success,omitempty"`
	Failure      int         `json:"failure,omitempty"`
	CanonicalIDs int         `json:"canonical_ids,omitempty"`
	Results      []Result    `json:"results,omitempty"`
}

// Result represents the status of a processed message.
type Result struct {
	MessageID      string      `json:"message_id,omitempty"`
	RegistrationID string      `json:"registration_id,omitempty"`
	Error          errorString `json:"error,omitempty"`
}

// Ok checks if the response contains any failures, canonical id changes or
// topic errors
func (r Response) Ok() bool {
	return r.Failure == 0 && r.CanonicalIDs == 0 && r.Error == ""
}
