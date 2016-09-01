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
	// Initial delay before first retry, without jitter.
	backoffInitialDelay = 1000
	// Maximum delay before a retry.
	maxBackoffDelay = 1024000
	// Http method for the api
	apiMethod = "POST"
)

// GcmSendEndpoint is the endpoint for sending messages to the GCM server.
var GcmSendEndpoint = "https://fcm.googleapis.com/fcm/send"

// Sender functions for sending messages to GCM
type Sender interface {
	Send(m *Message) (*Response, error)
}

// Client abstracts the interaction between the application server and the
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
type Client struct {
	APIKey     string
	RetryCount int
	HTTPClient *http.Client
}

// NewSender creates a new Sender and sets a timeout on the http.Client
func NewSender(apiKey string, retryCount int, timeout time.Duration) *Client {
	httpClient := new(http.Client)
	httpClient.Timeout = timeout
	return &Client{
		APIKey:     apiKey,
		RetryCount: retryCount,
		HTTPClient: httpClient,
	}
}

func (c *Client) send(m *Message) (*Response, int, error) {
	if err := c.validate(); err != nil {
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
	req.Header.Add("Authorization", fmt.Sprintf("key=%s", c.APIKey))

	resp, err := c.HTTPClient.Do(req)
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
func (c Client) Send(m *Message) (*Response, error) {
	r, backoff, err := c.send(m)
	if err != nil {
		return r, err
	}
	if backoff == 0 {
		return r, nil
	}

	for i := 0; i < c.RetryCount; i++ {
		time.Sleep(time.Second * time.Duration(2<<uint(backoff*i)))
		r, backoff, err = c.send(m)
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

// validate returns an error if the sender is not well-formed
func (c Client) validate() error {
	if c.APIKey == "" {
		return errors.New("the sender's API key must not be empty")
	}
	return nil
}
