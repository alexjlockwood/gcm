// Package gcm is Google Cloud Messaging for application servers implemented using the
// Go programming language.
package gcm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	// GcmSendEndpoint is the endpoint for sending messages to the GCM server.
	GcmSendEndpoint = "https://gcm-http.googleapis.com/gcm/send"
	// Initial delay before first retry, without jitter.
	backoffInitialDelay = 1000
	// Maximum delay before a retry.
	maxBackoffDelay = 1024000
	// Http method for the api
	apiMethod = "POST"
)

// Declared as a mutable variable for testing purposes.
var gcmSendEndpoint = GcmSendEndpoint

// Sender abstracts the interaction between the application server and the
// GCM server. The developer must obtain an API key from the Google APIs
// Console page and pass it to the Sender so that it can perform authorized
// requests on the application server's behalf. To send a message to one or
// more devices use the Sender's Send method.
//
// If your application server runs on Google AppEngine,
// you must use the "appengine/urlfetch" package to create the *http.Client
// as follows:
//
//	func handler(w http.ResponseWriter, r *http.Request) {
//		c := appengine.NewContext(r)
//		client := urlfetch.Client(c)
//		sender := &gcm.Sender{APIKey: key, Http: client}
//
//		/* ... */
//	}
type Sender struct {
	APIKey     string
	RetryCount int
	HTTPClient *http.Client
}

// NewSender creates a new Sender and sets a timeout on the http.Client
func NewSender(apiKey string, retryCount int, timeout time.Duration) *Sender {
	httpClient := new(http.Client)
	httpClient.Timeout = timeout
	return &Sender{
		APIKey:     apiKey,
		RetryCount: retryCount,
		HTTPClient: httpClient,
	}
}

func (s *Sender) send(m *Message) (*Response, int, error) {
	if err := s.Validate(); err != nil {
		return nil, -1, err
	} else if m == nil {
		return nil, -1, errors.New("Message cannot be nil")
	} else if err := m.Validate(); err != nil {
		return nil, -1, err
	}

	req, err := m.Request()
	if err != nil {
		return nil, -1, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("key=%s", s.APIKey))

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()

	if err := errorMap[resp.StatusCode]; err != nil {
		return nil, -1, err
	}

	if resp.StatusCode >= 500 {
		if retryAfter, err := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 32); err == nil {
			return nil, int(retryAfter), nil
		}
		return nil, 1, nil
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, -1, err
	}

	r := new(Response)
	err = json.Unmarshal(b, r)
	if err != nil {
		return nil, -1, err
	}

	if r.Ok() {
		return r, 0, nil
	}

	if r.Error != "" {
		return r, -1, r.Error
	}

	return r, 1, nil
}

// Send sends a message to the GCM server, retrying in case of
// service unavailability. A non-nil error is returned if a non-recoverable
// error occurs (i.e. if the response status is not "200 OK").
func (s *Sender) Send(m *Message) (*Response, error) {
	r, backoff, err := s.send(m)
	if err != nil {
		return r, err
	}
	for i := 0; i < s.RetryCount; i++ {
		time.Sleep(time.Second * time.Duration(2<<uint(backoff*i)))
		r, backoff, err = s.send(m)
		if err != nil {
			return r, err
		}
		if backoff == 0 {
			return r, nil
		}
		m.update(r)
	}
	return r, errors.New("Retry limit exceeded")
}

// Validate returns an error if the sender is not well-formed
func (s *Sender) Validate() error {
	if s.APIKey == "" {
		return errors.New("the sender's API key must not be empty")
	}
	return nil
}
