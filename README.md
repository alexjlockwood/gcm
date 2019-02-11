Fcm Http client
===

Library uses legacy http protocol for sending messages to devices using device token: https://firebase.google.com/docs/cloud-messaging/send-message#send_messages_using_the_legacy_app_server_protocols

Documentation for initial gcm library, which was replaced with FCM by changing message endpoint: http://godoc.org/github.com/alexjlockwood/gcm

Getting Started
---------------

To install use `go get`:

```bash
go get github.com/Smarp/fcm-http
```

Import with the following:

```go
import "github.com/Smarp/fcm-http"
```

Sample Usage
------------

Here is a quick sample illustrating how to send a message to the FCM server:

```go
package main

import (
	"fmt"
	"net/http"

	fcm "github.com/Smarp/fcm-http"
)

func main() {
	// Create the message to be sent.
	data := map[string]interface{}{"score": "5x1", "time": "15:10"}
	regIDs := []string{"4", "8", "15", "16", "23", "42"}
	msg := fcm.NewMessage(data, regIDs...)

	// Create a Sender to send the message.
	sender := &fcm.Sender{ApiKey: "sample_api_key"}

	// Send the message and receive the response after at most two retries.
	response, err := sender.Send(msg, 2)
	if err != nil {
		fmt.Println("Failed to send message:", err)
		return
	}

	/* ... */
}
```

Note for Google AppEngine users
-------------------------------

If your application server runs on Google AppEngine, you must import the `appengine/urlfetch` package and create the `Sender` as follows:

```go
package sample

import (
	"appengine"
	"appengine/urlfetch"

	fcm "github.com/Smarp/fcm-http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	sender := &fcm.Sender{ApiKey: "sample_api_key", Http: client}

	/* ... */
}        
```
